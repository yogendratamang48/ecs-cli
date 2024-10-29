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
