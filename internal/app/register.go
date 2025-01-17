package app

import (
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/api/http/v1routes/models"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/di"
	"github.com/gemyago/aws-messaging-boilerplate-go/internal/services/awsapi"
	"go.uber.org/dig"
)

type dummyMessagesSendersDeps struct {
	dig.In

	SnsDeps      awsapi.SNSMessageSenderDeps
	EventBusDeps awsapi.EventBusMessageSenderDeps

	DummyMessagesDetailType string `name:"config.aws.eventBus.dummyMessagesDetailType"`
	DummyMessagesTopicARN   string `name:"config.aws.sns.dummyMessagesTopicArn"`
}

type dummySendersOut struct {
	dig.Out

	SendDummySNSMessage         messageSender[models.DummyMessage] `name:"dummy-sns-message-sender"`
	SendDummyEventBridgeMessage messageSender[models.DummyMessage] `name:"dummy-eventbridge-message-sender"`
}

func Register(container *dig.Container) error {
	return di.ProvideAll(
		container,
		func(deps dummyMessagesSendersDeps) dummySendersOut {
			return dummySendersOut{
				SendDummySNSMessage: messageSender[models.DummyMessage](
					awsapi.NewSNSMessageSender[models.DummyMessage](deps.DummyMessagesTopicARN, deps.SnsDeps),
				),
				SendDummyEventBridgeMessage: messageSender[models.DummyMessage](
					awsapi.NewEventBusMessageSender[models.DummyMessage](deps.DummyMessagesDetailType, deps.EventBusDeps),
				),
			}
		},
		NewCommands,
	)
}
