# Banned domains

This example show how to use thunder to create two services that communicate
with each other using RabbitMQ. The first service is an API that receives an
email and stores it in a database. The second service is a worker that reads
the emails from the database and checks if the domain is banned. If it is, it
publishes a message to a queue that the first service reads and deletes the
email from the database.

Note that in the real world there is a gap between the time the email is stored
in the database and the time it's checked against the ban list.

## Setup

- Spin up a RabbitMQ container
- cd into the `email` directory and run `go run main.go`
- cd into the `ban` directory and run `go run main.go`
