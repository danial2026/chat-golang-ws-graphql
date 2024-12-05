package models

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"gorm.io/gorm"

	"chat-golang-react/chat/common"
)

// usages:
// delete message: by id if user has permission
// edit message: by id if user has permission
// forward message: by id
// reply message: by id
// list of the messages: by room name and sorted by creation date
// reaction on the message: by id if username is other than how created it

// message types:
const (
	MessageTextType      = "text"
	MessageLinkType      = "link"
	MessageLocationType  = "location"
	MessageImageType     = "image"
	MessageVoiceType     = "voice"
	MessageVideoType     = "video"
	MessagePDFType       = "pdf"
	MessageVideoCallType = "video_call"
	MessageVoiceCallType = "voice_call"
	MessageOtherMedia    = "others"
)

type MessageModel struct {
	// general fields
	Id       string `json:"id"`
	GUID     string `json:"guid"`
	Username string `json:"username"`
	UserId   string `json:"user_id"`
	RoomId   string `json:"room_id"`
	// types can be text, link, location, image, voice, video, pdf , video_call and other files
	Type string `json:"type"`
	// text
	TextContent string `json:"text_content"`
	// link url
	LinkUrl string `json:"link_url,omitempty"`
	// date
	CreatedAt  int    `json:"created_at"`
	UpdatedAt  int    `json:"updated_at,omitempty"`
	DeletedAt  int    `json:"deleted_at,omitempty"`
	DeletedFor string `json:"deleted_for,omitempty"`
}

func (m MessageModel) NewId() string {
	return common.GetULID()
}

func (m MessageModel) NewGUID() string {
	return common.GetUUID()
}

func (m MessageModel) ItemToStruct(item []byte) (MessageModel, error) {
	if item == nil {
		return MessageModel{}, nil
	}
	var msg MessageModel
	if err := json.Unmarshal(item, &msg); err != nil {
		return MessageModel{}, err
	}
	return msg, nil
}

func (m MessageModel) GetItem(ctx context.Context, db *gorm.DB, timeout time.Duration) (MessageModel, error) {
	log.Println("MessageModel/GetItem/Query: ", m.Id)
	var result MessageModel
	err := db.Table(os.Getenv("POSTGRESQLMESSAGESTABLE")).Where(`"id" = ?`, m.Id).Find(&result).Error
	if err != nil {
		return MessageModel{}, errors.New("Err in MessageModel/GetItem/Query: scan error")
	}
	return result, nil
}

func (m MessageModel) GetAllItems(ctx context.Context, db *gorm.DB, timeout time.Duration, limit, page int) ([]MessageModel, error) {
	log.Println("MessageModel/GetAllItems/Query: ")
	var result []MessageModel
	err := db.Table(os.Getenv("POSTGRESQLMESSAGESTABLE")).Where(`"deleted_at" = 0 AND "room_id" = ?`, m.RoomId).Order("created_at DESC").Limit(limit).Offset(limit * page).Find(&result).Error
	if err != nil {
		return []MessageModel{}, errors.New("Err in MessageModel/GetAllItems/Query: scan error : " + err.Error())
	}

	return result, nil
}

func (m MessageModel) Exists(ctx context.Context, db *gorm.DB, timeout time.Duration) (MessageModel, bool, error) {
	msg, getErr := m.GetItem(ctx, db, timeout)
	if getErr != nil {
		return msg, false, getErr
	}
	if msg.Id != m.Id {
		return msg, false, nil
	}
	if msg.DeletedAt > 0 {
		return msg, false, nil
	}
	// TODO : check if its not deleted for this user
	return msg, true, nil
}

func (m MessageModel) ToNewStruct() (MessageModel, error) {
	// set inputs

	msg := MessageModel{
		Id:          m.NewId(),
		GUID:        m.NewGUID(),
		Username:    m.Username,
		UserId:      m.UserId,
		RoomId:      m.RoomId,
		Type:        m.Type,
		TextContent: m.TextContent,
		LinkUrl:     m.LinkUrl,
		CreatedAt:   common.Now(),
	}

	return msg, nil
}

func (m MessageModel) ToNewItem() (MessageModel, []byte, error) {
	msg, creationErr := m.ToNewStruct()
	if creationErr != nil {
		return msg, nil, creationErr
	}
	data, marshalErr := json.Marshal(msg)
	if marshalErr != nil {
		return msg, nil, errors.New("Err in MessageModel/ToNewItem/Marshal: " + marshalErr.Error())
	}
	return msg, data, nil
}

func (m MessageModel) Create(ctx context.Context, db *gorm.DB, timeout time.Duration) (MessageModel, error) {
	room, _, itemErr := m.ToNewItem()
	if itemErr != nil {
		return room, itemErr
	}
	err := db.Table(os.Getenv("POSTGRESQLMESSAGESTABLE")).Create(&room).Error
	if err != nil {
		return room, errors.New("Err in MessageModel/Create: " + err.Error())
	}

	return room, nil
}

func (m MessageModel) DeleteItemFor(ctx context.Context, db *gorm.DB, timeout time.Duration, deleteFor string) error {

	log.Println("MessageModel/DeleteItemFor/Query: ", m.Id)
	err := db.Table(os.Getenv("POSTGRESQLMESSAGESTABLE")).Where(`"id" = `, m.Id).Updates(map[string]interface{}{
		"deleted_at":  int(time.Now().Unix()),
		"deleted_for": deleteFor,
	}).Error
	if err != nil {
		return errors.New("Err in MessageModel/DeleteItemFor: " + err.Error())
	}

	return nil
}
