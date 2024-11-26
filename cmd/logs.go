// cmd/logs.go
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yogendratamang48/ecs/pkg/aws"
)

func logsCmd() *cobra.Command {
	var (
		follow    bool
		since     time.Duration
		container string
	)

	cmd := &cobra.Command{
		Use:   "logs TASK_ID",
		Short: "View logs from a task",
		Long: `View CloudWatch logs from a task's containers.

Examples:
  # View logs from a task
  ecs logs 1234567890-abcd

  # Follow logs from a task
  ecs logs 1234567890-abcd -f

  # View logs from the last 1 hour
  ecs logs 1234567890-abcd --since=1h

  # View logs from a specific container
  ecs logs 1234567890-abcd --container=nginx`,

		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]

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

			// Get logs
			logChan, err := client.GetTaskLogs(context.Background(), taskID, follow, since, container)
			if err != nil {
				return fmt.Errorf("failed to get logs: %w", err)
			}

			// Print logs as they come
			for event := range logChan {
				timestamp := time.Unix(event.Timestamp/1000, 0).Format(time.RFC3339)
				fmt.Printf("%s %s\n", timestamp, event.Message)
			}

			return nil
		},
	}

	// Add flags
	cmd.Flags().BoolVarP(&follow, "follow", "f", false, "Follow log output")
	cmd.Flags().DurationVar(&since, "since", 10*time.Minute, "Show logs since duration (e.g., 5m, 1h)")
	cmd.Flags().StringVar(&container, "container", "", "Show logs from a specific container")

	return cmd
}
