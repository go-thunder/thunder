package publisher

import (
	thunderEvents "github.com/gothunder/thunder/pkg/events"
)

type PublisherGroup struct {
	publisher thunderEvents.EventPublisher
}

func newPublisherGroup(publisher thunderEvents.EventPublisher) *PublisherGroup {
	return &PublisherGroup{
		publisher: publisher,
	}
}
