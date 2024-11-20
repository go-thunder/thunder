# Structuring your project

## Directories

We recommend that you separate each part of the application, isolating the
business logic and following a clean architecture, the idea is that you should
be able to refactor each part individually without affecting other parts of the
app.

In this section, we'll give you an opinionated structure, it's based on
standards like CQRS, clean architecture, and hexagonal microservices. Feel free
to adapt this structure to your own needs, however, it should serve you well
especially as you scale your app codebase and complexity.

This will be the structure used for any examples and docs moving forward.

You can check [example](https://github.com/gothunder/thunder/tree/main/example)
to get a better idea of how this structure looks like.

### `internal`

This is a special directory for Go, everything that is inside can't be imported
outside of your app, think about this as similar to an src folder in languages
like JavaScript/Node.JS.

### `internal/features`

All domain logic

### `internal/features/commands`

Exposes methods that are used to change data, send events, making requests,
etc. You may also perform some validation and query existing data before making
any actions. Think of this as a controller or service layer.

It's strongly recommended to use database transactions for complex commands
since you may want to roll back your data if there are any failures.

You should also, always have the assumption that this command may be called
multiple times by your transports, maybe a user refreshed the page resending
the same request, or you received duplicated webhooks or even events were
consumed twice or more due to an instability or bug within your app.

### `internal/features/queries`

Exposes methods used to query the database, you may include default filters and
additional logic here before returning data. You shouldn't have anything in
here that changes the data available in the database.

### `internal/features/repository`

This is a simple abstraction for interacting with your database or ORM. You
should avoid adding business logic here. You should also avoid interacting with
the database / ORM without the use of your repository.

### `internal/features/<feature-name>`

This is where the code for any other feature that does not fit into the above
categories or transport modules should go.

### `pkg/grpc`

This is where `protoc` outputs the generated code for your gRPC services. You
shouldn't add any code here manually.

### Examples

#### `internal/features/controller`

Abstraction to use the same logic for resolvers, REST and consumers. Note that
it may not be as useful since the logic may vary a lot between transports.

#### `internal/features/domains`

Internal feature to work with domains, including rules, parsers, etc. If used
by other modules, it should be moved to `pkg`.

### `internal/transport-inbound`

Every communication that comes into your app, includes API routes, GraphQL
resolvers, event consumers, webhooks, etc.

This directory may be further split into subdirectories based on the type of
transport, that being:

- `graphql` for GraphQL resolvers and schemas
- `routes` for HTTP server handlers
- `webhooks` for webhooks
- `consumers` for event consumers

### `internal/transport-outbound`

Every communication that your app does, including event publishers, code that
interacts with external APIs, etc.

### `pkg`

This is similar to a dist folder in languages like JavaScript/Node.JS, here
you can add any code that is meant to be imported by other services. For
example event definitions.
