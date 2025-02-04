package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/yogendratamang48/ecs/pkg/types"
)

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
				TaskId:        taskId,
				TaskArn:       *task.TaskArn,
				Status:        string(*task.LastStatus),
				TaskDefFamily: taskDefFamily,
				LastStatus:    string(*task.LastStatus),
				DesiredStatus: string(*task.DesiredStatus),
				CreatedAt:     *task.CreatedAt,
				Group:         *task.Group,
				Cpu:           *task.Cpu,
				Memory:        *task.Memory,
				LaunchType:    string(task.LaunchType),
				CapacityProvider: func() string {
					if task.CapacityProviderName != nil {
						return string(*task.CapacityProviderName)
					}
					return "-"
				}(),
			}

			if task.StartedAt != nil {
				t.StartedAt = *task.StartedAt
			}

			if task.ContainerInstanceArn != nil {
				t.ContainerInstanceArn = *task.ContainerInstanceArn
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
			TaskId:            extractTaskId(*task.TaskArn),
			TaskArn:           *task.TaskArn,
			ClusterArn:        *task.ClusterArn,
			Cpu:               *task.Cpu,
			Memory:            *task.Memory,
			TaskDefinitionArn: *task.TaskDefinitionArn,
			Status:            string(*task.LastStatus),
			DesiredStatus:     string(*task.DesiredStatus),
			CreatedAt:         *task.CreatedAt,
			Group:             *task.Group,
			LaunchType:        string(task.LaunchType),
			CapacityProvider: func() string {
				if task.CapacityProviderName != nil {
					return string(*task.CapacityProviderName)
				}
				return "-"
			}(),
		}

		if task.ContainerInstanceArn != nil {
			taskDetail.ContainerInstanceArn = *task.ContainerInstanceArn
		}

		if task.Cpu != nil {
			taskDetail.Cpu = *task.Cpu
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
