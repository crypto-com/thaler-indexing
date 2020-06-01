package syncservice

import (
	"github.com/crypto-com/chainindex/usecase"
)

type DefaultSyncService struct {
	logger usecase.Logger
	config DefaultSyncServiceConfig

	blocksFeedSubscriber BlocksFeedSubscriber
	blocksProcessor      BlocksProcessor
	blockDataRepoWorker  BlockDataRepoWorker

	syncHeight                 *DefaultRWSerialUint64
	tendermintBlockHeight      *ChRWSerialUint64
	onTendermintHeightUpdateCh chan bool

	onExitCh chan<- bool
}

func NewDefaultSyncService(
	logger usecase.Logger,
	config DefaultSyncServiceConfig,

	blocksFeedSubscriber BlocksFeedSubscriber,
	blocksWorker BlocksProcessor,
	blockDataRepoWorker BlockDataRepoWorker,

	lastSyncHeight uint64,

	onExitCh chan<- bool,
) *DefaultSyncService {
	onTendermintHeightUpdateCh := make(chan bool, 1)
	return &DefaultSyncService{
		logger: logger.WithFields(usecase.LogFields{
			"module": "DefaultSyncService",
		}),
		config: config,

		blocksFeedSubscriber: blocksFeedSubscriber,
		blocksProcessor:      blocksWorker,
		blockDataRepoWorker:  blockDataRepoWorker,

		syncHeight:                 NewDefaultRWSerialUint64(lastSyncHeight),
		tendermintBlockHeight:      NewChRWSerialUint64(onTendermintHeightUpdateCh, lastSyncHeight),
		onTendermintHeightUpdateCh: onTendermintHeightUpdateCh,

		onExitCh: onExitCh,
	}
}

type DefaultSyncServiceConfig struct {
	BlockDataChSize uint
}

func (syncService *DefaultSyncService) Sync() error {
	blockDataCh := make(chan *usecase.BlockData, syncService.config.BlockDataChSize)
	onBlockStoredCh := make(chan uint64, syncService.config.BlockDataChSize)

	go syncService.blocksFeedSubscriber.Run(BlocksFeedSubscriberParams{
		TendermintHeight: syncService.tendermintBlockHeight,

		OnExitCh: syncService.onExitCh,
	})
	go syncService.blocksProcessor.Run(BlocksProcessorParams{
		LastSyncHeight:   syncService.syncHeight.Get(),
		TendermintHeight: syncService.tendermintBlockHeight,

		OnTendermintHeightUpdate: syncService.onTendermintHeightUpdateCh,

		BlockDataCh: blockDataCh,
		OnExitCh:    syncService.onExitCh,
	})
	go syncService.blockDataRepoWorker.Run(BlockDataRepoWorkerParams{
		BlockDataCh:     blockDataCh,
		OnBlockStoredCh: onBlockStoredCh,
		OnExitCh:        syncService.onExitCh,
	})

	go syncService.syncHeightUpdateWorker(onBlockStoredCh)

	return nil
}

func (syncService *DefaultSyncService) syncHeightUpdateWorker(onBlockStoredCh <-chan uint64) {
	for {
		latestSyncedBlockHeight := <-onBlockStoredCh

		syncService.syncHeight.SetIfLarger(latestSyncedBlockHeight)
	}
}

func (syncService *DefaultSyncService) GetStatus() usecase.SyncStatus {
	syncService.tendermintBlockHeight.RLock()
	defer syncService.tendermintBlockHeight.RUnlock()

	syncService.syncHeight.RLock()
	defer syncService.syncHeight.RUnlock()

	return usecase.SyncStatus{
		TendermintBlockHeight: syncService.tendermintBlockHeight.value,
		SyncBlockHeight:       syncService.syncHeight.value,
	}
}
