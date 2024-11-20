package consumers

import "github.com/gothunder/thunder/example/pkg/events"

func (c *ConsumerGroup) Topics() []string {
	return []string{
		events.TestTopic,
	}
}
