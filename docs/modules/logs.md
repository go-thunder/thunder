# Logs

This module wraps up the [zerolog](https://github.com/rs/zerolog) package.

The other modules will be expecting you to provide a `zerolog` instance to them.

Through the wrapper, you'll be able to configure the logging level via
environment variables, as well as enable pretty console output while in
development.

Additionally, there's a dedicated marshaller for error stack traces.

The logger is non-blocking, so you should make sure that all logs are flushed
before the app is closed (an example of the complete setup is provided below).

```go
package main

import (
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
    )
    app.Run()

    // This is required to flush the logs to stdout.
    // We only want to do this after the app has exited.
    thunderLogs.DiodeShutdown(w)
}
```
