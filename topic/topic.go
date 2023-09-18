package topic

import (
	"github.com/google/uuid"
	"github.com/nitagr/pubsub2/types"
)

var TopicMap = make(map[string][]string)
var ChannelSubscriberMap = make(map[string]chan types.Message)

func AddSubsriberToTopic(topic string) chan types.Message {
	subs := make(chan types.Message)
	subscriberId := uuid.New()
	TopicMap[topic] = append(TopicMap[topic], subscriberId.String())
	ChannelSubscriberMap[subscriberId.String()] = subs
	return subs
}
