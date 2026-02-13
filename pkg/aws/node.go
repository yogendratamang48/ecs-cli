package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

// ECSNode represents the information about a single ECS container instance
 type ECSNode struct { 
	ClusterArn         string `json:"clusterArn"`
	ContainerInstanceArn string `json:"containerInstanceArn"`
	RunningTasksCount  int64  `json:"runningTasksCount"`
	PendingTasksCount  int64  `json:"pendingTasksCount"`
	EcsAgentVersion    string `json:"ecsAgentVersion"`
}

// ListECSNodes retrieves the ECS nodes from the specified cluster
func ListECSNodes(clusterName string) ([]ECSNode, error) {
	sess, err := session.NewSession() 
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	svc := ecs.New(sess)

	input := &ecs.ListContainerInstancesInput{
		Cluster: aws.String(clusterName),
	}

	result, err := svc.ListContainerInstances(input)
	if err != nil {
		return nil, fmt.Errorf("failed to list container instances: %w", err)
	}

	nodes := make([]ECSNode, 0)
	for _, instanceArn := range result.ContainerInstanceArns {
		descInput := &ecs.DescribeContainerInstancesInput{
			Cluster:            aws.String(clusterName),
			ContainerInstances: []*string{instanceArn},
		}
		descResult, err := svc.DescribeContainerInstances(descInput)
		if err != nil {
			return nil, fmt.Errorf("failed to describe container instance: %w", err)
		}

		for _, node := range descResult.ContainerInstances {
			nodes = append(nodes, ECSNode{
				ClusterArn:         clusterName,
				ContainerInstanceArn: *node.ContainerInstanceArn,
				RunningTasksCount:  *node.RunningTasksCount,
				PendingTasksCount:  *node.PendingTasksCount,
				EcsAgentVersion:    *node.Ec2InstanceId,
			})
		}
	}

	return nodes, nil
}