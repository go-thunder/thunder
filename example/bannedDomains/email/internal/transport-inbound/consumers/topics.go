package consumers

import "github.com/gothunder/thunder/example/ban/pkg/events"

func (c *ConsumerGroup) Topics() []string {
	return []string{
		events.BanTopic,
	}
}
