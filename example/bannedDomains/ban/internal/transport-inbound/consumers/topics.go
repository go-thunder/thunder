package consumers

import "github.com/gothunder/thunder/example/email/pkg/events"

func (c *ConsumerGroup) Topics() []string {
	return []string{
		events.EmailTopic,
	}
}
