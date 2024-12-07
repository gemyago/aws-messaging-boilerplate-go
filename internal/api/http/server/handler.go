package server

import (
	"log/slog"
	"net/http"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/middleware"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/routes"
	sloghttp "github.com/samber/slog-http"
	"go.uber.org/dig"
)

type RouterAdapter struct {
	mux *http.ServeMux
}

func (r *RouterAdapter) PathValue(req *http.Request, paramName string) string {
	return req.PathValue(paramName)
}

func (r *RouterAdapter) HandleRoute(method, pathPattern string, h http.Handler) {
	r.mux.Handle(method+" "+pathPattern, h)
}

func NewMuxRouterAdapter(mux *http.ServeMux) *RouterAdapter {
	return &RouterAdapter{
		mux: mux,
	}
}

type RootHandlerDeps struct {
	dig.In

	RootLogger *slog.Logger
	Groups     []routes.MountFunc `group:"server"`
}

func NewRootHandler(deps RootHandlerDeps) http.Handler {
	mux := http.NewServeMux()

	for _, grp := range deps.Groups {
		grp(mux)
	}

	// Router wire-up
	chain := middleware.Chain(
		middleware.NewTracingMiddleware(middleware.NewTracingMiddlewareCfg()),
		sloghttp.NewWithConfig(deps.RootLogger, sloghttp.Config{
			DefaultLevel:     slog.LevelInfo,
			ClientErrorLevel: slog.LevelWarn,
			ServerErrorLevel: slog.LevelError,

			WithUserAgent:      true,
			WithRequestID:      false, // We handle it ourselves (tracing middleware)
			WithRequestHeader:  true,
			WithResponseHeader: true,
			WithSpanID:         true,
			WithTraceID:        true,
		}),
		middleware.NewRecovererMiddleware(deps.RootLogger),
	)
	return chain(mux)
}
