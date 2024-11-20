package publisher

import (
	"context"

	"github.com/gothunder/thunder/example/email/pkg/events"
)

func (pg *PublisherGroup) SendEmailEvent(ctx context.Context, event events.EmailPayload) error {
	return pg.publisher.Publish(ctx, events.EmailTopic, event)
}
