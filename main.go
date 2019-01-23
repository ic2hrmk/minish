package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ic2hrmk/minish/scheduler"
	subcontract "github.com/ic2hrmk/minish/shared/contract/scheduler"
)

func main() {
	//
	// Add beholder + some listeners
	//
	beholder := scheduler.NewBeholder()

	for i := 0; i < 10; i++ {
		listener := &TestBeholderConsumer{}

		id, err := beholder.AttachListener(listener)
		if err != nil {
			panic("attaching listener error" + err.Error())
		}

		listener.Identifier = id
	}

	taskID, err := beholder.AddTask("roman", 501*time.Millisecond, []byte{0, 1, 2})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("task deployed ->", taskID)

	err = beholder.CancelTask(taskID)
	fmt.Println("task manually canceled")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 30)
}

type TestBeholderConsumer struct {
	Identifier subcontract.EventListenerIdentifier
}

func (rcv *TestBeholderConsumer) Listen(event subcontract.Event) {
	fmt.Printf("[%s] event received: task [%s], canceled [%t]\n", rcv.Identifier, event.Task.TaskID, event.IsCanceled)
	sleepTime := time.Duration(rand.Int63n(10)) * time.Second

	fmt.Printf("[%s] will sleep for [%s] seconds\n", rcv.Identifier, sleepTime.String())
	time.Sleep(sleepTime)

	fmt.Printf("[%s] done!!!\n", rcv.Identifier,)
}



