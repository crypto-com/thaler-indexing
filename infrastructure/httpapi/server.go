package httpapi

import (
	"net/http"
	"time"
)

type Server struct {
	handler http.Handler

	config ServerConfig
}

type ServerConfig struct {
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
}

func NewServer(router http.Handler, config ServerConfig) *Server {
	return &Server{
		router,

		config,
	}
}

func (server *Server) ListenAndServe(address string) error {
	httpServer := &http.Server{
		Addr:         address,
		WriteTimeout: server.config.WriteTimeout,
		ReadTimeout:  server.config.ReadTimeout,
		IdleTimeout:  server.config.IdleTimeout,
		Handler:      server.handler,
	}

	if err := httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
