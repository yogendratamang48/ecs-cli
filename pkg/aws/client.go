// pkg/aws/client.go
package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/yogendratamang48/ecs/pkg/types"
)

// ECSClient wraps the ECS client and provides additional context
type ECSClient struct {
	*ecs.Client
	Context *types.Context
}

type CloudWatchClient struct {
	*cloudwatchlogs.Client
	Context *types.Context
}

type SSMClient struct {
	*ssm.Client
	Context *types.Context
}

func NewCloudWatchLogsClient(ctx *types.Context) (*CloudWatchClient, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(ctx.Region),
		config.WithSharedConfigProfile(ctx.Profile),
	)
	if err != nil {
		return nil, err
	}

	return &CloudWatchClient{
		Client:  cloudwatchlogs.NewFromConfig(cfg),
		Context: ctx,
	}, nil
}

// NewECSClient creates a new ECS client with the given context
func NewECSClient(ctx *types.Context) (*ECSClient, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(ctx.Region),
		config.WithSharedConfigProfile(ctx.Profile),
	)
	if err != nil {
		return nil, err
	}

	return &ECSClient{
		Client:  ecs.NewFromConfig(cfg),
		Context: ctx,
	}, nil
}

func NewSSMClient(ctx *types.Context) (*SSMClient, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(ctx.Region),
		config.WithSharedConfigProfile(ctx.Profile),
	)
	if err != nil {
		return nil, err
	}

	return &SSMClient{
		Client:  ssm.NewFromConfig(cfg),
		Context: ctx,
	}, nil
}
