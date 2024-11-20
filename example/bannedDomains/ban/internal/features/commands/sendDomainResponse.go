package commands

import (
	"context"

	"github.com/gothunder/thunder/example/ban/pkg/events"
	"github.com/rs/zerolog/log"
)

// SendDomainResponse sends the result of domain check to RabbitMQ.
func (cg CommandGroup) SendDomainResponse(ctx context.Context, domain string, id int) {
	log.Ctx(ctx).Debug().Msgf("Checking ban for '%s'", domain)
	isBanned := cg.domains.IsBanned(domain)

	if isBanned {
		log.Ctx(ctx).Debug().Msgf("Domain '%s' is banned!", domain)
		cg.publisherGroup.SendBanEvent(
			ctx,
			events.BanPayload{
				ID: id,
			},
		)
	} else {
		log.Ctx(ctx).Debug().Msgf("Domain '%s' is not banned!", domain)
	}
}
