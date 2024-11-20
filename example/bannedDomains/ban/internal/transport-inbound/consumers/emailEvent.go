package consumers

import (
	"context"
	"strings"

	"github.com/gothunder/thunder/example/ban/pkg/events"
	emailEvents "github.com/gothunder/thunder/example/email/pkg/events"
	thunderEvents "github.com/gothunder/thunder/pkg/events"
)

func (c *ConsumerGroup) emailEvent(ctx context.Context, payload emailEvents.EmailPayload) thunderEvents.HandlerResponse {
	domain := strings.Split(payload.Email, "@")[1]
	if c.domains.IsBanned(domain) {
		c.pg.SendBanEvent(ctx, events.BanPayload{
			ID: payload.ID,
		})
	}

	return thunderEvents.Success
}
