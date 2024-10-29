// pkg/types/task.go
package types

import "time"

type Task struct {
	TaskID               string    `json:"taskId" yaml:"taskId"`
	TaskARN              string    `json:"taskArn" yaml:"taskArn"`
	Status               string    `json:"status" yaml:"status"`
	TaskDefFamily        string    `json:"taskDefinitionFamily" yaml:"taskDefinitionFamily"`
	LastStatus           string    `json:"lastStatus" yaml:"lastStatus"`
	DesiredStatus        string    `json:"desiredStatus" yaml:"desiredStatus"`
	CreatedAt            time.Time `json:"createdAt" yaml:"createdAt"`
	StartedAt            time.Time `json:"startedAt" yaml:"startedAt"`
	Group                string    `json:"group" yaml:"group"`
	ContainerInstanceARN string    `json:"containerInstanceArn" yaml:"containerInstanceArn"`
}
