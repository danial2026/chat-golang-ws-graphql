package models

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"gorm.io/gorm"

	"chat-golang-react/chat/common"

	"log"
	"time"
)

// usage:
// create a room: by title, type, biography, profile picture:
// remove a room: by creator username, and roomId
// list of rooms by last activity

const ConversationsListRoom = "conversion-list-room"

// room types:
const (
	RoomMutualType  = "mutual"
	RoomGroupType   = "group"
	RoomChannelType = "channel"
)

type RoomModel struct {
	// general
	Id        string `json:"id"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	Biography string `json:"biography"`
	Creator   string `json:"creator"`
	CreatorId string `json:"creator_id"`
	// types is mutual, chat_group and channel
	Type string `json:"type"`
	// status
	IsReported bool `json:"is_reported"`
	IsActive   bool `json:"is_active"`
	IsPublic   bool `json:"is_public"`
	// date
	CreatedAt  int    `json:"created_at"`
	UpdatedAt  int    `json:"updated_at"`
	DeletedAt  int    `json:"deleted_at"`
	DeletedFor string `json:"deleted_for"`
}

func (r RoomModel) NewId() string {
	if (r.Id != "") {
		return r.Id
	}
	if (r.Type == RoomGroupType) {
		return "group-" + common.GetUUID()
	} else if (r.Type == RoomMutualType) {
		// this should never get called
		return common.GetUUID()
	} else if (r.Type == RoomChannelType) {
		return "channel-" + common.GetUUID()
	} else {
		return common.GetUUID()
	}
}

func (r RoomModel) ItemToStruct(item []byte) (RoomModel, error) {
	if item == nil {
		return RoomModel{}, nil
	}
	var room RoomModel
	if err := json.Unmarshal(item, &room); err != nil {
		return RoomModel{}, errors.New("Err in RoomModel/ItemToStruct/Unmarshal: " + err.Error())
	}
	return room, nil
}

func (r RoomModel) GetItem(ctx context.Context, db *gorm.DB, timeout time.Duration) (RoomModel, error) {
	log.Println("RoomModel/GetItem/Query: ", r.Id)
	var result RoomModel
	err := db.Table(os.Getenv("POSTGRESQLROOMSTABLE")).Where(`"id" = ?`, r.Id).Find(&result).Error
	if err != nil {
		return RoomModel{}, errors.New("Err in RoomModel/GetItem/Query: scan error : " + err.Error())
	}

	return result, nil
}

func (r RoomModel) GetAllItems(ctx context.Context, db *gorm.DB, timeout time.Duration, limit, page int) ([]RoomModel, error) {
	log.Println("RoomModel/GetAllItems/Query: ")
	var result []RoomModel
	err := db.Table(os.Getenv("POSTGRESQLROOMSTABLE")).Limit(limit).Offset(limit * page).Find(&result).Error
	if err != nil {
		return []RoomModel{}, errors.New("Err in RoomModel/GetAllItems/Query: scan error : " + err.Error())
	}

	return result, nil
}

func (r RoomModel) Exists(ctx context.Context, db *gorm.DB, timeout time.Duration) (RoomModel, bool, error) {
	room, getErr := r.GetItem(ctx, db, timeout)
	if getErr != nil {
		return RoomModel{}, false, getErr
	}
	if room.Id != r.Id {
		return RoomModel{}, false, nil
	}
	if !room.IsActive {
		return room, false, nil
	}
	if room.DeletedAt > 0 {
		return room, false, nil
	}
	return room, true, nil
}

func (r RoomModel) ToNewStruct(roomId string) (RoomModel, error) {
	room := RoomModel{
		Id:         r.NewId(),
		Title:      r.Title,
		Image:      r.Image,
		Biography:  r.Biography,
		Creator:    r.Creator,
		CreatorId:  r.CreatorId,
		DeletedFor: r.DeletedFor,
		Type:       r.Type,
		IsReported: false,
		IsActive:   true,
		IsPublic:   r.IsPublic,
		CreatedAt:  common.Now(),
		UpdatedAt:  0,
		DeletedAt:  0,
	}

	return room, nil
}

func (r RoomModel) ToNewItem() (RoomModel, []byte, error) {
	room, creationErr := r.ToNewStruct(r.Id)
	if creationErr != nil {
		return room, nil, creationErr
	}
	data, marshalErr := json.Marshal(room)
	if marshalErr != nil {
		return room, nil, errors.New("Err in RoomModel/ToNewItem/Marshal: " + marshalErr.Error())
	}
	return room, data, nil
}

func (r RoomModel) Create(ctx context.Context, db *gorm.DB, timeout time.Duration) (RoomModel, error) {
	// TODO: find out why i had to"in mutual type roomId is required"?
	// if r.Type == RoomMutualType && r.Id == "" {
	// 	return RoomModel{}, errors.New("err in RoomModel/Create/ToNewItem: for mutual type room roomId required")
	// }
	room, _, itemErr := r.ToNewItem()
	if itemErr != nil {
		return room, itemErr
	}
	err := db.Table(os.Getenv("POSTGRESQLROOMSTABLE")).Create(&room).Error
	if err != nil {
		return room, errors.New("Err in RoomModel/Create: " + err.Error())
	}

	return room, nil
}

func (r RoomModel) DeleteItemFor(ctx context.Context, db *gorm.DB, timeout time.Duration, deleteFor string) error {

	log.Println("RoomModel/DeleteItemFor/Query: ", r.Id)
	err := db.Table(os.Getenv("POSTGRESQLROOMSTABLE")).Where(`"id" = `, r.Id).Updates(map[string]interface{}{
		"deleted_at":  int(time.Now().Unix()),
		"deleted_for": deleteFor,
	}).Error
	if err != nil {
		return errors.New("Err in RoomModel/DeleteItemFor: " + err.Error())
	}

	return nil
}
