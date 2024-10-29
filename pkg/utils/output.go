// pkg/utils/output.go
package utils

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

// FormatOutput formats the data according to the specified format
func FormatOutput(data interface{}, format string) (string, error) {
	switch format {
	case "json":
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "", err
		}
		return string(b), nil
	case "yaml":
		b, err := yaml.Marshal(data)
		if err != nil {
			return "", err
		}
		return string(b), nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", format)
	}
}
