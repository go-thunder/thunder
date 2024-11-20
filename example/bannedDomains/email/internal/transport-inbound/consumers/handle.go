package consumers

import (
	"context"

	thunderEvents "github.com/gothunder/thunder/pkg/events"
	"github.com/gothunder/thunder/example/ban/pkg/events"
	emailEvents "github.com/gothunder/thunder/example/email/pkg/events"
	"github.com/rs/zerolog/log"
)

func (c *ConsumerGroup) Handle(ctx context.Context, topic string, decoder thunderEvents.EventDecoder) thunderEvents.HandlerResponse {
	switch {
	case topic == events.BanTopic:
		var formattedPayload events.BanPayload
		err := decoder.Decode(&formattedPayload)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to decode payload")
			return thunderEvents.DeadLetter
		}
		log.Ctx(ctx).Debug().Msgf("Got ban request for domain ID '%d'", formattedPayload.ID)

		return c.banEvent(ctx, formattedPayload)
	case topic == emailEvents.EmailTopic:
		log.Ctx(ctx).Debug().Msg("got email event")
		return thunderEvents.Retry
	default:
		return thunderEvents.DeadLetter
	}
}
