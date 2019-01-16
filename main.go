package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ic2hrmk/minish/scheduler"
)

func main() {
	executor := scheduler.NewExecutor()

	taskID, err := executor.Add("roman", 501*time.Millisecond, []byte{0, 1, 2})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("task deployed ->", taskID)

	err = executor.Cancel(taskID)
	fmt.Println("task manually canceled")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)
}
