package redisclient

import (
	"github.com/go-redis/redis"
)

type RedisClientPubSub struct {
	client *redis.Client
}

const (
	defaultNs = "NSpubsub"
)

// func NewRedisClient(client *redis.Client) *RedisClientPubSub {
// 	if client == nil {
// 		panic("")
// 	}

// 	if ns == "" {
// 		ns = defaultNs
// 	}
// 	if !strings.HasSuffix(ns, ":") {
// 		ns += ":" // GoMQ:
// 	}

// 	rsmq := &RedisClientPubSub{
// 		client: client,
// 	}

// 	return rsmq
// }
