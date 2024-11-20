module github.com/gothunder/thunder/example/ban

go 1.19

replace github.com/gothunder/thunder/example/email => ../email

require (
	github.com/gothunder/thunder v0.5.1
	github.com/gothunder/thunder/example/email v0.0.0-20230102180253-e0b111ffa5c9
	github.com/rs/zerolog v1.28.0
	go.uber.org/fx v1.18.2
)

require (
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/rabbitmq/amqp091-go v1.5.0 // indirect
	github.com/rotisserie/eris v0.5.4 // indirect
	github.com/vmihailenco/msgpack/v5 v5.3.5 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/dig v1.15.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.17.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
)
