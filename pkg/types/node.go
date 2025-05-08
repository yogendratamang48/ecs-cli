package types

import "time"

type Node struct {
	InstanceID       string    `json:"instance_id" yaml:"instance_id"`
	CapacityProvider string    `json:"capacity_provider" yaml:"capacity_provider"`
	RunningTasks     int32     `json:"running_tasks" yaml:"running_tasks"`
	PendingTasks     int32     `json:"pending_tasks" yaml:"pending_tasks"`
	Status           string    `json:"status" yaml:"status"`
	StatusReason     string    `json:"status_reason" yaml:"status_reason"`
	RegisteredAt     time.Time `json:"registered_at" yaml:"registered_at"`
}
