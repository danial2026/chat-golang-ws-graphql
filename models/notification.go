package models

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"chat-golang-react/chat/common"

	"gorm.io/gorm"
)

type Notification struct {
	Id               string `json:"id"`
	To               string `json:"to"`
	From             string `json:"from"`
	NotificationType string `json:"notification_type"`
	Message          string `json:"message"`
	// date
	CreatedAt int `json:"created_at"`
}

// TODO : move to notif service START
type NotificationResponse struct {
	FromId    string `json:"from_id"`
	RecieveId string `json:"recieve_id"`
	From      string `json:"from"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	// date
	CreatedAt int `json:"created_at"`
}

// TODO : move to notif service END

func (r Notification) NewId() string {
	return common.GetUUID()
}

func (r Notification) ToNewStruct(notifId string) (Notification, error) {
	// set inputs
	if notifId == "" {
		notifId = r.NewId()
	}
	notif := Notification{
		Id:               notifId,
		To:               r.Id,
		From:             r.From,
		NotificationType: r.NotificationType,
		Message:          r.Message,
		CreatedAt:        common.Now(),
	}

	return notif, nil
}

func (r Notification) ToNewItem() (Notification, []byte, error) {
	notif, creationErr := r.ToNewStruct(r.Id)
	if creationErr != nil {
		return notif, nil, creationErr
	}
	data, marshalErr := json.Marshal(notif)
	if marshalErr != nil {
		return notif, nil, errors.New("Err in Notification/ToNewItem/Marshal: " + marshalErr.Error())
	}
	return notif, data, nil
}

func (r Notification) Create(ctx context.Context, db *gorm.DB, timeout time.Duration) (Notification, error) {
	notif, _, itemErr := r.ToNewItem()
	if itemErr != nil {
		return notif, itemErr
	}
	err := db.Table(os.Getenv("POSTGRESQLNOTIFSTABLE")).Create(&notif).Error
	if err != nil {
		return notif, errors.New("Err in Notification/Create: " + err.Error())
	}

	return notif, nil
}
