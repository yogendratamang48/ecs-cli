// cmd/describe.go
// cmd/describe.go
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/yogendratamang48/ecs/pkg/aws"
	"gopkg.in/yaml.v2"
)

func describeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Show details of a specific resource",
		Long: `Show detailed information about a specific resource.

Valid resource types are:
  * services [SERVICE_NAME]    Show details of one or all services
  * tasks [TASK_ID]           Show details of one or all tasks (not implemented yet)

Examples:
  # Describe all services
  ecs describe services

  # Describe a specific service
  ecs describe services my-service`,
	}

	cmd.AddCommand(describeServicesCmd())
	// Future: cmd.AddCommand(describeTasksCmd())

	return cmd
}

func describeServicesCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:   "services [SERVICE_NAME]",
		Short: "Show details of services",
		Long: `Show detailed information about one or all services in the cluster.

Examples:
  # Describe all services
  ecs describe services

  # Describe a specific service
  ecs describe services my-service

  # Output in JSON format
  ecs describe services my-service -o json`,

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

			var serviceNames []string
			if len(args) > 0 {
				serviceNames = []string{args[0]}
			} else {
				// If no service name provided, get all services
				services, err := client.ListServices(context.Background())
				if err != nil {
					return fmt.Errorf("failed to list services: %w", err)
				}
				for _, svc := range services {
					serviceNames = append(serviceNames, svc.Name)
				}
			}

			// Get detailed service information
			services, err := client.DescribeServices(context.Background(), serviceNames)
			if err != nil {
				return fmt.Errorf("failed to describe services: %w", err)
			}

			// Handle different output formats
			switch outputFormat {
			case "json":
				data, err := json.MarshalIndent(services, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal to JSON: %w", err)
				}
				fmt.Println(string(data))

			case "yaml":
				data, err := yaml.Marshal(services)
				if err != nil {
					return fmt.Errorf("failed to marshal to YAML: %w", err)
				}
				fmt.Println(string(data))

			default:
				// Default formatted output
				for _, svc := range services {
					fmt.Printf("Name:           %s\n", svc.Name)
					fmt.Printf("Status:         %s\n", svc.Status)
					fmt.Printf("Task Definition: %s\n", svc.TaskDef)
					fmt.Printf("Desired Count:  %d\n", svc.DesiredCount)
					fmt.Printf("Running Count:  %d\n", svc.RunningCount)
					fmt.Printf("Pending Count:  %d\n", svc.PendingCount)
					fmt.Printf("Created At:     %s\n", svc.CreatedAt.Format(time.RFC3339))

					if len(svc.LoadBalancers) > 0 {
						fmt.Println("\nLoad Balancers:")
						for _, lb := range svc.LoadBalancers {
							fmt.Printf("  - Target Group:    %s\n", lb.TargetGroup)
							fmt.Printf("    Container Name:  %s\n", lb.ContainerName)
							fmt.Printf("    Container Port:  %d\n", lb.ContainerPort)
						}
					}

					if svc.NetworkConfig.Type != "" {
						fmt.Println("\nNetwork Configuration:")
						fmt.Printf("  Type:            %s\n", svc.NetworkConfig.Type)
						if len(svc.NetworkConfig.SubnetIds) > 0 {
							fmt.Printf("  Subnets:         %v\n", svc.NetworkConfig.SubnetIds)
						}
						if len(svc.NetworkConfig.SecurityGroups) > 0 {
							fmt.Printf("  Security Groups: %v\n", svc.NetworkConfig.SecurityGroups)
						}
						if svc.NetworkConfig.PublicIP != "" {
							fmt.Printf("  Public IP:       %s\n", svc.NetworkConfig.PublicIP)
						}
					}

					if len(svc.Events) > 0 {
						fmt.Println("\nRecent Events:")
						for _, event := range svc.Events[:5] { // Show only last 5 events
							fmt.Printf("  %s: %s\n",
								event.CreatedAt.Format(time.RFC3339),
								event.Message)
						}
					}

					fmt.Println()
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (json|yaml)")

	return cmd
}
