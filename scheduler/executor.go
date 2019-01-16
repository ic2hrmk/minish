package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ic2hrmk/minish/shared/contract/scheduler"
)

type executor struct {
	tasks  sync.Map
	stream chan scheduler.TaskIdentifier
}

func NewExecutor() scheduler.Executor {
	return &executor{
		stream: make(chan scheduler.TaskIdentifier),
	}
}

type executionTask struct {
	meta  *scheduler.Task
	timer *time.Timer
}

func (e *executor) Add(
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
	executionTask := &executionTask{
		meta:  task,
		timer: time.NewTimer(task.TTL),
	}

	//
	// Start counting
	//
	go func() {
		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, task.TTL)

		select {

		case <-executionTask.timer.C:
			fmt.Println("TASK TIMER EXPIRED")

		case <-ctx.Done():
			fmt.Println("END BY DEADLINE")

		}

		fmt.Println("ROUTINE FINISHED")

		e.tasks.Delete(executionTask.meta.TaskID)
	}()

	e.tasks.Store(task.TaskID, executionTask)

	return task.TaskID, nil
}

func (e *executor) Cancel(taskID scheduler.TaskIdentifier) error {
	// Look for a task
	storedTask, taskExists := e.tasks.Load(taskID)
	if !taskExists {
		return nil
	}

	// Stop started task
	storedTask.(*executionTask).timer.Stop()

	// Remove it from task registry
	e.tasks.Delete(taskID)

	return nil
}
