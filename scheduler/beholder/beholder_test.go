package beholder

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ic2hrmk/minish/scheduler"
)

func init(){
	DEBUG = true
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

func TestBeholder_Positive(t *testing.T) {

	b := NewBeholder()
	var err error

	listener := &TestBeholderConsumer{}

	if listener.Identifier, err = b.AttachListener(listener.Listen); err != nil {
		t.Fatal(err)
	}

	taskID, err := b.AddTask("roman", 200*time.Millisecond, []byte{0, 1, 2})
	t.Log("task deployed ->", taskID)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)
}

func TestBeholderWithNamedListeners_Positive(t *testing.T) {

	b := NewBeholder()

	namedListener := &TestBeholderConsumer{
		Identifier: "named",
	}

	if err := b.AttachNamedListener(namedListener.Identifier, namedListener.Listen); err != nil {
		t.Fatal(err)
	}

	taskID, err := b.AddTask("roman", 200*time.Millisecond, []byte{0, 1, 2})
	t.Log("task deployed ->", taskID)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)
}

func TestBeholderWithNamedListeners_PositiveWithTaskCancel(t *testing.T) {

	b := NewBeholder()

	namedListener := &TestBeholderConsumer{
		Identifier: "named",
	}

	if err := b.AttachNamedListener(namedListener.Identifier, namedListener.Listen); err != nil {
		t.Fatal(err)
	}

	taskID, err := b.AddTask("roman", 200*time.Millisecond, []byte{0, 1, 2})
	t.Log("task deployed ->", taskID)
	if err != nil {
		t.Fatal(err)
	}

	err = b.CancelTask(taskID)
	t.Log("task manually canceled")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Second)
}

func TestBeholderWithNamedListeners_NegativeAlreadyExists_(t *testing.T) {

	b := NewBeholder()

	namedListenerA := &TestBeholderConsumer{
		Identifier: "named",
	}

	namedListenerB := &TestBeholderConsumer{
		Identifier: "named",
	}

	if err := b.AttachNamedListener(namedListenerA.Identifier, namedListenerA.Listen); err != nil {
		t.Fatal(err)
	}

	if err := b.AttachNamedListener(namedListenerB.Identifier, namedListenerB.Listen); err == nil {
		t.Fatal(err)
	}
}





