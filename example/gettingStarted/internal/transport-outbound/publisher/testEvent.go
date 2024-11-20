package publisher

import (
	"context"

	"github.com/gothunder/thunder/example/pkg/events"
)

func (pg *PublisherGroup) SendTestEvent(ctx context.Context, event events.TestPayload) error {
	return pg.publisher.Publish(ctx, events.TestTopic, event)
}
