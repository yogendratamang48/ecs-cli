// cmd/scale.go
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yogendratamang48/ecs/pkg/aws"
)

func scaleCmd() *cobra.Command {
	var replicas int32

	cmd := &cobra.Command{
		Use:   "scale SERVICE_NAME",
		Short: "Scale a service",
		Long: `Scale a service by setting its desired count.

Example:
  # Scale a service to 2 replicas
  ecs scale my-service --replicas=2
  
  # Scale a service to 0 replicas (stop all tasks)
  ecs scale my-service --replicas=0`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			serviceName := args[0]

			// Get current context
			ctx, err := configManager.GetContext()
			if err != nil {
				return fmt.Errorf("failed to get current context: %w", err)
			}

			// Create ECS client
			client, err := aws.NewECSClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create ECS client: %w", err)
			}

			// Scale the service
			if err := client.ScaleService(context.Background(), serviceName, replicas); err != nil {
				return fmt.Errorf("failed to scale service: %w", err)
			}

			fmt.Printf("Successfully scaled service %s to %d replicas\n", serviceName, replicas)
			return nil
		},
	}

	cmd.Flags().Int32Var(&replicas, "replicas", 1, "Number of desired tasks")
	cmd.MarkFlagRequired("replicas")

	return cmd
}
