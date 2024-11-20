# Events

This module is a library with multiple adapters for processing messages.

We'll assume that you're providing a zerolog logger to the application, a log
entry will be created each time a message is consumed or published.

When consuming messages, a zerolog context will be provided with some default
fields like the topic name populated.

Currently, we support:

- RabbitMQ
- Channel (for testing/mocks)
- more to come in the future...

## Definitions

### `pkg/events`

It's recommended to define the topics and payloads this way you can change them
easily and import them later in other services.

```go
package events

const TestTopic = "topic.test"

type TestPayload struct {
    Hello string `json:"hello"`
}
```

## RabbitMQ

- Automatic reconnection with exponential backoffs
- Graceful shutdowns
- Panic recovery for consumers
- Publishing confirmation
- Re-publishing of messages that couldn't get delivered to the exchange
- Re-publishing of messages that would've been lost during a reconnection or
  instability
- Smart decoder based on the content type of the message
- Configurable parallel consumption of messages using goroutines
- Simple auto-creation of queues, exchanges, binds, and dead letter queues

```go
// main.go
package main

import (
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

        thunderEventRabbitmq.PublisherModule,
        thunderEventRabbitmq.InvokeConsumer,
    )
    app.Run()

    // This is required to flush the logs to stdout.
    // We only want to do this after the app has exited.
    thunderLogs.DiodeShutdown(w)
}
```

### Publisher

If you're using the fx module, the publisher will automatically be started and
closed.

When publishing an event, you'll be sending the struct that will be serialized
and sent to the exchange. If there's any error with the serialization, it'll be
returned back to you.

The message will be published asynchronously, and any errors will be treated
and retried by the module.

```go
type EventPublisher interface {
    // StartPublisher starts the background go routine that will publish messages
    // Returns an error if the publisher fails to start or reconnect
    StartPublisher(context.Context) error

    // Publish publishes a message to the given topic
    // The message is published asynchronously
    // The message will be republished if the connection is lost
    Publish(
        ctx context.Context,
        // The name of the event.
        topic string,
        // The payload of the event.
        payload interface{},
    ) error

    // Close gracefully closes the publisher, making sure all messages are published
    Close(context.Context) error
}
```

### Consumer

Make sure that you define a single handler that matches the interface below.

```go
type HandlerResponse int

const (
    // Default, we remove the message from the queue.
    Success HandlerResponse = iota

    // The message will be delivered to a server configured dead-letter queue.
    DeadLetter

    // Deliver this message to a different worker.
    Retry
)

type EventDecoder interface {
    // Decode decodes the payload into the given interface.
    // Returns an error if the payload cannot be decoded.
    Decode(v interface{}) error
}

type Handler interface {
    // The function that will be called when a message is received.
    Handle(ctx context.Context, topic string, decoder EventDecoder) HandlerResponse
    // The topics that will be subscribed to.
    Topics() []string
}
```

You can find an example of a consumer [here](https://github.com/gothunder/thunder/tree/main/example/internal/transport-inbound/consumers).
