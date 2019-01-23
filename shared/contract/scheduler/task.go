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

func GenerateTaskIdentifier() TaskIdentifier {
	return TaskIdentifier(uuid.New().String())
}
