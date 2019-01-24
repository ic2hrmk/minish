package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ic2hrmk/minish/scheduler"
	"github.com/ic2hrmk/minish/scheduler/beholder"
)

func main() {
	//
	// Add b + some listeners
	//
	b := beholder.NewBeholder()

	beholder.DEBUG = true

	listenerNamed := &TestBeholderConsumer{
		Identifier: "named",
	}

	if err := b.AttachNamedListener(listenerNamed.Identifier, listenerNamed); err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		listener := &TestBeholderConsumer{}

		id, err := b.AttachListener(listener)
		if err != nil {
			panic("attaching listener error" + err.Error())
		}

		listener.Identifier = id
	}

	taskID, err := b.AddTask("roman", 2000*time.Millisecond, []byte{0, 1, 2})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("task deployed ->", taskID)

	err = b.CancelTask(taskID)
	fmt.Println("task manually canceled")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 30)
}

type TestBeholderConsumer struct {
	Identifier scheduler.EventListenerIdentifier
}

func (rcv *TestBeholderConsumer) Listen(event scheduler.Event) {
	fmt.Printf("[%s] event received: task [%s], canceled [%t]\n", rcv.Identifier, event.Task.TaskID, event.IsCanceled)
	sleepTime := time.Duration(rand.Int63n(10)) * time.Second

	fmt.Printf("[%s] will sleep for [%s] seconds\n", rcv.Identifier, sleepTime.String())
	time.Sleep(sleepTime)

	fmt.Printf("[%s] done!!!\n", rcv.Identifier,)
}



