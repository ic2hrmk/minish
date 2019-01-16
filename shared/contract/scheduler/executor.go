package scheduler

import (
	"github.com/google/uuid"
	"time"
)

type (
	TaskIdentifier  string
	OwnerIdentifier string
)

type Task struct {
	TaskID  TaskIdentifier
	OwnerID OwnerIdentifier
	TTL     time.Duration
	Payload []byte
}

type Executor interface {
	Add(ownerID OwnerIdentifier, ttl time.Duration, payload []byte) (TaskIdentifier, error)
	Cancel(taskID TaskIdentifier) error
}

func GenerateTaskIdentifier() TaskIdentifier {
	return TaskIdentifier(uuid.New().String())
}
