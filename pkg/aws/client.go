// pkg/aws/client.go
package aws

import (
	"context"
	"strings"

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

func (c *ECSClient) ListTasks(ctx context.Context) ([]*types.Task, error) {
	var tasks []*types.Task
	var nextToken *string

	for {
		// List task ARNs
		input := &ecs.ListTasksInput{
			Cluster:    &c.Context.Cluster,
			NextToken:  nextToken,
			MaxResults: aws.Int32(100),
		}

		result, err := c.Client.ListTasks(ctx, input)
		if err != nil {
			return nil, err
		}

		if len(result.TaskArns) == 0 {
			break
		}

		// Describe tasks to get detailed information
		describeInput := &ecs.DescribeTasksInput{
			Cluster: &c.Context.Cluster,
			Tasks:   result.TaskArns,
		}

		describeResult, err := c.Client.DescribeTasks(ctx, describeInput)
		if err != nil {
			return nil, err
		}

		// Convert to our task type
		for _, task := range describeResult.Tasks {
			// Extract task ID from ARN
			taskId := extractTaskId(*task.TaskArn)

			taskDefParts := strings.Split(*task.TaskDefinitionArn, "/")
			taskDefFamily := taskDefParts[len(taskDefParts)-1]

			t := &types.Task{
				TaskID:        taskId,
				TaskARN:       *task.TaskArn,
				Status:        string(*task.LastStatus),
				TaskDefFamily: taskDefFamily,
				LastStatus:    string(*task.LastStatus),
				DesiredStatus: string(*task.DesiredStatus),
				CreatedAt:     *task.CreatedAt,
				Group:         *task.Group,
			}

			if task.StartedAt != nil {
				t.StartedAt = *task.StartedAt
			}

			if task.ContainerInstanceArn != nil {
				t.ContainerInstanceARN = *task.ContainerInstanceArn
			}

			tasks = append(tasks, t)
		}

		if result.NextToken == nil {
			break
		}
		nextToken = result.NextToken
	}

	return tasks, nil
}

// StopTask stops a task in the cluster
func (c *ECSClient) StopTask(ctx context.Context, taskId string) error {
	input := &ecs.StopTaskInput{
		Cluster: &c.Context.Cluster,
		Task:    &taskId,
		Reason:  aws.String("Stopped via ECS CLI"),
	}

	_, err := c.Client.StopTask(ctx, input)
	return err
}

// Helper function to extract task ID from ARN
func extractTaskId(arn string) string {
	parts := strings.Split(arn, "/")
	return parts[len(parts)-1]
}
