package queues

import (
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/di"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/services/awsapi"
	"go.uber.org/dig"
)

func Register(container *dig.Container) error {
	return di.ProvideAll(container,
		di.ProvideAs[*awsapi.MessagesPoller, messagePoller],
	)
}
