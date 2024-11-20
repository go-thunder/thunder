package consumers

import (
	"context"

	thunderEvents "github.com/gothunder/thunder/pkg/events"
	"github.com/gothunder/thunder/example/ban/pkg/events"
)

func (c *ConsumerGroup) banEvent(ctx context.Context, payload events.BanPayload) thunderEvents.HandlerResponse {
	id := payload.ID
	err := c.emailService.Delete(ctx, id)
	if err != nil {
		return thunderEvents.DeadLetter
	}

	return thunderEvents.Success
}
