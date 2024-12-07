package server

import (
	"net/http"

	"go.uber.org/dig"
)

type MuxRouterAdapter struct {
	mux *http.ServeMux
}

func (r *MuxRouterAdapter) PathValue(req *http.Request, paramName string) string {
	return req.PathValue(paramName)
}

func (r *MuxRouterAdapter) HandleRoute(method, pathPattern string, h http.Handler) {
	r.mux.Handle(method+" "+pathPattern, h)
}

func NewMuxRouterAdapter(mux *http.ServeMux) *MuxRouterAdapter {
	return &MuxRouterAdapter{
		mux: mux,
	}
}

type router interface {
	Handle(pattern string, handler http.Handler)
}

type Group struct {
	dig.Out

	Mount MountFunc `group:"server"`
}

type MountFunc func(r router)
