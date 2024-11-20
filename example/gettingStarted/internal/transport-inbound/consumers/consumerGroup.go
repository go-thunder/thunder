package consumers

import (
	thunderEvents "github.com/gothunder/thunder/pkg/events"
)

type ConsumerGroup struct{}

func newConsumerGroup() thunderEvents.Handler {
	return &ConsumerGroup{}
}
