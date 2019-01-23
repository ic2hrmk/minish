package scheduler

import "github.com/google/uuid"

type Event struct {
	Task       *Task
	IsCanceled bool
}

type EventListenerIdentifier string

type EventListener interface {
	Listen(event Event)
}

func GenerateEventListenerIdentifier() EventListenerIdentifier {
	return EventListenerIdentifier(uuid.New().String())
}
