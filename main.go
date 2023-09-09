package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/nitagr/pubsub2/topic"
)

var wg1 sync.WaitGroup
var wg2 sync.WaitGroup
var wg3 sync.WaitGroup

var mu1 sync.Mutex
var mu2 sync.Mutex

var totalMessagesSent int = 0
var totalMessagesAcknowledged int = 0

type Publisher struct {
	subscribers []chan interface{}
}

func NewPublisher(topicName string) *Publisher {
	return &Publisher{
		subscribers: topic.TopicMap[topicName],
	}
}

func (p *Publisher) Publish(message interface{}) {

	for _, subscriber := range p.subscribers {
		go func(sub chan interface{}) {
			mu1.Lock()
			sub <- message
			totalMessagesSent++
			mu1.Unlock()
		}(subscriber)

	}
}

func recieveMessagesOnChannels(sub chan interface{}, subscriberType string, ack chan interface{}) {
	go func() {
		for msg := range sub {
			fmt.Print(subscriberType)

			switch v := msg.(type) {
			case int:
				//fmt.Printf("Received an int: %d\n", v)
			case string:
				//fmt.Printf("Received a string: %s\n", v)
			case float64:
				//fmt.Printf("Received a float64: %f\n", v)
			default:
				fmt.Printf("Received an unknown type: %v\n", v)
			}

			data := fmt.Sprintf(" ack %s", msg)
			ack <- data
		}
	}()

}

func recieveAckOnChannels(ackChannel chan interface{}, wg *sync.WaitGroup) {
	// var wg sync.WaitGroup
	// wg.Add(1)
	go func() {
		// defer wg.Done()
		for ack := range ackChannel {
			mu2.Lock()
			// defer wg.Done()
			totalMessagesAcknowledged++
			fmt.Print(ack)
			mu2.Unlock()
		}
	}()

}

func main() {

	subs1 := topic.AddSubsriberToTopic(topic.TOPIC_CHANNEL_1)
	subs2 := topic.AddSubsriberToTopic(topic.TOPIC_CHANNEL_2)

	ackChan1 := make(chan interface{})
	ackChan2 := make(chan interface{})

	publisher1 := NewPublisher(topic.TOPIC_CHANNEL_1)
	publisher2 := NewPublisher(topic.TOPIC_CHANNEL_2)

	recieveMessagesOnChannels(subs1, "subscriber 1", ackChan1)
	recieveMessagesOnChannels(subs2, "subscriber 2", ackChan2)

	recieveAckOnChannels(ackChan1, &wg1)
	recieveAckOnChannels(ackChan2, &wg2)

	timeoutChannel := time.After(time.Second * 1)
	// wg1.Add(1)

	for {
		select {
		case <-timeoutChannel:
			fmt.Println("------")
			fmt.Println("totalMessagesSent %d", totalMessagesSent)
			fmt.Println("totalMessagesAcknowledged %d", totalMessagesAcknowledged)
			// wg1.Wait()

			return
		default:
			publisher1.Publish("Hello, World! publisher 1")
			publisher2.Publish(100)
		}
	}

}
