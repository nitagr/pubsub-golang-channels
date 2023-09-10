package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/nitagr/pubsub2/topic"
	"github.com/nitagr/pubsub2/types"
)

var wg1 sync.WaitGroup
var wg2 sync.WaitGroup
var wg3 sync.WaitGroup

var mu1 sync.Mutex
var mu2 sync.Mutex

var totalMessagesSent int = 0
var totalMessagesAcknowledged int = 0

type Publisher struct {
	subscribers []chan types.Message
	client      *redis.Client
}

func NewPublisher(topicName string, client *redis.Client) *Publisher {
	return &Publisher{
		subscribers: topic.TopicMap[topicName],
		client:      client,
	}
}

func (p *Publisher) Publish(message string) {

	for _, subscriber := range p.subscribers {
		go func(sub chan types.Message) {
			mu1.Lock()
			id := uuid.New()
			data := types.Message{
				Id:      id.String(),
				Message: message,
			}
			sub <- data

			key := topic.PUBSUB_MESSAGE_REDIS
			// p.client.HSet(key, id.String(), message)
			p.client.HIncrBy(key+":Sent", "totalsent", 1)
			totalMessagesSent++
			mu1.Unlock()
		}(subscriber)

	}
}

func recieveMessagesOnChannels(sub chan types.Message, subscriberType string, ack chan interface{}, redisClient *redis.Client) {
	go func() {
		for msg := range sub {

			ack <- msg
			key := topic.PUBSUB_MESSAGE_REDIS

			// redisClient.HDel(key, msg.Id)
			redisClient.HIncrBy(key+":Recieved", "totalrecieved", 1)

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
			_ = ack
			mu2.Unlock()
		}
	}()

}

func main() {

	// ns := "Nspubsub" // connecting to redis
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	subs1 := topic.AddSubsriberToTopic(topic.TOPIC_CHANNEL_1)
	subs2 := topic.AddSubsriberToTopic(topic.TOPIC_CHANNEL_2)

	ackChan1 := make(chan interface{})
	ackChan2 := make(chan interface{})

	publisher1 := NewPublisher(topic.TOPIC_CHANNEL_1, redisClient)
	publisher2 := NewPublisher(topic.TOPIC_CHANNEL_2, redisClient)

	recieveMessagesOnChannels(subs1, "subscriber 1", ackChan1, redisClient)
	recieveMessagesOnChannels(subs2, "subscriber 2", ackChan2, redisClient)

	recieveAckOnChannels(ackChan1, &wg1)
	recieveAckOnChannels(ackChan2, &wg2)

	timeoutChannel := time.After(time.Second * 1)
	// wg1.Add(1)

	for {
		select {
		case <-timeoutChannel:
			fmt.Println("------")
			time.Sleep(time.Second * 1)
			fmt.Println("totalMessagesSent %d", redisClient.HGet(topic.PUBSUB_MESSAGE_REDIS+":Sent", "totalsent"))
			fmt.Println("totalMessagesAcknowledged %d", redisClient.HGet(topic.PUBSUB_MESSAGE_REDIS+":Recieved", "totalrecieved"))
			return
		default:

			publisher1.Publish("PUB 1111")
			publisher2.Publish("PUB 2222")
		}
	}

}
