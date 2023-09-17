package types

import "github.com/go-redis/redis"

type Message struct {
	Id      string
	Message interface{}
}

type Publisher struct {
	Subscribers []chan Message
	Client      *redis.Client
}
