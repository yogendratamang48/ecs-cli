// pkg/aws/logs.go
package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type LogEvent struct {
	Timestamp int64
	Message   string
}

// GetTaskLogs retrieves logs for a specific task
func (c *ECSClient) GetTaskLogs(ctx context.Context, taskID string, follow bool, since time.Duration, container string) (<-chan LogEvent, error) {
	// Get task details to find the log configuration
	input := &ecs.DescribeTasksInput{
		Cluster: &c.Context.Cluster,
		Tasks:   []string{taskID},
	}

	result, err := c.Client.DescribeTasks(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe task: %w", err)
	}

	if len(result.Tasks) == 0 {
		return nil, fmt.Errorf("task %s not found", taskID)
	}

	task := result.Tasks[0]
	tfResult, err := c.Client.DescribeTaskDefinition(ctx, &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: task.TaskDefinitionArn,
	})
	if err != nil {
		return nil, fmt.Errorf("error loading task definition")
	}

	// Find the container to get logs from
	var targetContainer *ecsTypes.ContainerDefinition
	for _, c := range tfResult.TaskDefinition.ContainerDefinitions {
		if container == "" || *c.Name == container {
			targetContainer = &c
			break
		}
	}

	if targetContainer == nil {
		return nil, fmt.Errorf("container %s not found in task", container)
	}

	// Extract log configuration
	logConfig := targetContainer.LogConfiguration
	if logConfig == nil || string(logConfig.LogDriver) != "awslogs" {
		return nil, fmt.Errorf("awslogs driver not configured for container")
	}

	options := logConfig.Options
	logGroup := options["awslogs-group"]
	logStream := fmt.Sprintf("%s/%s/%s", options["awslogs-stream-prefix"], *targetContainer.Name, taskID)

	// Create CloudWatch Logs client
	cwlClient, err := NewCloudWatchLogsClient(c.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to create Cloudwatch client: %w", err)
	}

	// Calculate start time
	startTime := time.Now().Add(-since).UnixMilli()

	// Create channel for log events
	logChan := make(chan LogEvent)

	// Start goroutine to fetch logs
	go func() {
		defer close(logChan)

		var nextToken *string
		for {
			getLogsInput := &cloudwatchlogs.GetLogEventsInput{
				LogGroupName:  &logGroup,
				LogStreamName: &logStream,
				StartTime:     aws.Int64(startTime),
				NextToken:     nextToken,
			}

			logEvents, err := cwlClient.GetLogEvents(ctx, getLogsInput)
			if err != nil {
				fmt.Printf("Error fetching logs: %v\n", err)
				return
			}

			for _, event := range logEvents.Events {
				logChan <- LogEvent{
					Timestamp: *event.Timestamp,
					Message:   *event.Message,
				}
			}

			if !follow {
				break
			}

			// If following, wait before next poll
			if len(logEvents.Events) == 0 {
				time.Sleep(time.Second)
			}

			nextToken = logEvents.NextForwardToken
		}
	}()

	return logChan, nil
}
