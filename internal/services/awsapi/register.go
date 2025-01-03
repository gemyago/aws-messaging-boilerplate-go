package awsapi

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gemyago/aws-sqs-boilerplate-go/internal/di"
	"go.uber.org/dig"
)

func Register(rootCtx context.Context, container *dig.Container) error {
	return di.ProvideAll(container,
		newAWSConfigFactory(rootCtx),
		sqs.NewFromConfig,
		sns.NewFromConfig,
		eventbridge.NewFromConfig,
		NewMessagesPoller,
	)
}
