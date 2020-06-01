package main

import (
	"os"
	"time"

	"github.com/crypto-com/chainindex/usecase"
)

type Config struct {
	LogLevel usecase.LogLevel

	FileConfig
}

// TODO: Refactor to read an slice of all CLI and yaml options
// Override and merge configuration with CLI configuration
func (config *Config) OverrideByCLIConfig(cliConfig *CLIConfig) {
	if cliConfig.LoggerColor != nil {
		config.Logger.Color = *cliConfig.LoggerColor
	}
	if cliConfig.DatabaseSSL != nil {
		config.Database.SSL = *cliConfig.DatabaseSSL
	}
	if cliConfig.DatabaseHost != "" {
		config.Database.Host = cliConfig.DatabaseHost
	}
	if cliConfig.DatabasePort != nil {
		config.Database.Port = *cliConfig.DatabasePort
	}
	if cliConfig.DatabaseUsername != "" {
		config.Database.Username = cliConfig.DatabaseUsername
	}
	if cliConfig.DatabaseName != "" {
		config.Database.Name = cliConfig.DatabaseName
	}
	if cliConfig.DatabaseSchema != "" {
		config.Database.Schema = cliConfig.DatabaseSchema
	}
	config.Database.Password = os.Getenv("DB_PASSWORD")

	if cliConfig.TendermintHTTPRPCURL != "" {
		config.Tendermint.URL = cliConfig.TendermintHTTPRPCURL
	}
}

type FileConfig struct {
	Logger          LoggerConfig
	HTTPAPI         HTTPAPIConfig
	Tendermint      TendermintConfig
	Database        DatabaseConfig
	Synchronization SyncConfig
	Postgres        PostgresConfig
}

type HTTPAPIConfig struct {
	ListeningAddress string   `toml:"listening_address"`
	WriteTimeout     duration `toml:"write_timeout"`
	ReadTimeout      duration `toml:"read_timeout"`
	IdleTimeout      duration `toml:"idle_timeout"`
}

type LoggerConfig struct {
	Color bool `toml:"color"`
}

type TendermintConfig struct {
	URL string `toml:"http_rpc_url"`
}

type DatabaseConfig struct {
	SSL      bool   `toml:"ssl"`
	Host     string `toml:"host"`
	Port     uint32 `toml:"port"`
	Username string `toml:"username"`
	Password string
	Name     string `toml:"name"`
	Schema   string `toml:"schema"`
}

type SyncConfig struct {
	BlockDataChSize            uint     `toml:"block_data_channel_size"`
	BlockHeightPollingInterval duration `toml:"block_height_polling_interval"`
	BlockHeightChSize          uint     `toml:"block_height_channel_size"`
	MaxConcurrentBlockWorker   uint     `toml:"max_concurrent_block_worker"`
}

type PostgresConfig struct {
	MaxConns            int32    `toml:"pool_max_conns"`
	MinConns            int32    `toml:"pool_min_conns"`
	MaxConnLifeTime     duration `toml:"pool_max_conn_lifetime"`
	MaxConnIdleTime     duration `toml:"pool_max_conn_idle_time"`
	HealthCheckInterval duration `toml:"pool_health_check_interval"`
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type CLIConfig struct {
	LoggerColor *bool
	LogLevel    usecase.LogLevel

	DatabaseSSL      *bool
	DatabaseHost     string
	DatabasePort     *uint32
	DatabaseUsername string
	DatabaseName     string
	DatabaseSchema   string

	TendermintHTTPRPCURL string
}
