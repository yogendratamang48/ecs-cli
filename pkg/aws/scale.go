package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func (c *ECSClient) ScaleService(ctx context.Context, serviceName string, desiredCount int32) error {
	input := &ecs.UpdateServiceInput{
		Cluster:      &c.Context.Cluster,
		Service:      &serviceName,
		DesiredCount: &desiredCount,
	}
	_, err := c.Client.UpdateService(ctx, input)
	return err
}
