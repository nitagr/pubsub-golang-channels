package main

import (
	"fmt"
	"time"

	"github.com/nitagr/pubsub2/topic"
)

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
			sub <- message
		}(subscriber)

	}
}

func recieveMessagesOnChannels(sub chan interface{}, subscriberType string) {
	go func() {
		for msg := range sub {
			fmt.Print(subscriberType)
			switch v := msg.(type) {
			case int:
				fmt.Printf("Received an int: %d\n", v)
			case string:
				fmt.Printf("Received a string: %s\n", v)
			case float64:
				fmt.Printf("Received a float64: %f\n", v)
			default:
				fmt.Printf("Received an unknown type: %v\n", v)
			}
		}
	}()

}

func main() {

	subs1 := topic.AddSubsriberToTopic(topic.TOPIC_CHANNEL_1)
	subs2 := topic.AddSubsriberToTopic(topic.TOPIC_CHANNEL_2)

	publisher1 := NewPublisher(topic.TOPIC_CHANNEL_1)
	publisher2 := NewPublisher(topic.TOPIC_CHANNEL_2)

	recieveMessagesOnChannels(subs1, "subscriber 1")
	recieveMessagesOnChannels(subs2, "subscriber 2")

	for {
		time.Sleep(1 * time.Second)
		publisher1.Publish("Hello, World! publisher 1")
		publisher2.Publish("Hello, World! publisher 2")
		publisher1.Publish(100)
		publisher2.Publish(200)
	}

}
