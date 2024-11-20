package features

import (
	"github.com/gothunder/thunder/example/ban/internal/features/commands"
	"github.com/gothunder/thunder/example/ban/internal/features/domains"
	"go.uber.org/fx"
)

// Module is a collection of features that can be used in a fx application.
var Module = fx.Options(
	domains.Module,
	commands.Module,
)
