package main

import (
	"fmt"
	"log"
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
var goLimiter100 = make(chan int, 10)
var goLimiter50 = make(chan int, 10)
var goLimiter40 = make(chan int, 10)

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
	goLimiter100 <- 1
	go func(message string) {

		key := topic.PUBSUB_MESSAGE_REDIS
		p.client.Expire(key, time.Second*60*60)
		for _, subscriber := range p.subscribers {

			mu1.Lock()
			id := uuid.New()
			data := types.Message{
				Id:      id.String(),
				Message: message,
			}
			subscriber <- data

			p.client.HSet(key, id.String(), message)
			p.client.HIncrBy(key+":Sent", "totalsent", 1)
			totalMessagesSent++
			mu1.Unlock()

		}
		<-goLimiter100

	}(message)

}

func recieveMessagesOnChannels(sub chan types.Message, subscriberType string, ack chan interface{}, redisClient *redis.Client) {

	go func() {

		for msg := range sub {
			key := topic.PUBSUB_MESSAGE_REDIS
			id := uuid.New()
			redisClient.HSet(key+":RECV", id.String(), msg.Id)
			redisClient.HIncrBy(key+":Recieved", "totalrecieved", 1)
			totalMessagesAcknowledged++
		}
	}()

}

func removeReceivedMessages(wg *sync.WaitGroup, client *redis.Client) {
	key := topic.PUBSUB_MESSAGE_REDIS + ":RECV"
	receivedMessageIds, err := client.HGetAll(key).Result()
	defer wg.Done()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		fmt.Println("len recv", len(receivedMessageIds))
		err := client.Del(key).Err()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Deleted Successfully")

	}()

	for _, value := range receivedMessageIds {
		client.HDel(topic.PUBSUB_MESSAGE_REDIS, value)
	}

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

	timeoutChannel := time.After(time.Second * 2)

	for {
		select {
		case <-timeoutChannel:
			fmt.Println("------")

			fmt.Println("totalMessagesSent", (redisClient.HGet(topic.PUBSUB_MESSAGE_REDIS+":Sent", "totalsent")))
			fmt.Println("totalMessagesAcknowledged ", (redisClient.HGet(topic.PUBSUB_MESSAGE_REDIS+":Recieved", "totalrecieved")))
			fmt.Println("sent ", totalMessagesSent)
			fmt.Println("received ", totalMessagesAcknowledged)
			fmt.Println("---------- CHECKING IF SOME MESSAGES ARE UNSENT------------")
			time.Sleep(time.Second * 2)
			var wg sync.WaitGroup
			wg.Add(1)
			removeReceivedMessages(&wg, redisClient)
			wg.Wait()
			return
		default:
			// time.Sleep(time.Millisecond * 50)
			publisher1.Publish("PUB 1111")
			publisher2.Publish("PUB 2222")

		}
	}

}
