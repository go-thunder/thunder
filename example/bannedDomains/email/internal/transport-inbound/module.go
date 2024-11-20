package transportinbound

import (
	"github.com/gothunder/thunder/example/email/internal/transport-inbound/consumers"
	"github.com/gothunder/thunder/example/email/internal/transport-inbound/router"
	"go.uber.org/fx"
)

var Module = fx.Options(
	consumers.Module,
	router.Module,
)
