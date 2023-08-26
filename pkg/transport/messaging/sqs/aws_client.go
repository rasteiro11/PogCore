package sqs

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func newEndpointResolver() aws.EndpointResolverWithOptionsFunc {
	awsEndpoint := os.Getenv("AWS_ENDPOINT")
	awsRegion := os.Getenv("AWS_REGION")

	return func(_, region string, options ...interface{}) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}

		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}
}

func newProviderConfiguration(ctx context.Context) (aws.Config, error) {
	options := []func(*config.LoadOptions) error{
		config.WithEndpointResolverWithOptions(newEndpointResolver()),
	}

	return config.LoadDefaultConfig(ctx, options...)
}

func newProvider(ctx context.Context) (*sqs.Client, error) {
	cfg, err := newProviderConfiguration(ctx)
	if err != nil {
		return nil, err
	}

	return sqs.NewFromConfig(cfg), nil
}
