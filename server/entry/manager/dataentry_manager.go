package manager

import (
	"context"
	jsonEncoder "encoding/json"
	"log"

	"github.com/Carmind-Mindia/user-hub/server/domain"
	data_json "github.com/Carmind-Mindia/user-hub/server/entry/json"
	"github.com/Carmind-Mindia/user-hub/server/services"
	"github.com/rabbitmq/amqp091-go"
)

type DataEntryManager struct {
	repository domain.UserRepository
}

func NewDataEntryManager(repo domain.UserRepository) *DataEntryManager {
	return &DataEntryManager{
		repository: repo,
	}
}

func (d *DataEntryManager) ProcessData(data data_json.ZoneNotification) error {

	fcmtokens, repoErr := d.repository.GetFCMTokensByUserNames(context.Background(), data.Emails)

	if repoErr != nil {
		return repoErr
	}

	var FCMTokens []string
	for _, fcmtoken := range fcmtokens {
		FCMTokens = append(FCMTokens, fcmtoken.FCMToken)
	}

	data.FCMTokens = FCMTokens

	zoneNotificationBytes, _ := jsonEncoder.Marshal(data)
	err := services.GlobalChannel.PublishWithContext(context.Background(), "carmind", "notification.zone.fastemail.ready", false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        zoneNotificationBytes,
	})
	if err != nil {
		log.Println(err)
	}

	return err
}
