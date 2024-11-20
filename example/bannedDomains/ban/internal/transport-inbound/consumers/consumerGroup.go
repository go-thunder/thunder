package consumers

import (
	"github.com/gothunder/thunder/example/ban/internal/features/domains"
	"github.com/gothunder/thunder/example/ban/internal/transport-outbound/publisher"
	thunderEvents "github.com/gothunder/thunder/pkg/events"
	"go.uber.org/fx"
)

type ConsumerGroup struct {
	domains *domains.Domains
	pg      *publisher.PublisherGroup
}

type ConsumerGroupOptions struct {
	fx.In

	Domains        *domains.Domains
	PublisherGroup *publisher.PublisherGroup
}

func newConsumerGroup(options ConsumerGroupOptions) thunderEvents.Handler {
	return &ConsumerGroup{
		domains: options.Domains,
		pg:      options.PublisherGroup,
	}
}
