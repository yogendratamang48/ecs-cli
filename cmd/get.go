// cmd/service.go
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

// getCmd represents the get command group
func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or many resources",
		Long: `Display one or many resources from an ECS cluster.
        
Valid resource types are:
  * services (alias: svc, svcs)
  * tasks
  * nodes (alias: node, servers, instances)
        
Examples:
  # List all services in the current context
  ecs get services
  
  # List all services using the short alias
  ecs get svc
  
  # List all tasks in the current context
  ecs get tasks

  # List nodes in the cluster
  ecs get nodes`,
	}

	// Add subcommands to 'get'
	cmd.AddCommand(getServicesCmd()) // get services
	cmd.AddCommand(getTasksCmd())    // get tasks
	cmd.AddCommand(getNodesCmd())    // get nodes

	return cmd
}

func getServicesCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:     "services",
		Aliases: []string{"svc", "svcs", "service", "services"},
		Short:   "List services",
		Long:    `Display all services in the current ECS cluster context.`,
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

			// Get services
			services, err := client.ListServices(context.Background())
			if err != nil {
				return fmt.Errorf("failed to list services: %w", err)
			}
			switch outputFormat {
			case "json":
				data, err := json.MarshalIndent(services, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal to JSON: %w", err)
				}
				fmt.Println(string(data))
				return nil
			case "yaml":
				data, err := yaml.Marshal(services)
				if err != nil {
					return fmt.Errorf("failed to marshal to YAML: %w", err)
				}
				fmt.Println(string(data))
				return nil
			case "":
				// Display services
				headers := []string{
					"NAME",
					"STATUS",
					"DESIRED",
					"RUNNING",
					"PENDING",
					"AGE",
				}
				table := utils.NewTableFormatter(headers)

				for _, svc := range services {
					age := time.Since(svc.CreatedAt).Round(time.Second)

					row := []string{
						svc.Name,
						svc.Status,
						fmt.Sprintf("%d", svc.DesiredCount),
						fmt.Sprintf("%d", svc.RunningCount),
						fmt.Sprintf("%d", svc.PendingCount),
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
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (json|yaml|wide)")

	return cmd
}

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
			case "wide":
				headers := []string{
					"TASK ID",
					"STATUS",
					"TASK DEFINITION",
					"STARTED",
					"AGE",
					"CPU",
					"MEMORY",
					"LAUNCH TYPE",
					"CAPACITY PROVIDER",
				}
				table := utils.NewTableFormatter(headers)
				for _, task := range tasks {
					age := time.Since(task.CreatedAt).Round(time.Second)
					started := "-"
					if !task.StartedAt.IsZero() {
						started = formatAge(time.Since(task.StartedAt))
					}

					row := []string{
						task.TaskId,
						task.Status,
						task.TaskDefFamily,
						started,
						formatAge(age),
						task.Cpu,
						task.Memory,
						task.LaunchType,
						task.CapacityProvider,
					}
					table.AppendRow(row)
				}

				table.Render()
				return nil
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
						task.TaskId,
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
	cmd.Flags().StringVarP(&outputFormat, "output", "o", "", "Output format (json|yaml|wide)")

	return cmd
}

func getNodesCmd() *cobra.Command {
	var outputFormat string
	cmd := &cobra.Command{
		Use:     "nodes",
		Aliases: []string{"node", "nodes", "servers", "instances"},
		Short:   "List nodes",
		Long:    `Display all nodes in the current ECS cluster context.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := configManager.GetContext()
			if err != nil {
				return fmt.Errorf("failed to get current context: %w", err)
			}
			// create ecs client
			client, err := aws.NewECSClient(ctx)
			if err != nil {
				return fmt.Errorf("failed to create ECS client: %w", err)
			}
			// get nodes
			nodes, err := client.ListNodes(context.Background())
			if err != nil {
				return fmt.Errorf("failed to list nodes: %w", err)
			}
			switch outputFormat {
			case "json":
				data, err := json.MarshalIndent(nodes, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal to JSON: %w", err)
				}
				fmt.Println(string(data))
				return nil
			case "yaml":
				data, err := yaml.Marshal(nodes)
				if err != nil {
					return fmt.Errorf("failed to marshal to YAML: %w", err)
				}
				fmt.Println(string(data))
				return nil
			case "":
				headers := []string{
					"INSTANCE ID",
					"CAPACITY PROVIDER",
					"STATUS",
					"RUNNING TASKS",
					"PENDING TASKS",
					"AGE",
				}
				table := utils.NewTableFormatter(headers)
				for _, node := range nodes {
					age := time.Since(node.RegisteredAt).Round(time.Second)
					row := []string{
						node.InstanceID,
						node.CapacityProvider,
						node.Status,
						fmt.Sprintf("%d", node.RunningTasks),
						fmt.Sprintf("%d", node.PendingTasks),
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

// formatAge returns a human-readable string of the age
func formatAge(age time.Duration) string {
	if age < time.Hour {
		return fmt.Sprintf("%dm", int(age.Minutes()))
	} else if age < 24*time.Hour {
		return fmt.Sprintf("%dh", int(age.Hours()))
	} else if age < 30*24*time.Hour {
		return fmt.Sprintf("%dd", int(age.Hours()/24))
	} else if age < 365*24*time.Hour {
		return fmt.Sprintf("%dM", int(age.Hours()/(24*30)))
	}
	return fmt.Sprintf("%dy", int(age.Hours()/(24*365)))
}
