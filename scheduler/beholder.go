package scheduler

import (
	"time"
)

type Beholder interface {
	AddTask(ownerID OwnerIdentifier, ttl time.Duration, payload []byte) (TaskIdentifier, error)
	CancelTask(taskID TaskIdentifier) error

	AttachListener(EventListenerMethod) (EventListenerIdentifier, error)
	AttachNamedListener(EventListenerIdentifier, EventListenerMethod) error
	DeleteListener(EventListenerIdentifier) error
}
