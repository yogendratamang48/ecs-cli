// pkg/types/task_detail.go
package types

import (
	"time"
)

type TaskDetail struct {
	TaskID               string             `json:"taskId" yaml:"taskId"`
	TaskARN              string             `json:"taskArn" yaml:"taskArn"`
	ClusterARN           string             `json:"clusterArn" yaml:"clusterArn"`
	TaskDefinitionARN    string             `json:"taskDefinitionArn" yaml:"taskDefinitionArn"`
	ContainerInstanceARN string             `json:"containerInstanceArn,omitempty" yaml:"containerInstanceArn,omitempty"`
	Status               string             `json:"status" yaml:"status"`
	DesiredStatus        string             `json:"desiredStatus" yaml:"desiredStatus"`
	CPU                  string             `json:"cpu,omitempty" yaml:"cpu,omitempty"`
	Memory               string             `json:"memory,omitempty" yaml:"memory,omitempty"`
	CreatedAt            time.Time          `json:"createdAt" yaml:"createdAt"`
	StartedAt            time.Time          `json:"startedAt,omitempty" yaml:"startedAt,omitempty"`
	StoppedAt            time.Time          `json:"stoppedAt,omitempty" yaml:"stoppedAt,omitempty"`
	StoppedReason        string             `json:"stoppedReason,omitempty" yaml:"stoppedReason,omitempty"`
	Group                string             `json:"group" yaml:"group"`
	LaunchType           string             `json:"launchType" yaml:"launchType"`
	Containers           []ContainerDetail  `json:"containers" yaml:"containers"`
	NetworkInterfaces    []NetworkInterface `json:"networkInterfaces,omitempty" yaml:"networkInterfaces,omitempty"`
}

type ContainerDetail struct {
	Name            string        `json:"name" yaml:"name"`
	Image           string        `json:"image" yaml:"image"`
	Status          string        `json:"status" yaml:"status"`
	RuntimeID       string        `json:"runtimeId,omitempty" yaml:"runtimeId,omitempty"`
	ExitCode        *int32        `json:"exitCode,omitempty" yaml:"exitCode,omitempty"`
	CreatedAt       time.Time     `json:"createdAt" yaml:"createdAt"`
	StartedAt       time.Time     `json:"startedAt,omitempty" yaml:"startedAt,omitempty"`
	FinishedAt      time.Time     `json:"finishedAt,omitempty" yaml:"finishedAt,omitempty"`
	HealthStatus    string        `json:"healthStatus,omitempty" yaml:"healthStatus,omitempty"`
	NetworkBindings []PortBinding `json:"networkBindings,omitempty" yaml:"networkBindings,omitempty"`
}

type PortBinding struct {
	ContainerPort int32  `json:"containerPort" yaml:"containerPort"`
	HostPort      int32  `json:"hostPort" yaml:"hostPort"`
	Protocol      string `json:"protocol" yaml:"protocol"`
}

type NetworkInterface struct {
	AttachmentID string `json:"attachmentId" yaml:"attachmentId"`
	PrivateIPv4  string `json:"privateIpv4" yaml:"privateIpv4"`
	PublicIPv4   string `json:"publicIpv4,omitempty" yaml:"publicIpv4,omitempty"`
	SubnetID     string `json:"subnetId" yaml:"subnetId"`
}
