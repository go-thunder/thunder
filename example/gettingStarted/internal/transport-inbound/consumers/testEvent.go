package consumers

import (
	"context"

	"github.com/gothunder/thunder/example/pkg/events"
	thunderEvents "github.com/gothunder/thunder/pkg/events"
)

func (c *ConsumerGroup) testEvent(ctx context.Context, payload events.TestPayload) thunderEvents.HandlerResponse {
	// implement your logic here

	return thunderEvents.Success
}
