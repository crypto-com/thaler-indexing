package httpapi

import "net/http"

type Router interface {
	Use(middlewar func(http.Handler) http.Handler)
	Get(path string, handler func(http.ResponseWriter, *http.Request))
	Handler() http.Handler
}

type RoutePath interface {
	Vars(*http.Request) map[string]string
}
