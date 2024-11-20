# Modules

Here is the list of modules in this project:

- [log](./logs.md): Provides a [zerolog](https://github.com/rs/zerolog) logger.
- graphql: Provides a GraphQL handler using [gqlgen](https://gqlgen.com/).
- chi: Provides a [chi](https://go-chi.io/#/) multiplexer with the default
       middlewares. Also exposes a method to start a server with graceful
       shutdown.
- mocks: Provides mocks for consumer and publisher, along with a channel to
         receive intercepted messages.
- [rabbitmq](./events.md): Provides RabbitMQ publisher and consumer with
            [amqp]("https://github.com/rabbitmq/amqp091-go").
- [gRPC](./grpc.md): Provides a gRPC server and client.

## TODO

- [ent](https://entgo.io/): ORM that can be integrated with GraphQL through
  `gqlgen`.
