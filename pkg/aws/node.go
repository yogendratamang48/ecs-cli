package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/yogendratamang48/ecs/pkg/types"
)

func (c *ECSClient) ListNodes(ctx context.Context) ([]*types.Node, error) {
	input := &ecs.ListContainerInstancesInput{
		Cluster: &c.Context.Cluster,
	}

	result, err := c.Client.ListContainerInstances(ctx, input)

	if err != nil {
		return nil, fmt.Errorf("no container instances found in cluster")
	}
	if len(result.ContainerInstanceArns) == 0 {
		return nil, err
	}

	var nodes []*types.Node
	describeInput := &ecs.DescribeContainerInstancesInput{
		Cluster:            &c.Context.Cluster,
		ContainerInstances: result.ContainerInstanceArns,
	}
	describeResult, err := c.DescribeContainerInstances(ctx, describeInput)
	if err != nil {
		return nil, err
	}

	for _, instance := range describeResult.ContainerInstances {
		node := &types.Node{
			InstanceID:        string(*instance.Ec2InstanceId),
			ContainerInstance: string(*instance.ContainerInstanceArn),
			CapacityProvider:  string(*instance.CapacityProviderName),
			Status:            string(*instance.Status),
			RegisteredAt:      *instance.RegisteredAt,
			RunningTasks:      instance.RunningTasksCount,
			PendingTasks:      instance.PendingTasksCount,
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}
