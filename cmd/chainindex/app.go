package main

import (
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/crypto-com/chainindex/internal/primptr"
	"github.com/crypto-com/chainindex/usecase"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (app *App) Run(args []string) error {
	cliApp := &cli.App{
		Name:                 filepath.Base(args[0]),
		Usage:                "Crypto.com Chain Index",
		Version:              "v0.0.2",
		Copyright:            "(c) 2020 Crypto.com",
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config,c",
				Value: "./config/config.toml",
				Usage: "TOML `FILE` to load configuration from",
			},
			&cli.StringFlag{
				Name:  "logLevel,l",
				Value: "info",
				Usage: "Log level (Default: info, Allowed values: fatal,error,info,debug)",
			},
			&cli.BoolFlag{
				Name:  "color",
				Value: true,
				Usage: "Display colored log",
			},

			&cli.BoolFlag{
				Name:    "dbSSL",
				Usage:   "Enable Postgres SSL mode",
				EnvVars: []string{"DB_SSL"},
			},
			&cli.StringFlag{
				Name:    "dbHost",
				Usage:   "Postgres database hostname",
				EnvVars: []string{"DB_HOST"},
			},
			&cli.UintFlag{
				Name:    "dbPort",
				Usage:   "Postgres database port",
				EnvVars: []string{"DB_PORT"},
			},
			&cli.StringFlag{
				Name:    "dbUsername",
				Usage:   "Postgres username",
				EnvVars: []string{"DB_USERNAME"},
			},
			&cli.StringFlag{
				Name:    "dbName",
				Usage:   "Postgres database name",
				EnvVars: []string{"DB_NAME"},
			},
			&cli.StringFlag{
				Name:    "dbSchema",
				Usage:   "Postgres schema name",
				EnvVars: []string{"DB_SCHEMA"},
			},

			&cli.StringFlag{
				Name:    "tendermintURL",
				Usage:   "Tendermint HTTP RPC URL",
				EnvVars: []string{"TENDERMINT_URL"},
			},
		},
		Action: func(ctx *cli.Context) error {
			var err error

			if args := ctx.Args(); args.Len() > 0 {
				return fmt.Errorf("Unexpected arguments: %q", args.Get(0))
			}

			configPath := ctx.String("config")

			cliConfig := CLIConfig{
				LogLevel: parseLogLevel(ctx.String("logLevel")),

				DatabaseHost:     ctx.String("dbHost"),
				DatabaseUsername: ctx.String("dbUsername"),
				DatabaseName:     ctx.String("dbName"),
				DatabaseSchema:   ctx.String("dbSchema"),

				TendermintHTTPRPCURL: ctx.String("tendermintURL"),
			}
			if ctx.IsSet("color") {
				cliConfig.LoggerColor = primptr.Bool(ctx.Bool("color"))
			}
			if ctx.IsSet("dbSSL") {
				cliConfig.DatabaseSSL = primptr.Bool(ctx.Bool("dbSSL"))
			}
			if ctx.IsSet("dgPort") {
				cliConfig.DatabasePort = primptr.Uint32(uint32(ctx.Uint("dbPort")))
			}

			serverApp, err := NewServer(configPath, &cliConfig)
			if err != nil {
				return fmt.Errorf("error creating server: %v", err)
			}

			if err = serverApp.Run(); err != nil {
				return fmt.Errorf("Error when starting server: %v", err)
			}

			return nil
		},
	}

	err := cliApp.Run(args)
	if err != nil {
		return err
	}

	return nil
}

func parseLogLevel(level string) usecase.LogLevel {
	switch level {
	case "panic":
		return usecase.LOG_LEVEL_PANIC
	case "error":
		return usecase.LOG_LEVEL_ERROR
	case "info":
		return usecase.LOG_LEVEL_INFO
	case "debug":
		return usecase.LOG_LEVEL_DEBUG
	default:
		panic("Unsupported log level: " + level)
	}
}
