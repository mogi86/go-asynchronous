package main

import (
	"fmt"
	"net/http"

	publisher "github.com/mogi86/go-asynchronous/internal/http"
)

func main() {
	mux := http.NewServeMux()

	fmt.Println("Server start...")

	// publish message to AWS SQS
	mux.Handle("/publish", http.HandlerFunc(publisher.PublishMessage))

	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		fmt.Printf("build server failed. %+v\n", err)
	}
}
