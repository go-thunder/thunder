package transportinbound

import (
	"github.com/gothunder/thunder/example/internal/transport-inbound/consumers"
	"go.uber.org/fx"
)

var Module = fx.Options(
	consumers.Module,
)
