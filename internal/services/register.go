package services

import (
	"context"
	"time"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/di"
	"go.uber.org/dig"
)

func Register(rootCtx context.Context, container *dig.Container) error {
	return di.ProvideAll(container,
		// NewAWSConfigFactory(rootCtx),
		NewTimeProvider,
		di.ProvideValue(time.NewTicker),
		NewShutdownHooks,
	)
}
