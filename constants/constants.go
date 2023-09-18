package constants

import "time"

const (
	DefaultExpiration    = 5 * time.Minute
	PurgeTime            = 10 * time.Minute
	TOPIC_CHANNEL_1      = "TOPIC_CHANNEL_1"
	TOPIC_CHANNEL_2      = "TOPIC_CHANNEL_2"
	PUBSUB_MESSAGE_REDIS = "PUBSUB_MESSAGE_REDIS"
	RECEIVED             = ":RECEIVED"
	SENT                 = ":SENT"
	ACKNOWLEDGE          = ":ACKNOWLEDGE"
	RETRY                = ":RETRY"
)
