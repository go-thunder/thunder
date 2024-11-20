package commands

import (
	"github.com/gothunder/thunder/example/ban/internal/features/domains"
	"github.com/gothunder/thunder/example/ban/internal/transport-outbound/publisher"
	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

type CommandGroup struct {
	publisherGroup *publisher.PublisherGroup
	log            *zerolog.Logger
	domains        *domains.Domains
}

type CommandGroupInput struct {
	fx.In

	PublisherGroup *publisher.PublisherGroup
	Log            *zerolog.Logger
	Domains        *domains.Domains
}

func NewCommandGroup(input CommandGroupInput) *CommandGroup {
	return &CommandGroup{
		publisherGroup: input.PublisherGroup,
		log:            input.Log,
		domains:        input.Domains,
	}
}
