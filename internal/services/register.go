package services

import (
	"time"

	"github.com/gemyago/aws-sqs-boilerplate-go/internal/di"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		NewTimeProvider,
		di.ProvideValue(time.NewTicker),
		NewShutdownHooks,
	)
}
