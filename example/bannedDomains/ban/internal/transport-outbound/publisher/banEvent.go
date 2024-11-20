package publisher

import (
	"context"

	"github.com/gothunder/thunder/example/ban/pkg/events"
)

func (pg *PublisherGroup) SendBanEvent(ctx context.Context, event events.BanPayload) error {
	return pg.publisher.Publish(ctx, events.BanTopic, event)
}
