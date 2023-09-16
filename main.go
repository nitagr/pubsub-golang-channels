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
var mu3 sync.Mutex

var totalMessagesSent int = 0
var totalMessagesAcknowledged int = 0
var goLimiter100 = make(chan int, 50)
var goLimiter50 = make(chan int, 50)
var goLimiter40 = make(chan int, 50)

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

	go func(message string) {

		goLimiter100 <- 1
		key := topic.PUBSUB_MESSAGE_REDIS
		p.client.Expire(key, time.Second*60*60)

		for _, subscriber := range p.subscribers {

			id := uuid.New()
			data := types.Message{
				Id:      id.String(),
				Message: message,
			}
			subscriber <- data

			mu1.Lock()
			p.client.HSet(key, id.String(), message)
			p.client.HIncrBy(key+":Sent", "totalsent", 1)
			totalMessagesSent++
			mu1.Unlock()

		}
		<-goLimiter100

	}(message)

}

func recieveMessagesOnChannels(sub chan types.Message, subscriberType string, ack chan interface{}, redisClient *redis.Client) {

	goLimiter50 <- 1
	go func() {
		for msg := range sub {

			ack <- msg
			key := topic.PUBSUB_MESSAGE_REDIS
			mu3.Lock()
			// time.Sleep(time.Millisecond * 200)
			redisClient.HDel(key, msg.Id)
			redisClient.HIncrBy(key+":Recieved", "totalrecieved", 1)
			mu3.Unlock()

		}
		<-goLimiter50
	}()

}

func recieveAckOnChannels(ackChannel chan interface{}, wg *sync.WaitGroup) {
	// var wg sync.WaitGroup
	// wg.Add
	goLimiter40 <- 1
	go func() {
		// defer wg.Done()
		for ack := range ackChannel {
			mu2.Lock()
			// defer wg.Done()
			totalMessagesAcknowledged++
			_ = ack
			mu2.Unlock()
		}
		<-goLimiter40
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

			fmt.Println("totalMessagesSent", (redisClient.HGet(topic.PUBSUB_MESSAGE_REDIS+":Sent", "totalsent")))
			fmt.Println("totalMessagesAcknowledged ", (redisClient.HGet(topic.PUBSUB_MESSAGE_REDIS+":Recieved", "totalrecieved")))
			return
		default:

			publisher1.Publish("PUB 1111")
			publisher2.Publish("PUB 2222")
		}
	}

}
