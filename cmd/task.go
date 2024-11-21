// cmd/task.go
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yogendratamang48/ecs/pkg/aws"
	"github.com/yogendratamang48/ecs/pkg/utils"
	"gopkg.in/yaml.v2"
)

func getTasksCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "List tasks",
		Long:  `Display all tasks in the current ECS cluster context.`,
		RunE: func(cmd *cobra.Command, args []string) error {
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

			// Get tasks
			tasks, err := client.ListTasks(context.Background())
			if err != nil {
				return fmt.Errorf("failed to list tasks: %w", err)
			}

			// Handle different output formats
			switch outputFormat {
			case "json":
				data, err := json.MarshalIndent(tasks, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal to JSON: %w", err)
				}
				fmt.Println(string(data))
				return nil

			case "yaml":
				data, err := yaml.Marshal(tasks)
				if err != nil {
					return fmt.Errorf("failed to marshal to YAML: %w", err)
				}
				fmt.Println(string(data))
				return nil

			case "":
				// Default table output
				headers := []string{
					"TASK ID",
					"STATUS",
					"TASK DEFINITION",
					"STARTED",
					"AGE",
				}

				table := utils.NewTableFormatter(headers)

				for _, task := range tasks {
					age := time.Since(task.CreatedAt).Round(time.Second)
					started := "-"
					if !task.StartedAt.IsZero() {
						started = formatAge(time.Since(task.StartedAt))
					}

					row := []string{
						task.TaskID,
						task.Status,
						task.TaskDefFamily,
						started,
						formatAge(age),
					}
					table.AppendRow(row)
				}

				table.Render()
				return nil

			default:
				return fmt.Errorf("unsupported output format: %s", outputFormat)
			}
		},
	}

	// Add flags
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (json|yaml)")

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
