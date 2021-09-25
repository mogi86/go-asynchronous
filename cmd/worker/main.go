package main

import (
	"fmt"
	"github.com/mogi86/go-asynchronous/internal/worker"
)

var QueueName string

func main() {
	err := worker.Worker(QueueName, worker.GetMessage)
	if err != nil {
		fmt.Printf("somthing happened")
	}
}
