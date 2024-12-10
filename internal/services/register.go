package services

import (
	"context"
	"time"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/di"
	"go.uber.org/dig"
)

func Register(rootCtx context.Context, container *dig.Container) error {
	return di.ProvideAll(container,
		newAWSConfigFactory(rootCtx),
		newSqsClient,
		NewTimeProvider,
		di.ProvideValue(time.NewTicker),
		NewShutdownHooks,
	)
}
