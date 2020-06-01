package httpapi

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type StatusHandler struct{}

func NewStatusHandler() *StatusHandler {
	return &StatusHandler{}
}

func (handler *StatusHandler) Health(resp http.ResponseWriter, _ *http.Request) {
	resp.WriteHeader(200)
}

func (handler *StatusHandler) Status(resp http.ResponseWriter, _ *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	message, _ := jsoniter.Marshal(Status{
		ServerStatus:     SERVER_STATUS_HEALTHY,
		TendermintStatus: TENDERMINT_STATUS_HEALTHY,
	})
	http.Error(resp, string(message), 200)
}

type Status struct {
	ServerStatus     ServerStatus     `json:"server_status"`
	TendermintStatus TendermintStatus `json:"tendermint_status"`
}

type ServerStatus = string

var (
	SERVER_STATUS_HEALTHY       ServerStatus = "healthy"
	SERVER_STATUS_SYNCHRONIZING              = "synchronizing"
	SERVER_STATUS_DEGRADED                   = "degraded"
)

type TendermintStatus = string

var (
	TENDERMINT_STATUS_HEALTHY  TendermintStatus = "healthy"
	TENDERMINT_STATUS_DEGRADED                  = "degraded"
	TENDERMINT_STATUS_OFFLINE                   = "offline"
)
