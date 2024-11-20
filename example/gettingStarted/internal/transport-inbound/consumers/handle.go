package consumers

import (
	"context"

	"github.com/gothunder/thunder/example/pkg/events"
	thunderEvents "github.com/gothunder/thunder/pkg/events"
	"github.com/rs/zerolog/log"
)

func (c *ConsumerGroup) Handle(ctx context.Context, topic string, decoder thunderEvents.EventDecoder) thunderEvents.HandlerResponse {
	switch {
	case topic == events.TestTopic:
		var formattedPayload events.TestPayload
		err := decoder.Decode(&formattedPayload)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to decode payload")
			return thunderEvents.DeadLetter
		}

		return c.testEvent(ctx, formattedPayload)
	default:
		return thunderEvents.DeadLetter
	}
}
