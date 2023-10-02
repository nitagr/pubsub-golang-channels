package types

import (
	"sync"
	"time"
)

type SafeClosingChannel struct {
	C    chan interface{}
	once sync.Once
}
type Message struct {
	Id             string
	TopicName      string
	SubscriberId   string
	Message        interface{}
	CreatedAt      time.Time
	UpdatedAt      time.Time
	RetryFrequency int16
}
