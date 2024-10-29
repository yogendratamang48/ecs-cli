// pkg/types/service.go
package types

import (
	"time"
)

type Service struct {
	Name         string    `json:"name" yaml:"name"`
	Status       string    `json:"status" yaml:"status"`
	TaskDef      string    `json:"taskDefinition" yaml:"taskDefinition"`
	DesiredCount int       `json:"desiredCount" yaml:"desiredCount"`
	RunningCount int       `json:"runningCount" yaml:"runningCount"`
	PendingCount int       `json:"pendingCount" yaml:"pendingCount"`
	CreatedAt    time.Time `json:"createdAt" yaml:"createdAt"`
}
