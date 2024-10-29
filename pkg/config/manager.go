// pkg/config/manager.go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/yogendratamang48/ecs/pkg/types"
)

type Manager struct {
	configFile string
}

// NewManager creates a new config manager
func NewManager() *Manager {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("could not get home directory: %v", err))
	}

	configDir := filepath.Join(home, ".ecs")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		panic(fmt.Sprintf("could not create config directory: %v", err))
	}

	configFile := filepath.Join(configDir, "config.yaml")

	// Initialize Viper configuration
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// Create config file if it doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		file, err := os.Create(configFile)
		if err != nil {
			panic(fmt.Sprintf("could not create config file: %v", err))
		}
		file.Close()

		viper.Set("contexts", map[string]interface{}{})
		viper.Set("current-context", "")

		if err := viper.WriteConfig(); err != nil {
			panic(fmt.Sprintf("could not write initial config: %v", err))
		}
	}

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Sprintf("could not read config file: %v", err))
	}

	return &Manager{
		configFile: configFile,
	}
}

// GetContext returns the current context
func (m *Manager) GetContext() (*types.Context, error) {
	currentContextName := viper.GetString("current-context")
	if currentContextName == "" {
		return nil, fmt.Errorf("no current context set")
	}

	contexts := viper.GetStringMap("contexts")
	ctxInterface, ok := contexts[currentContextName]
	if !ok {
		return nil, fmt.Errorf("current context '%s' not found", currentContextName)
	}

	var ctx types.Context
	if err := m.convertToStruct(ctxInterface, &ctx); err != nil {
		return nil, err
	}

	return &ctx, nil
}

// SetContext saves a new context and sets it as current
func (m *Manager) SetContext(ctx *types.Context) error {
	contexts := viper.GetStringMap("contexts")
	if contexts == nil {
		contexts = make(map[string]interface{})
	}

	contexts[ctx.Name] = ctx
	viper.Set("contexts", contexts)
	viper.Set("current-context", ctx.Name)

	return viper.WriteConfig()
}

// ListContexts returns all configured contexts and the current context name
func (m *Manager) ListContexts() ([]types.Context, string, error) {
	var contextList []types.Context
	currentContext := viper.GetString("current-context")

	contexts := viper.GetStringMap("contexts")
	for _, ctxInterface := range contexts {
		var ctx types.Context
		if err := m.convertToStruct(ctxInterface, &ctx); err != nil {
			return nil, "", err
		}
		contextList = append(contextList, ctx)
	}

	return contextList, currentContext, nil
}

// Helper method to convert map[string]interface{} to Context struct
func (m *Manager) convertToStruct(in interface{}, out *types.Context) error {
	config := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   out,
		TagName:  "mapstructure",
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(in)
}

// GetConfigFile returns the path to the config file
func (m *Manager) GetConfigFile() string {
	return m.configFile
}

// DeleteContext removes a context
func (m *Manager) DeleteContext(name string) error {
	contexts := viper.GetStringMap("contexts")
	if _, exists := contexts[name]; !exists {
		return fmt.Errorf("context '%s' not found", name)
	}

	delete(contexts, name)
	viper.Set("contexts", contexts)

	// If we're deleting the current context, clear it
	if viper.GetString("current-context") == name {
		viper.Set("current-context", "")
	}

	return viper.WriteConfig()
}

// UseContext sets the current context
func (m *Manager) UseContext(name string) error {
	contexts := viper.GetStringMap("contexts")
	if _, exists := contexts[name]; !exists {
		return fmt.Errorf("context '%s' not found", name)
	}

	viper.Set("current-context", name)
	return viper.WriteConfig()
}

// ValidateContext checks if a context is valid
func (m *Manager) ValidateContext(ctx *types.Context) error {
	if ctx.Name == "" {
		return fmt.Errorf("context name cannot be empty")
	}
	if ctx.Cluster == "" {
		return fmt.Errorf("cluster name cannot be empty")
	}
	return nil
}

func (m *Manager) ViewConfig() (string, error) {
	contexts := viper.GetStringMap("contexts")
	currentContext := viper.GetString("current-context")

	// Format the configuration as YAML
	output := fmt.Sprintf("Current-context: %s\n\nContexts:\n", currentContext)

	for name, ctxInterface := range contexts {
		var ctx types.Context
		if err := m.convertToStruct(ctxInterface, &ctx); err != nil {
			return "", err
		}

		output += fmt.Sprintf("\n%s:\n", name)
		output += fmt.Sprintf("  cluster: %s\n", ctx.Cluster)
		output += fmt.Sprintf("  profile: %s\n", ctx.Profile)
		output += fmt.Sprintf("  region: %s\n", ctx.Region)
	}

	return output, nil
}
