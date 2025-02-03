package notify

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"karlota.aasumitro.id/internal/common"
)

type notificationHandler struct {
	service  INotificationService
	rabbitMQ *amqp.Connection
}

func (handler notificationHandler) EmailEvent() {
	ch, err := handler.rabbitMQ.Channel()
	if err != nil {
		log.Println(err.Error())
		return
	}
	q, err := ch.QueueDeclare(
		common.EmailEventTopic, false,
		false, false, false, nil)
	if err != nil {
		log.Println("failed to declare queue: ", err)
		return
	}
	messages, err := ch.Consume(
		q.Name, "", false,
		false, false, false, nil)
	if err != nil {
		log.Println("failed to consume queue: ", err)
		return
	}
	go func() {
		for d := range messages {
			if err := handler.service.SendEmail(d.Body); err != nil {
				log.Println("failed to send email: ", err.Error())
				continue
			}
			if err := d.Ack(false); err != nil {
				log.Println("failed to acknowledge email queue: ", err.Error())
				continue
			}
		}
	}()
}

func NewNotificationHandler(
	service INotificationService,
	rabbitMQ *amqp.Connection,
) {
	handler := &notificationHandler{service: service, rabbitMQ: rabbitMQ}
	handler.EmailEvent()
	// todo register queue worker
	// so let say user did not online then we will store the notify inside the queue table
	// and then will proceed the item when user back online
	// this will applied to, new conversation message, new group invite, etc.
}
