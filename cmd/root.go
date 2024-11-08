package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ecs",
	Short: "A CLI to interact with AWS ECS Cluster and Services",
	Long:  `A CLI tool similar to kubectl for managing AWS ECS clusters, services, and tasks.`,
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ecs.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(configCmd())
	rootCmd.AddCommand(getCmd())
	rootCmd.AddCommand(describeCmd())
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

// getCmd represents the get command group
func getCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or many resources",
		Long: `Display one or many resources from an ECS cluster.
        
Valid resource types are:
  * services (alias: svc)
  * tasks
        
Examples:
  # List all services in the current context
  ecs get services
  
  # List all services using the short alias
  ecs get svc
  
  # List all tasks in the current context
  ecs get tasks`,
	}

	// Add subcommands to 'get'
	cmd.AddCommand(getServicesCmd()) // This adds the services command
	cmd.AddCommand(getTasksCmd())    // This adds the tasks command

	return cmd
}

func deleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete resources",
		Long: `Delete resources from the ECS cluster.
        
Valid resource types are:
  * task TASK_ID    Delete (stop) a specific task
  * service NAME    Delete a service (not implemented yet)
        
Examples:
  # Delete a specific task
  ecs delete task 1234567890-abcd-efgh-ijkl
  
  # Delete a service (not implemented yet)
  ecs delete service my-service`,
	}

	// Add delete subcommands
	cmd.AddCommand(deleteTaskCmd())
	// Future commands can be added here, like:
	// cmd.AddCommand(deleteServiceCmd())

	return cmd
}
