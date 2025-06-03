// cmd/port_forward.go
package cmd

import (
	"context"
	"fmt"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/yogendratamang48/ecs/pkg/aws"
)

func portForwardCmd() *cobra.Command {
	var container string

	cmd := &cobra.Command{
		Use:   "port-forward TASK_ID LOCAL_PORT:CONTAINER_PORT",
		Short: "Forward a local port to a port in a container",
		Long: `Forward a local port to a port in a container using AWS ECS execute-command.

The container name will be auto-detected if not specified with the -c flag.

Examples:
  # Forward local port 8080 to container port 80 (auto-detects container name)
  ecs port-forward 1234567890-abcd 8080:80

  # Forward local port 8080 to container port 80 in a specific container
  ecs port-forward 1234567890-abcd 8080:80 -c nginx`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			taskID := args[0]
			portMapping := args[1]

			// Parse port mapping (LOCAL_PORT:CONTAINER_PORT)
			re := regexp.MustCompile(`^(\d+):(\d+)$`)
			matches := re.FindStringSubmatch(portMapping)
			if len(matches) != 3 {
				return fmt.Errorf("invalid port mapping format: %s, expected LOCAL_PORT:CONTAINER_PORT", portMapping)
			}

			localPort, err := strconv.Atoi(matches[1])
			if err != nil {
				return fmt.Errorf("invalid local port: %s", matches[1])
			}

			containerPort, err := strconv.Atoi(matches[2])
			if err != nil {
				return fmt.Errorf("invalid container port: %s", matches[2])
			}

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

			// Start port forwarding
			fmt.Printf("Forwarding local port %d to container port %d in task %s...\n", localPort, containerPort, taskID)
			if err := client.PortForward(context.Background(), taskID, containerName, localPort, containerPort); err != nil {
				return fmt.Errorf("failed to forward port: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&container, "container", "c", "", "Container name (will be auto-detected if not specified)")

	return cmd
}
