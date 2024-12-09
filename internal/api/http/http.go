package http

import (
	"errors"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/server"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1controllers"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/handlers"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/di"
	"go.uber.org/dig"
)

// Use apigen to generate v1routes
//go:generate apigen ./v1schemas/api.yaml ./v1routes --verbose

func NewV1RoutesApp(router *server.MuxRouterAdapter) *handlers.HTTPApp {
	return handlers.NewHTTPApp(router)
}

func Register(container *dig.Container) error {
	return errors.Join(
		v1controllers.Register(container),
		di.ProvideAll(container,
			NewV1RoutesApp,
			server.MakeHandlersGroupFactory(handlers.RegisterMessagesRoutes),
		),
	)
}
