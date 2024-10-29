// pkg/types/context.go
package types

// Context represents an ECS CLI configuration context
type Context struct {
	Name    string `mapstructure:"name" yaml:"name"`
	Cluster string `mapstructure:"cluster" yaml:"cluster"`
	Profile string `mapstructure:"profile" yaml:"profile"`
	Region  string `mapstructure:"region" yaml:"region"`
}
