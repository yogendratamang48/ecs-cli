// cmd/exec.go
package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yogendratamang48/ecs/pkg/aws"
)

func execCmd() *cobra.Command {
	var container string

	cmd := &cobra.Command{
		Use:   "exec TASK_ID [-- COMMAND]",
		Short: "Execute a command in a running container",
		Long: `Execute a command in a running container using AWS ECS execute-command.

The container name will be auto-detected if not specified with the -c flag.

Note: AWS ECS execute-command only supports interactive mode.

Examples:
  # Execute a shell in a container (auto-detects container name)
  ecs exec 1234567890-abcd -- /bin/sh

  # Execute a specific command in a container (auto-detects container name)
  ecs exec 1234567890-abcd -- ls -la

  # Execute a command in a specific container
  ecs exec 1234567890-abcd -c nginx -- ls -la`,
		Args: cobra.MinimumNArgs(1),
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

			// Parse command
			var command string
			var hasCommand bool

			if len(args) > 1 {
				// Skip the first argument (taskID) and any "--" separator
				cmdArgs := args[1:]
				if cmdArgs[0] == "--" {
					cmdArgs = cmdArgs[1:]
				}

				// Check if there are any command arguments after the separator
				if len(cmdArgs) > 0 {
					command = strings.Join(cmdArgs, " ")
					hasCommand = true
				}
			}

			// If no command is provided, return an error
			if !hasCommand {
				return fmt.Errorf("error: you must specify at least one command for the container")
			}

			// If container name is not specified, try to detect it
			containerName := container
			if containerName == "" {
				detectedContainer, err := client.GetContainerNameForTask(context.Background(), taskID)
				if err != nil {
					return fmt.Errorf("failed to detect container name: %w", err)
				}
				containerName = detectedContainer
				fmt.Printf("Auto-detected container: %s\n", containerName)
			}

			// Execute the command in interactive mode
			if err := client.ExecuteCommand(context.Background(), taskID, true, containerName, command); err != nil {
				return fmt.Errorf("failed to execute command: %w", err)
			}

			return nil
		},
	}

	// Add container flag
	cmd.Flags().StringVarP(&container, "container", "c", "", "Container name (will be auto-detected if not specified)")

	return cmd
}
