package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeTriggersName    = "triggers"
	queueTriggersRoutingKey = "triggers"
)

type App struct {
	rmqConn           *amqp.Connection
	rmqCh             *amqp.Channel
	triggersQueueName string
	logger            *logger
	client            *client
}

func New() (newApp *App, err error) {

	newApp = &App{}

	rmqConn, rmqCh, err := newApp.prepareRMQChannel()
	if err != nil {
		return nil, fmt.Errorf("preparing rmq channel: %w", err)
	}
	newApp.rmqConn = rmqConn
	newApp.rmqCh = rmqCh

	err = newApp.prepareRMQExchangeTriggers()
	if err != nil {
		rmqConn.Close()
		rmqConn.Close()
		return nil, fmt.Errorf("preparing rmq exchange %s: %w", exchangeTriggersName, err)
	}

	logger, err := newLogger(rmqCh)
	if err != nil {
		rmqConn.Close()
		rmqConn.Close()
		return nil, fmt.Errorf("creating logger: %w", err)
	}
	newApp.logger = logger

	client := newClient()
	newApp.client = client

	return newApp, nil

}

type triggerRecord struct {
	RecordID string `json:"record_id"`
}

func (a *App) Run() (err error) {

	defer func() {
		a.CloseRMQChan()
		a.CloseRMQConn()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	msgs, err := a.rmqCh.Consume(a.triggersQueueName, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("start consuming: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go a.listenDelivery(cancel, msgs)

	log.Printf("Waiting for triggers. To exit press CTRL+C")

	select {
	case <-shutdown:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}

}

func (a *App) listenDelivery(cancel context.CancelFunc, msgs <-chan amqp.Delivery) {

	for msg := range msgs {

		if errSending := a.logger.send(a.rmqCh, debugLvl, msg.Body); errSending != nil {
			log.Printf("producing log debug msg: %s", errSending.Error())
		}

		recordMsg := &triggerRecord{}
		_ = json.Unmarshal(msg.Body, recordMsg)
		if recordMsg.RecordID == "" {
			if errRejecting := msg.Reject(false); errRejecting != nil {
				log.Printf("rejecting msg %s: %s: %s", string(msg.Body), "invalid msg payload", errRejecting.Error())
			} else {
				log.Printf("msg `%s` rejected. Reason: %s", string(msg.Body), "invalid msg payload")
			}
			continue
		}

		phone, code, errGettingPhone := a.client.getPhone(recordMsg.RecordID)
		if errGettingPhone != nil {
			if errors.Is(errGettingPhone, errNoConnectionToPhonesServer) {
				log.Printf("connecting to phones server: %s", errGettingPhone.Error())
				cancel()
				return
			}
			if errRejecting := msg.Reject(true); errRejecting != nil {
				log.Printf("rejecting msg %s: %s: %s", string(msg.Body), errGettingPhone.Error(), errRejecting.Error())
			} else {
				log.Printf("msg `%s` rejected. Reason: %s", string(msg.Body), errGettingPhone.Error())
			}
			continue
		}

		if code == http.StatusInternalServerError {
			if errRejecting := msg.Reject(true); errRejecting != nil {
				log.Printf("rejecting msg %s: %s: %s", string(msg.Body), "failed to process", errRejecting.Error())
			} else {
				log.Printf("msg `%s` rejected. Reason: %s", string(msg.Body), "failed to process")
			}
			continue
		}

		if code == http.StatusNotFound {
			msgNotFound := "Not found: " + recordMsg.RecordID
			if errSending := a.logger.send(a.rmqCh, errorLvl, []byte(msgNotFound)); errSending != nil {
				log.Printf("producing log error msg: %s", errSending.Error())
			}
			if errAck := msg.Ack(false); errAck != nil {
				log.Printf("ack msg `%s`: %s", string(msg.Body), errAck.Error())
			}
			continue
		}

		if code == http.StatusOK {
			msgPhone := msgPhone(recordMsg.RecordID, string(phone))
			if errSending := a.logger.send(a.rmqCh, infoLvl, []byte(msgPhone)); errSending != nil {
				log.Printf("producing log error msg: %s", errSending.Error())
			}
			if errAck := msg.Ack(false); errAck != nil {
				log.Printf("ack msg `%s`: %s", string(msg.Body), errAck.Error())
			}
			continue
		}
	}

}

func msgPhone(id, phone string) (msg string) {
	return "{\"id\": \"" + id + "\", number:\"" + phone + "\"}"
}

func (a *App) prepareRMQChannel() (conn *amqp.Connection, ch *amqp.Channel, err error) {

	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, nil, fmt.Errorf("connecting to RMQ: %w", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("opening channel: %w", err)
	}

	return conn, ch, nil

}

func (a *App) prepareRMQExchangeTriggers() (err error) {

	err = a.rmqCh.ExchangeDeclare(exchangeTriggersName, "direct", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("declaring exhange %s: %w", exchangeTriggersName, err)
	}

	q, err := a.rmqCh.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		return fmt.Errorf("declaring queue: %w", err)
	}

	a.triggersQueueName = q.Name

	err = a.rmqCh.QueueBind(q.Name, queueTriggersRoutingKey, exchangeTriggersName, false, nil)
	if err != nil {
		return fmt.Errorf("binding queue %s with routing key %s: %w", q.Name, queueTriggersRoutingKey, err)
	}

	return nil
}

func (a *App) CloseRMQChan() error {
	return a.rmqConn.Close()
}

func (a *App) CloseRMQConn() error {
	return a.rmqConn.Close()
}

func (a *App) RMQConnIsClosed() bool {
	return a.rmqConn.IsClosed()
}
