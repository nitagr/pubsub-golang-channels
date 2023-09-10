package topic

import (
	"github.com/nitagr/pubsub2/types"
)

const (
	TOPIC_CHANNEL_1      = "TOPIC_CHANNEL_1"
	TOPIC_CHANNEL_2      = "TOPIC_CHANNEL_2"
	PUBSUB_MESSAGE_REDIS = "PUBSUB_MESSAGE_REDIS"
)

var TopicMap = make(map[string][]chan types.Message)

func AddSubsriberToTopic(topic string) chan types.Message {
	subs := make(chan types.Message)
	TopicMap[topic] = append(TopicMap[topic], subs)
	return subs
}
