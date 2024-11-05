// pkg/types/service_detail.go
package types

import (
	"time"
)

type ServiceDetail struct {
	Name          string         `json:"name" yaml:"name"`
	Status        string         `json:"status" yaml:"status"`
	TaskDef       string         `json:"taskDefinition" yaml:"taskDefinition"`
	DesiredCount  int            `json:"desiredCount" yaml:"desiredCount"`
	RunningCount  int            `json:"runningCount" yaml:"runningCount"`
	PendingCount  int            `json:"pendingCount" yaml:"pendingCount"`
	CreatedAt     time.Time      `json:"createdAt" yaml:"createdAt"`
	LoadBalancers []LoadBalancer `json:"loadBalancers,omitempty" yaml:"loadBalancers,omitempty"`
	NetworkConfig NetworkConfig  `json:"networkConfig" yaml:"networkConfig"`
	Events        []ServiceEvent `json:"events" yaml:"events"`
}

type LoadBalancer struct {
	Type          string `json:"type" yaml:"type"`
	TargetGroup   string `json:"targetGroup" yaml:"targetGroup"`
	ContainerName string `json:"containerName" yaml:"containerName"`
	ContainerPort int    `json:"containerPort" yaml:"containerPort"`
}

type NetworkConfig struct {
	Type           string   `json:"type" yaml:"type"`
	SubnetIds      []string `json:"subnetIds,omitempty" yaml:"subnetIds,omitempty"`
	SecurityGroups []string `json:"securityGroups,omitempty" yaml:"securityGroups,omitempty"`
	PublicIP       string   `json:"publicIP,omitempty" yaml:"publicIP,omitempty"`
}

type ServiceEvent struct {
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	Message   string    `json:"message" yaml:"message"`
}
