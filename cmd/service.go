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

func getServicesCmd() *cobra.Command {
	var outputFormat string

	cmd := &cobra.Command{
		Use:     "services",
		Aliases: []string{"svc", "svc", "service"},
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
