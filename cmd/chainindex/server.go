package main

import (
	"fmt"

	"github.com/crypto-com/chainindex/adapter"
	httpapiadapter "github.com/crypto-com/chainindex/adapter/httpapi"
	"github.com/crypto-com/chainindex/adapter/rdbviewrepo"
	"github.com/crypto-com/chainindex/adapter/syncservice"
	tendermintadapter "github.com/crypto-com/chainindex/adapter/tendermint"
	"github.com/crypto-com/chainindex/infrastructure"
	"github.com/crypto-com/chainindex/infrastructure/httpapi"
	"github.com/crypto-com/chainindex/infrastructure/tendermint"
	"github.com/crypto-com/chainindex/internal/filereader/toml"
	"github.com/crypto-com/chainindex/usecase"
	"github.com/crypto-com/chainindex/usecase/viewrepo"
)

type Server struct {
	*ServerContext
}

func NewServer(configPath string, cliConfig *CLIConfig) (*Server, error) {
	configReader, configFileErr := toml.FromFile(configPath)
	if configFileErr != nil {
		return nil, configFileErr
	}

	var config Config
	readConfigErr := configReader.Read(&config)
	if readConfigErr != nil {
		return nil, readConfigErr
	}
	config.OverrideByCLIConfig(cliConfig)

	context := NewContext(&config)
	context.logger.SetLogLevel(cliConfig.LogLevel)

	return &Server{
		ServerContext: context,
	}, nil
}

func (server *Server) Run() error {
	tendermintClient := tendermint.NewHTTPClient(server.config.Tendermint.URL)

	pgxConnPool, err := infrastructure.NewPgxConnPool(infrastructure.PgxConnPoolConfig{
		PgxConnConfig: infrastructure.PgxConnConfig{
			Host:     server.config.Database.Host,
			Port:     server.config.Database.Port,
			Username: server.config.Database.Username,
			Password: server.config.Database.Password,
			Database: server.config.Database.Name,
			SSL:      server.config.Database.SSL,
		},
		MaxConns:          server.config.Postgres.MaxConns,
		MinConns:          server.config.Postgres.MinConns,
		MaxConnLifeTime:   server.config.Postgres.MaxConnLifeTime.Duration,
		MaxConnIdleTime:   server.config.Postgres.MaxConnIdleTime.Duration,
		HealthCheckPeriod: server.config.Postgres.HealthCheckInterval.Duration,
	}, server.logger)
	if err != nil {
		return fmt.Errorf("error creating connection pool to Postgres: %v", err)
	}
	rDbConn := infrastructure.NewPgxRDbConn(pgxConnPool)
	rDBTypeConv := new(infrastructure.PgxRDbTypeConv)

	_, err = rDbConn.Exec("SELECT 1")
	if err != nil {
		server.logger.Panicf("error connecting to Database: %v", err)
	}

	blockActivityDataRepo := adapter.NewDefaultRDbBlockActivityDataRepo(infrastructure.PostgresStmtBuilder, rDBTypeConv)
	blockDataRepo := adapter.NewRDbBlockDataRepo(rDbConn, infrastructure.PostgresStmtBuilder, rDBTypeConv, blockActivityDataRepo)
	blockViewRepo := rdbviewrepo.NewRDbBlockViewRepo(rDbConn, infrastructure.PostgresStmtBuilder, rDBTypeConv)
	activityViewRepo := rdbviewrepo.NewRDbActivityViewRepo(rDbConn, infrastructure.PostgresStmtBuilder, rDBTypeConv)
	councilNodeViewRepo := rdbviewrepo.NewRDbCouncilNodeViewRepo(rDbConn, infrastructure.PostgresStmtBuilder, rDBTypeConv)
	rewardViewRepo := rdbviewrepo.NewRDbRewardViewRepo(rDbConn, infrastructure.PostgresStmtBuilder, rDBTypeConv)
	stakingAccountViewRepo := rdbviewrepo.NewRDbStkaingAccountViewRepo(rDbConn, infrastructure.PostgresStmtBuilder, rDBTypeConv)

	// TODO: channel to receive error instead of exit
	onExitCh := make(chan bool)

	syncService := server.getDefaultSyncService(
		tendermintClient,
		blockDataRepo,
		blockViewRepo,

		onExitCh,
	)
	if err := syncService.Sync(); err != nil {
		server.logger.Panicf("error starting sync service: %v", err)
	}

	server.startHTTPAPIServer(
		syncService,
		activityViewRepo,
		rewardViewRepo,
		blockViewRepo,
		councilNodeViewRepo,
		stakingAccountViewRepo,

		onExitCh,
	)

	<-onExitCh

	return nil
}

func (server *Server) getDefaultSyncService(
	tendermintClient tendermintadapter.Client,
	blockDataRepo usecase.BlockDataRepository,
	blockViewRepo viewrepo.BlockViewRepo,

	onExitCh chan<- bool,
) usecase.SyncService {
	config := syncservice.DefaultSyncServiceConfig{
		BlockDataChSize: server.config.Synchronization.BlockDataChSize,
	}

	tendermintBlocksFeed := tendermint.NewPollingBlocksFeed(
		server.logger,
		tendermintClient,
		tendermint.PollingBlocksFeedOptions{
			PollingInterval: server.config.Synchronization.BlockHeightPollingInterval.Duration,
		},
	)
	blocksFeedSubscriber := syncservice.NewDefaultBlocksFeedSubscriber(
		server.logger,
		tendermintBlocksFeed,
	)
	// blocksWorker := syncservice.NewPeriodicBlocksWorker(server.logger, tendermintClient)
	blocksWorker := syncservice.NewBatchBlocksProcessor(
		server.logger,
		tendermintClient,
		int(server.config.Synchronization.MaxConcurrentBlockWorker),
	)
	repoWorker := syncservice.NewDefaultBlockDataRepoWorker(
		server.logger,
		blockDataRepo,
	)

	lastSyncHeight, err := blockViewRepo.LatestBlockHeight()
	if err != nil {
		server.logger.Panicf("error getting last sync height: %v", err)
	}

	return syncservice.NewDefaultSyncService(
		server.logger,
		config,

		blocksFeedSubscriber,
		blocksWorker,
		repoWorker,

		lastSyncHeight,

		onExitCh,
	)
}

func (server *Server) startHTTPAPIServer(
	syncService usecase.SyncService,
	activityViewRepo viewrepo.ActivityViewRepo,
	rewardViewRepo viewrepo.RewardViewRepo,
	blockViewRepo viewrepo.BlockViewRepo,
	councilNodeViewRepo viewrepo.CouncilNodeViewRepo,
	stakingAccountViewRepo viewrepo.StakingAccountViewRepo,

	onExitCh chan<- bool,
) {
	router := httpapi.NewMuxRouter()

	routePath := httpapi.NewMuxRoutePath()

	statusHandler := httpapiadapter.NewStatusHandler()

	activitiesHandler := httpapiadapter.NewActivitiesHandler(
		server.logger,
		routePath,
		activityViewRepo,
	)
	blocksHandler := httpapiadapter.NewBlocksHandler(
		server.logger,
		routePath,
		blockViewRepo,
	)
	councilNodeHandler := httpapiadapter.NewCouncilNodesHandler(
		server.logger,
		routePath,
		councilNodeViewRepo,
	)
	chainStatusHandler := httpapiadapter.NewChainStatusHandler(
		server.logger,
		syncService,
		activityViewRepo,
		rewardViewRepo,
		councilNodeViewRepo,
	)
	searchHandler := httpapiadapter.NewSearchHandler(
		server.logger,
		activityViewRepo,
		blockViewRepo,
		stakingAccountViewRepo,
		councilNodeViewRepo,
	)

	httpapiadapter.NewRoutesRegistry(
		router,

		statusHandler,
		chainStatusHandler,
		activitiesHandler,
		blocksHandler,
		councilNodeHandler,
		searchHandler,
	).RegisterHandlers()

	router.Use(httpapiadapter.LoggerMiddleware(server.logger))
	httpAPIServer := httpapi.NewServer(router.Handler(), httpapi.ServerConfig{
		WriteTimeout: server.config.HTTPAPI.WriteTimeout.Duration,
		ReadTimeout:  server.config.HTTPAPI.ReadTimeout.Duration,
		IdleTimeout:  server.config.HTTPAPI.IdleTimeout.Duration,
	})

	server.logger.Info("HTTP API server start listening on " + server.config.HTTPAPI.ListeningAddress)
	if err := httpAPIServer.ListenAndServe(server.config.HTTPAPI.ListeningAddress); err != nil {
		server.logger.Errorf("error listening and serving HTTP API server: %v", err)
	}

	onExitCh <- true
}
