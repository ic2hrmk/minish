package scheduler

import (
	"time"
)

type Beholder interface {
	AddTask(ownerID OwnerIdentifier, ttl time.Duration, payload []byte) (TaskIdentifier, error)
	CancelTask(taskID TaskIdentifier) error

	AttachListener(EventListener) (EventListenerIdentifier, error)
	AttachNamedListener(EventListenerIdentifier, EventListener) error
	DeleteListener(EventListenerIdentifier) error
}
