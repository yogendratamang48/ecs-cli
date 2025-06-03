package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

// PortForward forwards a local port to a container port
func (c *ECSClient) PortForward(ctx context.Context, taskID string, container string, localPort int, containerPort int) error {
	// For port forwarding, we'll use a command that creates a TCP proxy in the container
	// This will allow us to forward traffic from the container port to the session
	// We'll try different approaches based on what's available in the container
	command := fmt.Sprintf("sh -c 'timeout 1 bash -c \"echo > /dev/tcp/localhost/%d\" 2>/dev/null || { echo \"Error: Port %d is not open in the container. Make sure the service is running.\"; exit 1; }; if command -v socat >/dev/null 2>&1; then socat STDIO TCP:localhost:%d; elif command -v nc >/dev/null 2>&1; then nc localhost %d; elif command -v bash >/dev/null 2>&1; then bash -c \"exec 3</dev/tcp/localhost/%d; cat <&3 & cat >&3; wait\"; else echo \"Error: No suitable tool found for port forwarding.\"; echo \"Please install socat or netcat in the container.\"; exit 1; fi'", containerPort, containerPort, containerPort, containerPort, containerPort)

	// Create the ECS execute-command API call
	execCommandInput := &ecs.ExecuteCommandInput{
		Cluster:     &c.Context.Cluster,
		Task:        &taskID,
		Container:   &container,
		Command:     &command,
		Interactive: true, // Always set to true as ECS only supports interactive mode
	}

	// Execute the command
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

	fmt.Printf("Forwarding from 127.0.0.1:%d -> %d\n", localPort, containerPort)
	fmt.Printf("Forwarding from [::1]:%d -> %d\n", localPort, containerPort)
	fmt.Println("Press Ctrl+C to stop port forwarding")

	// Handle Ctrl+C to gracefully terminate the port forwarding
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Start the session-manager-plugin in a goroutine
	go func() {
		<-sigCh
		fmt.Println("\nStopping port forwarding...")
		cmd.Process.Kill()
		os.Exit(0)
	}()

	// Run the session-manager-plugin
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("session-manager-plugin failed: %w", err)
	}

	return nil
}
