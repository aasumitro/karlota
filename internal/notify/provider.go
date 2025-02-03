package notify

import (
	"karlota.aasumitro.id/config"
	"karlota.aasumitro.id/internal/utils/mailer"
)

func NewNotificationProvider(infra *config.Infra) {
	service := NewNotificationService(&mailer.SMTPClient{
		Host:      config.GetMailHost(),
		Port:      config.GetMailPort(),
		Username:  config.GetMailUsername(),
		Password:  config.GetMailPassword(),
		TLS:       true,
		LocalName: "https://" + config.GetMailPassword(),
	})
	NewNotificationHandler(service, infra.RabbitMQConsumer)
}
