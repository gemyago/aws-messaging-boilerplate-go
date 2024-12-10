package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

//go:generate mockery --name=MessageSender --filename=mock_message_sender.go --config ../../.mockery-funcs.yaml

type AWSConfigDeps struct {
	Region string `config:"aws.region"`
}

func NewAWSConfigFactory(ctx context.Context) func(deps AWSConfigDeps) (aws.Config, error) {
	return func(deps AWSConfigDeps) (aws.Config, error) {
		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(deps.Region))
		if err != nil {
			return aws.Config{}, fmt.Errorf("failed to load aws configuration, %w", err)
		}
		return cfg, nil
	}
}

type Message struct {
	Id       int64  `json:"id"` //nolint:revive,stylecheck // Id is used to match apigen generated code
	Name     string `json:"name"`
	Comments string `json:"comments,omitempty"`
}

type MessageSender func(ctx context.Context, message *Message) error
