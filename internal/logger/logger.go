package logger

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

type LogConsumer struct {
	exchangeName string
	routingKeys  []string
}

func New() *LogConsumer {

	newLogConsumer := &LogConsumer{}

	newLogConsumer.exchangeName = "logger"

	routingKeys := []string{"debug", "info", "error", ""}
	newLogConsumer.routingKeys = routingKeys

	return newLogConsumer

}

func (l *LogConsumer) Run() error {

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("connecting to RMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("opening channel: %w", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(l.exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declaring exhange: %w", err)
	}

	for _, rKey := range l.routingKeys {

		q, err := ch.QueueDeclare("", true, false, true, false, nil)
		if err != nil {
			return fmt.Errorf("declaring queue with routing key %s: %w", rKey, err)
		}

		err = ch.QueueBind(q.Name, rKey, l.exchangeName, false, nil)
		if err != nil {
			return fmt.Errorf("binding queue %s with routing key %s: %w", q.Name, rKey, err)
		}

		msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("queue name: %s, routing key: %s, starting delivery: %w", q.Name, rKey, err)
		}

		go listenDelivery(rKey, msgs)
	}

	log.Printf("Waiting for logs. To exit press CTRL+C")
	<-shutdown

	return nil
}

func listenDelivery(lvl string, msgs <-chan amqp.Delivery) {
	if lvl == "" {
		lvl = "common"
	}
	for d := range msgs {
		log.Printf("{lvl: %s}, {msg: %s}", lvl, d.Body)
		if err := d.Ack(false); err != nil {
			log.Printf("ack msg `%s`: %s", string(d.Body), err.Error())
		}
	}
}
