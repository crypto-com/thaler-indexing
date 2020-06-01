package httpapi

import (
	"net/http"

	"github.com/crypto-com/chainindex/usecase"
)

func LoggerMiddleware(logger usecase.Logger) func(http.Handler) http.Handler {
	logger = logger.WithFields(usecase.LogFields{
		"module": "httpapi",
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			logger.Infof("%s %s", req.Method, req.RequestURI)

			next.ServeHTTP(resp, req)
		})
	}
}
