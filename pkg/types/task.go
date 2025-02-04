// pkg/types/task.go
package types

import "time"

type Task struct {
	TaskId               string    `json:"taskId" yaml:"taskId"`
	TaskArn              string    `json:"taskArn" yaml:"taskArn"`
	Status               string    `json:"status" yaml:"status"`
	Cpu                  string    `json:"cpu" yaml:"cpu"`
	Memory               string    `json:"memory" yaml:"memory"`
	LaunchType           string    `json:"launchType" yaml:"launchType"`
	TaskDefFamily        string    `json:"taskDefinitionFamily" yaml:"taskDefinitionFamily"`
	LastStatus           string    `json:"lastStatus" yaml:"lastStatus"`
	DesiredStatus        string    `json:"desiredStatus" yaml:"desiredStatus"`
	CreatedAt            time.Time `json:"createdAt" yaml:"createdAt"`
	StartedAt            time.Time `json:"startedAt" yaml:"startedAt"`
	Group                string    `json:"group" yaml:"group"`
	ContainerInstanceArn string    `json:"containerInstanceArn" yaml:"containerInstanceArn"`
	CapacityProvider     string    `json:"capacityProvider" yaml:"capacityProvider"`
}
