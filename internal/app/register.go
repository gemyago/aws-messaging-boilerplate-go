package app

import (
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/di"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container, NewCommands)
}
