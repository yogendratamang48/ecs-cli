// cmd/delete.go
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yogendratamang48/ecs/pkg/aws"
)

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
		Long: `Delete resources from the ECS cluster.
        
Valid resource types are:
  * task TASK_ID    Delete (stop) a specific task
        
Examples:
  # Delete a specific task
  ecs delete task 1234567890-abcd-efgh-ijkl`,
	}

	// Add delete subcommands
	cmd.AddCommand(deleteTaskCmd())
	// Future commands can be added here, like:
	// cmd.AddCommand(deleteServiceCmd())

	return cmd
}

func deleteTaskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "task TASK_ID",
		Short: "Delete (stop) a task",
		Long:  `Delete (stop) a specific task from the ECS cluster.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskId := args[0]

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

			// Stop the task
			if err := client.StopTask(context.Background(), taskId); err != nil {
				return fmt.Errorf("failed to stop task: %w", err)
			}

			fmt.Printf("Task %s stopped successfully\n", taskId)
			return nil
		},
	}
}
