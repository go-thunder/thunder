package consumer

import (
	"context"

	"github.com/gothunder/thunder/pkg/events"
	"github.com/rotisserie/eris"
)

func (r *rabbitmqConsumer) Subscribe(
	ctx context.Context,
	handler events.Handler,
) error {
	if r.config.DisableConsumer {
		r.logger.Info().Msg("consumer is disabled, skipping subscription")
		// Block until context is cancelled
		<-ctx.Done()
		return nil
	}

	for {
		// Start the go routines that will consume messages
		err := r.startGoRoutines(handler)
		if err != nil {
			return eris.Wrap(err, "failed to start go routines")
		}

		// Check if the channel reconnects
		err = <-r.chManager.NotifyReconnection
		if err != nil {
			return eris.Wrap(err, "failed to reconnect to the amqp channel")
		}

		r.logger.Info().Msg("restarting consumer after reconnection")
	}
}

func (r *rabbitmqConsumer) startGoRoutines(handler events.Handler) error {
	// Declare exchange, queues, and bind them together
	err := r.declare(handler.Topics())
	if err != nil {
		return err
	}

	msgs, err := r.chManager.Channel.Consume(
		r.config.QueueName,
		r.config.ConsumerName,
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		return eris.Wrap(err, "failed to consume messages")
	}

	// We'll keep track of the go routines that we start
	r.wg.Add(r.config.ConsumerConcurrency)
	for i := 0; i < r.config.ConsumerConcurrency; i++ {
		go func() {
			r.handler(msgs, handler)
			// The handler will return when the channel is closed
			r.wg.Done()
		}()
	}
	r.logger.Info().Msgf("processing messages on %v goroutines", r.config.ConsumerConcurrency)

	return nil
}
