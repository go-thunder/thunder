package consumers

import (
	"context"

	thunderEvents "github.com/gothunder/thunder/pkg/events"
	"github.com/gothunder/thunder/example/email/pkg/events"
	"github.com/rs/zerolog/log"
)

func (c *ConsumerGroup) Handle(ctx context.Context, topic string, decoder thunderEvents.EventDecoder) thunderEvents.HandlerResponse {
	switch {
	case topic == events.EmailTopic:
		var formattedPayload events.EmailPayload
		err := decoder.Decode(&formattedPayload)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to decode payload")
			return thunderEvents.DeadLetter
		}
		log.Ctx(ctx).Debug().Msgf("Got the email '%s' to check for ban", formattedPayload.Email)

		return c.emailEvent(ctx, formattedPayload)
	default:
		return thunderEvents.DeadLetter
	}
}
