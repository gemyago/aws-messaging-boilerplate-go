package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

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
