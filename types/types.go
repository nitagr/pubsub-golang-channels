package types

import "time"

type Message struct {
	Id           string
	TopicName    string
	SubscriberId string
	Message      interface{}
	CreatedAt    time.Time
}
