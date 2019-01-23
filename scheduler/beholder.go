package scheduler

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ic2hrmk/minish/shared/contract/scheduler"
)

type beholder struct {
	tasks               sync.Map
	listeners           sync.Map
	notificationTimeout time.Duration
}

const (
	defaultNotificationTimeout = 3 * time.Second
)

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
			fmt.Println("TASK TIMER EXPIRED")
			isCanceled = false

		case <-ctx.Done():
			fmt.Println("END BY DEADLINE")
			isCanceled = true

		}

		fmt.Println("ROUTINE FINISHED")

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

func (rcv *beholder) AttachListener(l scheduler.EventListener) (scheduler.EventListenerIdentifier, error) {
	// Create unique listener identifier
	listenerIdentifier := scheduler.GenerateEventListenerIdentifier()

	// Attach listener
	rcv.listeners.Store(listenerIdentifier, l)

	return listenerIdentifier, nil
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

func (rcv *beholder) notifyListeners(event scheduler.Event) {
	wg := sync.WaitGroup{}

	rcv.listeners.Range(func(identifier, listener interface{}) bool {
		wg.Add(1)

		go func() {
			defer wg.Done()


			ctx := context.Background()
			ctx, _ = context.WithTimeout(ctx, rcv.notificationTimeout)


			done := make(chan struct{})

			go func() {
				log.Printf("DISPATCH EVENT TO [%s] \n", identifier.(scheduler.EventListenerIdentifier))
				listener.(scheduler.EventListener).Listen(event)

				log.Printf("[%s] NOTIFICATION METHOD IS DONE\n", identifier.(scheduler.EventListenerIdentifier))
				done <- struct{}{}
			}()

			select {
			case <-done:
				log.Printf("LISTENER NOTIFICATION OK\n")

			case <-ctx.Done():
				log.Printf("NOTIFICATION METHOD IS TIMED OUT\n")

			}
		}()

		return true
	})

	wg.Wait()
	log.Println("NOTIFICATION ENDED")
}
