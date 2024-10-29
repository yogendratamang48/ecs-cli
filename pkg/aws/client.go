// pkg/aws/client.go
package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/yogendratamang48/ecs/pkg/types"
)

// ECSClient wraps the ECS client and provides additional context
type ECSClient struct {
	*ecs.Client
	Context *types.Context
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

// ListServices returns all services in the cluster
func (c *ECSClient) ListServices(ctx context.Context) ([]*types.Service, error) {
	var services []*types.Service
	var nextToken *string

	for {
		// List service ARNs
		input := &ecs.ListServicesInput{
			Cluster:    &c.Context.Cluster,
			NextToken:  nextToken,
			MaxResults: aws.Int32(100),
		}

		result, err := c.Client.ListServices(ctx, input)
		if err != nil {
			return nil, err
		}

		if len(result.ServiceArns) == 0 {
			break
		}

		// Describe services to get detailed information
		describeInput := &ecs.DescribeServicesInput{
			Cluster:  &c.Context.Cluster,
			Services: result.ServiceArns,
		}

		describeResult, err := c.Client.DescribeServices(ctx, describeInput)
		if err != nil {
			return nil, err
		}

		// Convert to our service type
		for _, svc := range describeResult.Services {
			services = append(services, &types.Service{
				Name:         *svc.ServiceName,
				Status:       string(*svc.Status),
				TaskDef:      *svc.TaskDefinition,
				DesiredCount: int(svc.DesiredCount),
				RunningCount: int(svc.RunningCount),
				PendingCount: int(svc.PendingCount),
				CreatedAt:    *svc.CreatedAt,
			})
		}

		if result.NextToken == nil {
			break
		}
		nextToken = result.NextToken
	}

	return services, nil
}
