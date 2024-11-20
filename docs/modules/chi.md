# Chi

This module uses [chi](https://go-chi.io/#/) as the HTTP multiplexer, providing
a easy way to add or remove routes.

## Starting the server

To start the server, you need to invoke the `thunderChi.StartServer` function.
This is ideally done in the `main` package using [Uber's
fx](https://uber-go.github.io/fx/). The snippet below shows how to do this
using the `fx.Invoke` method.

```go
// main.go
package main

import (
    thunderLogs "github.com/gothunder/thunder/pkg/log"
	thunderChi "github.com/gothunder/thunder/pkg/router/chi"
	transportInbound "github.com/gothunder/thunder/example/email/internal/transport-inbound"

    "github.com/rs/zerolog/diode"
    "go.uber.org/fx"
)

func main() {
    var w diode.Writer

    app := fx.New(
        // The order of these options isn't important.
        thunderLogs.Module,
        fx.Populate(&w),

        thunderChi.Module, // This is the module we're interested in.
        fx.Invoke(thunderChi.StartServer), // This starts the server.
        transportInbound.Module, // This is the module where the routes are defined.
    )
    app.Run()

    // This is required to flush the logs to stdout.
    // We only want to do this after the app has exited.
    thunderLogs.DiodeShutdown(w)
}
```

## Defining routes

To define routes, you need to create a module that implements the
`HTTPHandler` interface, that is:

- `Method() string` - The HTTP method of the route, e.g. `GET`, `POST`, etc.
- `Pattern() string` - The pattern of the route, e.g. `/users`.
- `ServeHTTP(w http.ResponseWriter, r *http.Request)` - The handler function
  that will be called when the route is hit.

The snippet below shows how to define a route that returns a list of users.

```go
// transport-inbound/router/handlers.go
package router

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gothunder/thunder/pkg/router"
)

func NewUsersHandler() router.HandlerOutput {
	return router.HandlerOutput{
		Handler: UsersHandler{},
	}
}

type UsersHandler struct {}

func (h UsersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    users := []string{"user1", "user2", "user3"}

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(fmt.Sprintf("%s", users)))
}

func (h UsersHandler) Method() string {
	return http.MethodGet
}

func (h UsersHandler) Pattern() string {
	return "/users"
}
```

## Adding routes

```go
// transport-inbound/router/module.go
package router

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewUsersHandler,
	),
)
```

## TODO

- Elaborate further on how to use this module.
