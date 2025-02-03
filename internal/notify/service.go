package notify

import (
	"html/template"
	"log"
	"net/mail"
	"sync"
	"time"

	"karlota.aasumitro.id/internal/model/entity"
	"karlota.aasumitro.id/internal/utils/mailer"
)

type (
	INotificationService interface {
		SendEmail(message []byte) error
	}

	notificationService struct {
		mc *mailer.SMTPClient
	}
)

const sleepDuration = 1 * time.Second

func (service notificationService) SendEmail(message []byte) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	// extract data from message
	notify := &entity.Notify{}
	if err := notify.Decode(message); err != nil {
		return err
	}
	// build body
	params := struct{ HTMLContent template.HTML }{
		HTMLContent: template.HTML(notify.Body)}
	body, err := mailer.BuildMailBody(params)
	if err != nil {
		log.Println("Error building mail:", err)
		return err
	}
	time.Sleep(sleepDuration)
	return service.mc.Send(&mailer.Message{
		From:    mail.Address{Name: service.mc.Username, Address: service.mc.Username},
		To:      []mail.Address{{Address: notify.Target}},
		Subject: notify.Subject,
		HTML:    body,
	})
}

func NewNotificationService(mc *mailer.SMTPClient) INotificationService {
	return &notificationService{mc: mc}
}
