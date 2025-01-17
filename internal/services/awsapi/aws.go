package awsapi

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/dig"
)

//go:generate mockery --name=MessageSender --filename=mock_message_sender.go --config ../../../.mockery-funcs.yaml

type AWSConfigDeps struct {
	dig.In `ignore-unexported:"true"`

	Region       string `name:"config.aws.region"`
	BaseEndpoint string `name:"config.aws.baseEndpoint" optional:"true"`

	loadOpts []func(*config.LoadOptions) error
}

func newAWSConfigFactory(ctx context.Context) func(deps AWSConfigDeps) (aws.Config, error) {
	return func(deps AWSConfigDeps) (aws.Config, error) {
		opts := append([]func(*config.LoadOptions) error{
			config.WithRegion(deps.Region),
		}, deps.loadOpts...)
		// BaseEndpoint is defined for local/test modes only and points on localstack instance
		if deps.BaseEndpoint != "" {
			opts = append(opts,
				config.WithBaseEndpoint(deps.BaseEndpoint),
				config.WithCredentialsProvider(aws.AnonymousCredentials{}),
			)
		}
		cfg, err := config.LoadDefaultConfig(ctx, opts...)
		if err != nil {
			return aws.Config{}, fmt.Errorf("failed to load aws configuration, %w", err)
		}
		return cfg, nil
	}
}

type MessageSender[TMessage any] func(
	ctx context.Context,
	message *TMessage,
) error
