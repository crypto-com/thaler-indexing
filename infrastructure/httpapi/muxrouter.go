package httpapi

import (
	"net/http"

	"github.com/gorilla/mux"
)

type MuxRouter struct {
	instance *mux.Router
}

func NewMuxRouter() *MuxRouter {
	return &MuxRouter{
		instance: mux.NewRouter(),
	}
}

func (router *MuxRouter) Use(middleware func(http.Handler) http.Handler) {
	router.instance.Use(middleware)
}

func (router *MuxRouter) Get(path string, handler func(http.ResponseWriter, *http.Request)) {
	router.instance.HandleFunc(path, handler).Methods("GET")
}

func (router *MuxRouter) Handler() http.Handler {
	return router.instance
}

type MuxRoutePath struct{}

func NewMuxRoutePath() *MuxRoutePath {
	return &MuxRoutePath{}
}

func (path *MuxRoutePath) Vars(req *http.Request) map[string]string {
	return mux.Vars(req)
}
