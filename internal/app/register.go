package app

import (
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/api/http/v1routes/models"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/di"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/services/awsapi"
	"go.uber.org/dig"
)

type dummyMessagesSnsDeps struct {
	dig.In

	SnsDeps  awsapi.SNSMessageSenderDeps
	TopicARN string `name:"config.aws.sns.dummyMessagesTopicArn"`
}

func Register(container *dig.Container) error {
	return di.ProvideAll(
		container,
		func(deps dummyMessagesSnsDeps) messageSender[models.DummyMessage] {
			return messageSender[models.DummyMessage](
				awsapi.NewSNSMessageSender[models.DummyMessage](deps.TopicARN, deps.SnsDeps),
			)
		},
		NewCommands,
	)
}
