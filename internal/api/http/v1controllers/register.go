package v1controllers

import (
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/app"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/di"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		di.ProvideAs[*app.Commands, dummyMessagesCommands],

		NewMessagesController,
	)
}
