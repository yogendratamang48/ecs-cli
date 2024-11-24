// cmd/config.go
package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yogendratamang48/ecs/pkg/config"
	"github.com/yogendratamang48/ecs/pkg/types"
)

var configManager *config.Manager

func init() {
	configManager = config.NewManager()
}
func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configDir := filepath.Join(home, ".ecs")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		panic(err)
	}
	configFile := filepath.Join(configDir, "config.yaml")
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// If the config file doesn't exist, create it with default settings
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create empty config file
		file, err := os.Create(configFile)
		if err != nil {
			panic(err)
		}
		file.Close()

		// Initialize with empty contexts
		viper.Set("contexts", map[string]interface{}{})
		viper.Set("current-context", "")

		// Write initial config
		if err := viper.WriteConfig(); err != nil {
			panic(err)
		}
	}

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Modify ecs cli configuration",
		Long:  `Modify ecs cli configuration including contexts, credentials, and other settings.`,
	}

	cmd.AddCommand(configSetContextCmd())
	cmd.AddCommand(configGetContextsCmd())
	cmd.AddCommand(configUseContextCmd())
	cmd.AddCommand(configDeleteContextCmd())
	cmd.AddCommand(configCurrentContextCmd())
	cmd.AddCommand(configViewCmd())

	return cmd
}

func configSetContextCmd() *cobra.Command {
	var ctx types.Context

	cmd := &cobra.Command{
		Use:   "set-context NAME",
		Short: "Set a new context for ECS CLI",
		Long: `Set a new context for ECS CLI with a given NAME (alias).

Example:
  # Set a context named "prod" for production cluster
  ecs config set-context prod --cluster production-cluster --profile prod-profile --region us-west-2`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx.Name = args[0]

			if err := configManager.SetContext(&ctx); err != nil {
				return fmt.Errorf("failed to save context: %w", err)
			}

			fmt.Printf("Context '%s' created and set as current context\n", ctx.Name)
			fmt.Printf("Cluster: %s\n", ctx.Cluster)
			fmt.Printf("Profile: %s\n", ctx.Profile)
			fmt.Printf("Region: %s\n", ctx.Region)

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&ctx.Cluster, "cluster", "", "ECS cluster name")
	flags.StringVar(&ctx.Profile, "profile", "default", "AWS profile name")
	flags.StringVar(&ctx.Region, "region", "us-east-1", "AWS region")
	cmd.MarkFlagRequired("cluster")

	return cmd
}

func configGetContextsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get-contexts",
		Short: "Display all configured contexts",
		Long:  "Display all configured contexts and indicate which one is currently active",
		RunE: func(cmd *cobra.Command, args []string) error {
			contexts, currentContext, err := configManager.ListContexts()
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
			if err != nil {
				return fmt.Errorf("failed to list contexts: %w", err)
			}
			printContextHeaders(w, false)
			for _, ctx := range contexts {
				printContext(ctx.Name, &ctx, w, false, ctx.Name == currentContext)
			}
			w.Flush()
			return nil
		},
	}
}

func configUseContextCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "use-context NAME",
		Short: "Set the current context",
		Long: `Set the current context to use for ECS commands.

Example:
  ecs config use-context prod`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := configManager.UseContext(args[0]); err != nil {
				return err
			}
			fmt.Printf("Switched to context %q\n", args[0])
			return nil
		},
	}
}

func configDeleteContextCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete-context NAME",
		Short: "Delete a context from the configuration",
		Long: `Delete a specified context from the ECS CLI configuration.

Example:
  ecs config delete-context old-context`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := configManager.DeleteContext(args[0]); err != nil {
				return err
			}
			fmt.Printf("Context %q deleted\n", args[0])
			return nil
		},
	}
}

func configCurrentContextCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "current-context",
		Short: "Display the current context",
		Long:  "Display the name and details of the current context in use",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, err := configManager.GetContext()
			if err != nil {
				return err
			}

			fmt.Printf("Current context: %s\n", ctx.Name)
			fmt.Printf("Cluster: %s\n", ctx.Cluster)
			fmt.Printf("Profile: %s\n", ctx.Profile)
			fmt.Printf("Region: %s\n", ctx.Region)
			return nil
		},
	}
}

func configViewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Display merged configuration settings",
		Long:  "Display the current state of the merged configuration settings",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := configManager.ViewConfig()
			if err != nil {
				return err
			}
			fmt.Println(config)
			return nil
		},
	}
}

func printContextHeaders(out io.Writer, nameOnly bool) error {
	columnNames := []string{"CURRENT", "NAME", "CLUSTER", "PROFILE", "REGION"}
	if nameOnly {
		columnNames = columnNames[:1]
	}
	_, err := fmt.Fprintf(out, "%s\n", strings.Join(columnNames, "\t"))
	return err
}

func printContext(name string, ctx *types.Context, w io.Writer, nameOnly bool, current bool) error {
	if nameOnly {
		_, err := fmt.Fprintf(w, "%s\n", name)
		return err
	}
	prefix := " "
	if current {
		prefix = "*"
	}
	_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", prefix, name, ctx.Cluster, ctx.Profile, ctx.Region)
	return err

}
