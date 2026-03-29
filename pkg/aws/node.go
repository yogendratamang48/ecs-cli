package aws

import (
	"context"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// ECSClient contains the AWS ECS client
type ECSClient struct {
	Client *ecs.Client
}

// ListNodes lists the container instances and sorts them by RegisteredAt
func (c *ECSClient) ListNodes(ctx context.Context) ([]*types.Node, error) {
	// Load the AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	
	// Create the ECS service client
	client := ecs.NewFromConfig(cfg)
	
	// List container instances
	result, err := client.ListContainerInstances(ctx, &ecs.ListContainerInstancesInput{
		Cluster: aws.String("your-cluster-name"), // Replace with your cluster name
	})
	if err != nil {
		return nil, err
	}
	
	// Describe container instances in batches
	var nodes []*types.Node
	for _, arn := range result.ContainerInstanceArns {
		describeResp, err := client.DescribeContainerInstances(ctx, &ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String("your-cluster-name"),
			ContainerInstances: []string{*arn},
		})
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, describeResp.ContainerInstances...)
	}

	// Sort nodes by RegisteredAt
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].RegisteredAt.Before(*nodes[j].RegisteredAt)
	})
	return nodes, nil
}