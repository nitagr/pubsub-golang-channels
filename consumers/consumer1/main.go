package main

import (
	"fmt"

	"github.com/r3labs/sse/v2"
)

func main() {

	client := sse.NewClient("http://localhost:8080/events")

	client.Subscribe("messages", func(msg *sse.Event) {
		// Got some data!
		fmt.Printf("message recieved %s", msg.Data)
	})
	fmt.Print("Client")
}
