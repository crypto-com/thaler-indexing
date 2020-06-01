package main

import (
	"os"

	"github.com/crypto-com/chainindex/infrastructure"
	"github.com/crypto-com/chainindex/usecase"
)

type ServerContext struct {
	logger usecase.Logger
	config *Config
}

func NewContext(config *Config) *ServerContext {
	var logger usecase.Logger
	if config.Logger.Color {
		logger = infrastructure.NewZerologLoggerWithColor(os.Stdout)
	} else {
		logger = infrastructure.NewZerologLogger(os.Stdout)
	}

	return &ServerContext{
		logger,
		config,
	}
}
