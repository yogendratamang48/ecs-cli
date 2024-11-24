package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// These variables will be set by GoReleaser
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information of ecs cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built at: %s\n", date)
	},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ecs",
	Short:   "A CLI to interact with AWS ECS Cluster and Services",
	Long:    `A CLI tool similar to kubectl for managing AWS ECS clusters, services, and tasks.`,
	Version: version,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	// Run: getServices,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(configCmd())
	rootCmd.AddCommand(getCmd())
	rootCmd.AddCommand(describeCmd())
	rootCmd.AddCommand(deleteCmd())
}
