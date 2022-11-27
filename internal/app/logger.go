package app

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const exchangeLoggerName = "logger"

const (
	debugLvl = "debug"
	infoLvl  = "info"
	errorLvl = "error"
)

type logger struct {
}

func newLogger(rmqCh *amqp.Channel) (newLogger *logger, err error) {

	newLogger = &logger{}

	err = newLogger.prepareRMQExchange(rmqCh)
	if err != nil {
		return nil, fmt.Errorf("preparing logger exchange: %w", err)
	}

	return newLogger, nil

}

func (l *logger) prepareRMQExchange(ch *amqp.Channel) (err error) {

	err = ch.ExchangeDeclare(exchangeLoggerName, "direct", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("preparing exchange: %w", err)
	}

	logLvls := []string{debugLvl, infoLvl, errorLvl}
	for _, logLvl := range logLvls {

		queue, err := ch.QueueDeclare("", true, false, true, false, nil)
		if err != nil {
			return fmt.Errorf("declaring queue %s: %w", logLvl, err)
		}

		err = ch.QueueBind(queue.Name, logLvl, exchangeLoggerName, false, nil)
		if err != nil {
			return fmt.Errorf("binding queue %s: %w", logLvl, err)
		}

	}

	return nil
}

func (l *logger) send(ch *amqp.Channel, lvl string, msg []byte) (err error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx, exchangeLoggerName, lvl, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})

	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx, exchangeLoggerName, "", false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})

	if err != nil {
		return err
	}

	return nil

}
