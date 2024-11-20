# Getting started

The modules are imported to the root of your app, you can then provide them to
the rest of your project, below we're using [Uber's
fx](https://uber-go.github.io/fx/) to perform the dependency injection

Install the library:

```bash
go get github.com/gothunder/thunder
```

Then, create a `main.go` file like this:

```go
// main.go
package main

import (
    transportinbound "github.com/gothunder/thunder/example/internal/transport-inbound"
    transportoutbound "github.com/gothunder/thunder/example/internal/transport-outbound"
    thunderEventRabbitmq "github.com/gothunder/thunder/pkg/events/rabbitmq"
    thunderLogs "github.com/gothunder/thunder/pkg/log"

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

        thunderEventRabbitmq.PublisherModule,
        thunderEventRabbitmq.InvokeConsumer,
    )
    app.Run()

    // This is required to flush the logs to stdout.
    // We only want to do this after the app has exited.
    thunderLogs.DiodeShutdown(w)
}
```
