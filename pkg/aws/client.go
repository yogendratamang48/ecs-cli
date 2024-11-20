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

func (c *ECSClient) DescribeServices(ctx context.Context, serviceNames []string) ([]*types.ServiceDetail, error) {
	input := &ecs.DescribeServicesInput{
		Cluster:  &c.Context.Cluster,
		Services: serviceNames,
	}

	result, err := c.Client.DescribeServices(ctx, input)
	if err != nil {
		return nil, err
	}

	var services []*types.ServiceDetail
	for _, svc := range result.Services {
		// Create load balancer info
		var loadBalancers []types.LoadBalancer
		for _, lb := range svc.LoadBalancers {
			loadBalancers = append(loadBalancers, types.LoadBalancer{
				Type:          string(*lb.TargetGroupArn),
				TargetGroup:   *lb.TargetGroupArn,
				ContainerName: *lb.ContainerName,
				ContainerPort: int(*lb.ContainerPort),
			})
		}

		// Create network config
		networkConfig := types.NetworkConfig{
			Type: string(svc.NetworkConfiguration.AwsvpcConfiguration.AssignPublicIp),
		}
		if svc.NetworkConfiguration != nil && svc.NetworkConfiguration.AwsvpcConfiguration != nil {
			networkConfig.SubnetIds = svc.NetworkConfiguration.AwsvpcConfiguration.Subnets
			networkConfig.SecurityGroups = svc.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups
			networkConfig.PublicIP = string(svc.NetworkConfiguration.AwsvpcConfiguration.AssignPublicIp)
		}

		// Create events
		var events []types.ServiceEvent
		for _, event := range svc.Events {
			events = append(events, types.ServiceEvent{
				CreatedAt: *event.CreatedAt,
				Message:   *event.Message,
			})
		}

		serviceDetail := &types.ServiceDetail{
			Name:          *svc.ServiceName,
			Status:        string(*svc.Status),
			TaskDef:       *svc.TaskDefinition,
			DesiredCount:  int(svc.DesiredCount),
			RunningCount:  int(svc.RunningCount),
			PendingCount:  int(svc.PendingCount),
			CreatedAt:     *svc.CreatedAt,
			LoadBalancers: loadBalancers,
			NetworkConfig: networkConfig,
			Events:        events,
		}

		services = append(services, serviceDetail)
	}

	return services, nil
}

func (c *ECSClient) DescribeTasks(ctx context.Context, taskIds []string) ([]*types.TaskDetail, error) {
	input := &ecs.DescribeTasksInput{
		Cluster: &c.Context.Cluster,
		Tasks:   taskIds,
	}

	result, err := c.Client.DescribeTasks(ctx, input)
	if err != nil {
		return nil, err
	}

	var tasks []*types.TaskDetail

	for _, task := range result.Tasks {
		taskDetail := &types.TaskDetail{
			TaskID:            extractTaskId(*task.TaskArn),
			TaskARN:           *task.TaskArn,
			ClusterARN:        *task.ClusterArn,
			TaskDefinitionARN: *task.TaskDefinitionArn,
			Status:            string(*task.LastStatus),
			DesiredStatus:     string(*task.DesiredStatus),
			CreatedAt:         *task.CreatedAt,
			Group:             *task.Group,
			LaunchType:        string(task.LaunchType),
		}

		if task.ContainerInstanceArn != nil {
			taskDetail.ContainerInstanceARN = *task.ContainerInstanceArn
		}

		if task.Cpu != nil {
			taskDetail.CPU = *task.Cpu
		}

		if task.Memory != nil {
			taskDetail.Memory = *task.Memory
		}

		if task.StartedAt != nil {
			taskDetail.StartedAt = *task.StartedAt
		}

		if task.StoppedAt != nil {
			taskDetail.StoppedAt = *task.StoppedAt
		}

		if task.StoppedReason != nil {
			taskDetail.StoppedReason = *task.StoppedReason
		}

		// Add container details
		// fmt.Println(string(*task.Containers[0].Name))
		for _, container := range task.Containers {
			if strings.HasPrefix(string(*container.Name), "ecs-service-connect-") {
				continue
			}
			containerDetail := types.ContainerDetail{
				Name:         string(*container.Name),
				Image:        string(*container.Image),
				Status:       string(*container.LastStatus),
				CreatedAt:    *task.CreatedAt,
				HealthStatus: string(container.HealthStatus),
			}

			if container.RuntimeId != nil {
				containerDetail.RuntimeID = *container.RuntimeId
			}

			if container.ExitCode != nil {
				containerDetail.ExitCode = container.ExitCode
			}

			// Add port bindings
			for _, binding := range container.NetworkBindings {
				containerDetail.NetworkBindings = append(containerDetail.NetworkBindings, types.PortBinding{
					ContainerPort: *binding.ContainerPort,
					HostPort:      *binding.HostPort,
					Protocol:      string(binding.Protocol),
				})
			}

			taskDetail.Containers = append(taskDetail.Containers, containerDetail)
		}

		// Add network interfaces
		if task.Attachments != nil {
			for _, attachment := range task.Attachments {
				if *attachment.Type == "ElasticNetworkInterface" {
					var networkInterface types.NetworkInterface
					for _, detail := range attachment.Details {
						switch *detail.Name {
						case "networkInterfaceId":
							networkInterface.AttachmentID = *detail.Value
						case "privateIPv4Address":
							networkInterface.PrivateIPv4 = *detail.Value
						case "publicIPv4Address":
							networkInterface.PublicIPv4 = *detail.Value
						case "subnetId":
							networkInterface.SubnetID = *detail.Value
						}
					}
					taskDetail.NetworkInterfaces = append(taskDetail.NetworkInterfaces, networkInterface)
				}
			}
		}

		tasks = append(tasks, taskDetail)
	}

	return tasks, nil
}
