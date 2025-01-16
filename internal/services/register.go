package services

import (
	"context"
	"errors"
	"time"

	"github.com/gemyago/aws-messaging-boilerplate-go/internal/di"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/services/awsapi"
	"go.uber.org/dig"
)

func Register(rootCtx context.Context, container *dig.Container) error {
	return errors.Join(
		awsapi.Register(rootCtx, container),
		di.ProvideAll(container,
			NewTimeProvider,
			di.ProvideValue(time.NewTicker),
			NewShutdownHooks,
		),
	)
}
