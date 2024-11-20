package main

import (
	thunderEventRabbitmq "github.com/gothunder/thunder/pkg/events/rabbitmq"
	thunderLogs "github.com/gothunder/thunder/pkg/log"
	"github.com/gothunder/thunder/example/ban/internal/features"
	transportinbound "github.com/gothunder/thunder/example/ban/internal/transport-inbound"
	transportoutbound "github.com/gothunder/thunder/example/ban/internal/transport-outbound"

	"github.com/rs/zerolog/diode"
	"go.uber.org/fx"
)

func main() {
	var w diode.Writer

	app := fx.New(
		// The order of these options isn't important.
		thunderLogs.Module,
		fx.Populate(&w),

		transportinbound.Module,
		transportoutbound.Module,
		features.Module,

		thunderEventRabbitmq.PublisherModule,
		thunderEventRabbitmq.InvokeConsumer,
	)
	app.Run()

	// This is required to flush the logs to stdout.
	// We only want to do this after the app has exited.
	thunderLogs.DiodeShutdown(w)
}
