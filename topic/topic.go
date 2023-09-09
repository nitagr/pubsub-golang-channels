package topic

const (
	TOPIC_CHANNEL_1 = "TOPIC_CHANNEL_1"
	TOPIC_CHANNEL_2 = "TOPIC_CHANNEL_2"
)

var TopicMap = make(map[string][]chan interface{})

func AddSubsriberToTopic(topic string) chan interface{} {
	subs := make(chan interface{})
	TopicMap[topic] = append(TopicMap[topic], subs)
	return subs
}
