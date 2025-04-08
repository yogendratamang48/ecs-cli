package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

// ExecuteCommand runs a command on a container in a task
func (c *ECSClient) ExecuteCommand(ctx context.Context, taskID string, interactive bool, container string, command string) error {
	// Note: AWS ECS execute-command API only supports interactive mode
	// So we'll ignore the interactive parameter and always use interactive mode
	
	// Create the ECS execute-command API call
	execCommandInput := &ecs.ExecuteCommandInput{
		Cluster:     &c.Context.Cluster,
		Task:        &taskID,
		Container:   &container,
		Command:     &command,
		Interactive: true, // Always set to true as ECS only supports interactive mode
	}

	// Execute the command
	fmt.Printf("Starting session with task %s...\n", taskID)
	execCommandResult, err := c.Client.ExecuteCommand(ctx, execCommandInput)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	// Get the session details
	sessionId := *execCommandResult.Session.SessionId
	streamUrl := *execCommandResult.Session.StreamUrl
	tokenValue := *execCommandResult.Session.TokenValue

	// Prepare the input for the session-manager-plugin
	// The plugin expects a JSON string with the session details
	sessionInput := struct {
		SessionId    string `json:"sessionId"`
		StreamUrl    string `json:"streamUrl"`
		TokenValue   string `json:"tokenValue"`
		ClientMode   string `json:"clientMode"`
		ResponseMode string `json:"responseMode"`
	}{
		SessionId:    sessionId,
		StreamUrl:    streamUrl,
		TokenValue:   tokenValue,
		ClientMode:   "interactive",
		ResponseMode: "json",
	}

	fmt.Println("Starting interactive session...")

	sessionInputJSON, err := json.Marshal(sessionInput)
	if err != nil {
		return fmt.Errorf("failed to marshal session input: %w", err)
	}

	// Find the session-manager-plugin executable
	pluginPath, err := exec.LookPath("session-manager-plugin")
	if err != nil {
		return fmt.Errorf("session-manager-plugin not found: %w\nPlease install the Session Manager plugin: https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html", err)
	}

	// Start the session-manager-plugin process
	cmd := exec.CommandContext(ctx, pluginPath,
		string(sessionInputJSON),
		c.Context.Region,
		"StartSession")

	// Connect stdin/stdout/stderr to the current process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the plugin and wait for it to complete
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("session-manager-plugin failed: %w", err)
	}

	return nil
}
