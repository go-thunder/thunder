package consumer

import (
	"strings"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rotisserie/eris"
)

func (r *rabbitmqConsumer) declare(routingKeys []string) error {
	r.chManager.ChannelMux.RLock()
	defer r.chManager.ChannelMux.RUnlock()

	dlxName := r.config.QueueName + "_dlx"
	err := r.deadLetterDeclare(dlxName)
	if err != nil {
		return eris.Wrap(err, "failed to declare dead letter")
	}

	err = r.queueDeclare(dlxName)
	if err != nil {
		return eris.Wrap(err, "failed to declare queue")
	}

	err = r.queueBindDeclare(routingKeys)
	if err != nil {
		return eris.Wrap(err, "failed to declare queue bind")
	}

	err = r.chManager.Channel.Qos(
		r.config.PrefetchCount, 0, false,
	)
	if err != nil {
		return eris.Wrap(err, "failed to set QoS")
	}

	return nil
}

func (r *rabbitmqConsumer) queueDeclare(dlxName string) error {
	err := r.chManager.Channel.ExchangeDeclare(
		r.config.ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return eris.Wrap(err, "failed to declare exchange")
	}

	_, err = r.chManager.Channel.QueueDeclare(
		r.config.QueueName,
		true,
		false,
		false,
		false,
		amqp091.Table{
			"x-queue-type":           "quorum",
			"x-dead-letter-exchange": dlxName,
		},
	)
	if err != nil {
		return eris.Wrap(err, "failed to declare queue")
	}

	return nil
}

func (r *rabbitmqConsumer) queueBindDeclare(routingKeys []string) error {
	for _, routingKey := range routingKeys {
		err := r.chManager.Channel.QueueBind(
			r.config.QueueName,
			routingKey,
			r.config.ExchangeName,
			false,
			nil,
		)
		if err != nil {
			return eris.Wrapf(err, "failed to bind queue, topic: %s", routingKey)
		}
	}

	return nil
}

func (r *rabbitmqConsumer) deadLetterDeclare(dlxName string) error {
	dlxProps := amqp091.Table{
		"x-queue-type":             "quorum",
		amqp091.QueueMessageTTLArg: 1000 * 60 * 60 * 24 * 14, // 14 days
		amqp091.QueueMaxLenArg:     10000,                    // 10k messages
	}

	// Try passive declare first
	_, err := r.chManager.Channel.QueueDeclarePassive(
		dlxName,
		true,
		false,
		false,
		false,
		dlxProps,
	)

	if err != nil {
		if !strings.Contains(err.Error(), "NOT_FOUND") &&
			!strings.Contains(err.Error(), "inequivalent arg") {
			return eris.Wrap(err, "unexpected error checking queue")
		}

		// Queue doesn't exist or has different properties - create/recreate
		if strings.Contains(err.Error(), "inequivalent arg") {
			// Delete existing queue if properties don't match
			_, err := r.chManager.Channel.QueueDelete(dlxName, false, false, false)
			if err != nil {
				return eris.Wrap(err, "failed to delete existing queue")
			}
		}

		// Declare exchange
		err = r.chManager.Channel.ExchangeDeclare(
			dlxName,
			"fanout",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return eris.Wrap(err, "failed to declare exchange")
		}

		// Declare queue
		_, err = r.chManager.Channel.QueueDeclare(
			dlxName,
			true,
			false,
			false,
			false,
			dlxProps,
		)
		if err != nil {
			return eris.Wrap(err, "failed to declare queue")
		}

		// Bind queue
		err = r.chManager.Channel.QueueBind(
			dlxName,
			"",
			dlxName,
			false,
			nil,
		)
		if err != nil {
			return eris.Wrap(err, "failed to bind queue")
		}
	}

	return nil
}
