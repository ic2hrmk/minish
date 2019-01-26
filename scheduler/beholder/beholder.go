package beholder

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ic2hrmk/minish/scheduler"
)

type beholder struct {
	tasks               sync.Map
	listeners           sync.Map
	notificationTimeout time.Duration
}

const (
	defaultNotificationTimeout = 3 * time.Second
	defaultDebugMessagesShow   = false
)

var DEBUG = defaultDebugMessagesShow

func NewBeholder() scheduler.Beholder {
	return &beholder{
		notificationTimeout: defaultNotificationTimeout,
	}
}

type beholderTask struct {
	meta  *scheduler.Task
	timer *time.Timer
}

func (rcv *beholder) AddTask(
	ownerID scheduler.OwnerIdentifier,
	ttl time.Duration,
	payload []byte,
) (
	scheduler.TaskIdentifier,
	error,
) {
	//
	// Create task
	//
	task := &scheduler.Task{
		OwnerID: ownerID,
		TaskID:  scheduler.GenerateTaskIdentifier(),
		TTL:     ttl,
		Payload: payload,
	}

	//
	// We added execution task to have a control under task timer
	//
	executionTask := &beholderTask{
		meta:  task,
		timer: time.NewTimer(task.TTL),
	}

	//
	// Start counting
	//
	go func() {
		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, task.TTL)

		var isCanceled bool

		select {

		case <-executionTask.timer.C:
			isCanceled = false
			debugLogf("task [ID=%s] ended", task.TaskID)

		case <-ctx.Done():
			isCanceled = true
			debugLogf("task [ID=%s] canceled", task.TaskID)

		}

		rcv.notifyListeners(scheduler.Event{
			Task:       task,
			IsCanceled: isCanceled,
		})

		rcv.tasks.Delete(executionTask.meta.TaskID)
	}()

	rcv.tasks.Store(task.TaskID, executionTask)

	return task.TaskID, nil
}

func (rcv *beholder) CancelTask(taskID scheduler.TaskIdentifier) error {
	// Look for a task
	storedTask, taskExists := rcv.tasks.Load(taskID)
	if !taskExists {
		return nil
	}

	// Stop started task
	storedTask.(*beholderTask).timer.Stop()

	// Remove it from task registry
	rcv.tasks.Delete(taskID)

	return nil
}

func (rcv *beholder) AttachListener(l scheduler.EventListenerMethod) (scheduler.EventListenerIdentifier, error) {
	// Create unique listener identifier
	listenerIdentifier := scheduler.GenerateEventListenerIdentifier()

	// Attach listener
	rcv.listeners.Store(listenerIdentifier, l)

	return listenerIdentifier, nil
}

func (rcv *beholder) AttachNamedListener(
	identifier scheduler.EventListenerIdentifier, l scheduler.EventListenerMethod,
) (
	error,
) {
	if _, alreadyExists := rcv.listeners.Load(identifier); alreadyExists {
		return fmt.Errorf("event listener with same identifier [%s] is already attached", identifier)
	}

	// Attach listener
	rcv.listeners.Store(identifier, l)

	return nil
}

func (rcv *beholder) DeleteListener(listenerIdentifier scheduler.EventListenerIdentifier) error {
	// Look for a listener
	attachedListener, listenerExists := rcv.listeners.Load(listenerIdentifier)
	if !listenerExists {
		return nil
	}

	// Kill listener
	rcv.listeners.Delete(attachedListener)

	return nil
}

// Asynchronously notifies all attached listeners with deadline
func (rcv *beholder) notifyListeners(event scheduler.Event) {
	wg := sync.WaitGroup{}

	rcv.listeners.Range(func(identifier, listener interface{}) bool {
		wg.Add(1)

		go func() {
			defer wg.Done()

			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, rcv.notificationTimeout)

			done := make(chan struct{})

			go func() {
				debugLogf("dispatch event to listener  [ID=%s]\n", identifier.(scheduler.EventListenerIdentifier))
				listener.(scheduler.EventListenerMethod)(event)

				debugLogf("listener [ID=%s] ended listen\n", identifier.(scheduler.EventListenerIdentifier))
				done <- struct{}{}
			}()

			select {

			case <-done:
				// That's ok, nothing to do

			case <-ctx.Done():
				debugLogf("notification for listener [ID=%s] is timed out\n",
					identifier.(scheduler.EventListenerIdentifier))

			}

			cancel()
		}()

		return true
	})

	wg.Wait()

	debugLogf("notification ended")
}

func debugLogf(message string, data ...interface{}) {
	if DEBUG {
		log.Printf("[BEHOLDER] "+message, data...)
	}
}
