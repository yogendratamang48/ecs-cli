// ListServices returns all services in the cluster
package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/yogendratamang48/ecs/pkg/types"
)

func (c *ECSClient) ListServices(ctx context.Context) ([]*types.Service, error) {
	var services []*types.Service
	var nextToken *string

	for {
		// List service ARNs
		input := &ecs.ListServicesInput{
			Cluster:    &c.Context.Cluster,
			NextToken:  nextToken,
			MaxResults: aws.Int32(10),
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
