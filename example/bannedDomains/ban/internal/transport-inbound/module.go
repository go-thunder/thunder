package transportinbound

import (
	"github.com/gothunder/thunder/example/ban/internal/transport-inbound/consumers"
	"go.uber.org/fx"
)

var Module = fx.Options(
	consumers.Module,
)
