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

// usages:
// joinMember: with username and room id
// removeMember: with username and room id
// listOfMembers: with roomId and sorted by last date online
// is admin: by the room id and username
// mute room: by the roomId and username
// make a user admin

type RoomMembershipModel struct {
	Id       string `json:"id"`
	RoomId   string `json:"room_id"`
	JoinBy   string `json:"join_by"`
	JoinById string `json:"join_by_id"`
	Username string `json:"username"`
	UserId   string `json:"user_id"`
	IsAdmin  bool   `json:"is_admin"`
	// mute related's
	MuteUntil int `json:"mute_until"`
	// date
	JoinAt  int `json:"join_at"`
	LeaveAt int `json:"leave_at"`
}

func (m RoomMembershipModel) NewId() string {
	return common.GetUUID()
}

func (m RoomMembershipModel) ItemToStruct(item []byte) (RoomMembershipModel, error) {
	if item == nil {
		return RoomMembershipModel{}, nil
	}
	var member RoomMembershipModel
	if err := json.Unmarshal(item, &member); err != nil {
		return RoomMembershipModel{}, errors.New("Err in RoomMembershipModel/ItemToStruct/UnmarshalMap: " + err.Error())
	}
	return member, nil
}

func (m RoomMembershipModel) GetItem(ctx context.Context, db *gorm.DB, timeout time.Duration) (RoomMembershipModel, error) {
	log.Println("RoomMembershipModel/GetItem/Query: ", m.Id)
	var result RoomMembershipModel
	err := db.Table(os.Getenv("POSTGRESQLMEMBERSHIPTABLE")).Where(`"id" = ?`, m.Id).Find(&result).Error
	if err != nil {
		return RoomMembershipModel{}, errors.New("Err in RoomMembershipModel/GetItem/Query: scan error")
	}
	return result, nil
}

func (m RoomMembershipModel) GetItemByUserIdAndRoomId(ctx context.Context, db *gorm.DB, timeout time.Duration) (RoomMembershipModel, error) {
	log.Println("RoomMembershipModel/GetItemByUserIdAndRoomId/Query: ", m.UserId, m.RoomId, m.Username)
	var result RoomMembershipModel
	err := db.Table(os.Getenv("POSTGRESQLMEMBERSHIPTABLE")).Where(`"room_id" = ? AND "user_id" = ?`, m.RoomId, m.UserId).Find(&result).Error
	if err != nil {
		return RoomMembershipModel{}, errors.New("Err in RoomMembershipModel/GetItem/Query: scan error")
	}
	return result, nil
}

func (m RoomMembershipModel) Exists(ctx context.Context, db *gorm.DB, timeout time.Duration) (bool, error) {
	member, getErr := m.GetItemByUserIdAndRoomId(ctx, db, timeout)
	if getErr != nil {
		return false, getErr
	}
	if member.RoomId != m.RoomId {
		return false, nil
	}
	if member.LeaveAt > 0 {
		return false, nil
	}
	return true, nil
}

func (m RoomMembershipModel) ToNewStruct() (RoomMembershipModel, error) {
	member := RoomMembershipModel{
		Id:        m.NewId(),
		RoomId:    m.RoomId,
		JoinBy:    m.JoinBy,
		JoinById:  m.JoinById,
		Username:  m.Username,
		UserId:    m.UserId,
		IsAdmin:   m.IsAdmin,
		MuteUntil: 0,
		JoinAt:    common.Now(),
		LeaveAt:   0,
	}

	return member, nil
}

func (m RoomMembershipModel) ToNewItem() (RoomMembershipModel, []byte, error) {
	member, creationErr := m.ToNewStruct()
	if creationErr != nil {
		return member, nil, creationErr
	}
	data, marshalErr := json.Marshal(member)
	if marshalErr != nil {
		return member, nil, errors.New("Err in RoomMembershipModel/ToNewItem/MarshalMap: " + marshalErr.Error())
	}
	return member, data, nil
}

func (m RoomMembershipModel) Create(ctx context.Context, db *gorm.DB, timeout time.Duration) (RoomMembershipModel, error) {
	member, _, itemErr := m.ToNewItem()
	if itemErr != nil {
		return member, itemErr
	}
	err := db.Table(os.Getenv("POSTGRESQLMEMBERSHIPTABLE")).Create(&member).Error
	if err != nil {
		return RoomMembershipModel{}, errors.New("Err in RoomMembershipModel/Create/Create: Insert error")
	}
	return RoomMembershipModel{}, nil
}

func (m RoomMembershipModel) GetByRoomId(ctx context.Context, db *gorm.DB, timeout time.Duration) ([]RoomMembershipModel, error) {
	var results []RoomMembershipModel
	err := db.Table(os.Getenv("POSTGRESQLMEMBERSHIPTABLE")).Where(`"room_id" = '` + m.RoomId + "'").Find(&results).Error
	if err != nil {
		log.Println("RoomMembershipModel/GetByRoomId/Query: ", err)
		return []RoomMembershipModel{}, errors.New("Err in RoomMembershipModel/GetByRoomId/Query: " + err.Error())
	}
	return results, nil
}

func (m RoomMembershipModel) GetByRoomIdWithPagination(ctx context.Context, db *gorm.DB, timeout time.Duration, limit, page int) ([]RoomMembershipModel, error) {
	var results []RoomMembershipModel
	err := db.Table(os.Getenv("POSTGRESQLMEMBERSHIPTABLE")).Where(`"room_id" = '` + m.RoomId + "'").Limit(limit).Offset(limit * page).Find(&results).Error
	if err != nil {
		log.Println("RoomMembershipModel/GetByRoomIdWithPagination/Query: ", err)
		return []RoomMembershipModel{}, errors.New("Err in RoomMembershipModel/GetByRoomIdWithPagination/Query: " + err.Error())
	}
	return results, nil
}

func (m RoomMembershipModel) GetByUserId(ctx context.Context, db *gorm.DB, timeout time.Duration, limit, page int) ([]RoomMembershipModel, error) {
	var results []RoomMembershipModel
	err := db.Table(os.Getenv("POSTGRESQLMEMBERSHIPTABLE")).Where(`"user_id" = '` + m.UserId + "'" + ` and "leave_at" = '0' `).Limit(limit).Offset(limit * page).Find(&results).Error
	if err != nil {
		log.Println("RoomMembershipModel/GetByUserId/Query: ", err)
		return []RoomMembershipModel{}, errors.New("Err in RoomMembershipModel/GetByUserId/Query: " + err.Error())
	}
	return results, nil
}

func (m RoomMembershipModel) GetGroupsByUserId(ctx context.Context, db *gorm.DB, timeout time.Duration, limit, page int) ([]RoomMembershipModel, error) {
	var results []RoomMembershipModel
	err := db.Table(os.Getenv("POSTGRESQLMEMBERSHIPTABLE")).Where(`"user_id" = '` + m.UserId + "'" + ` AND "leave_at" = '0' AND "room_id" LIKE 'group-%'`).Limit(limit).Offset(limit * page).Find(&results).Error
	if err != nil {
		log.Println("RoomMembershipModel/GetGroupsByUserId/Query: ", err)
		return []RoomMembershipModel{}, errors.New("Err in RoomMembershipModel/GetGroupsByUserId/Query: " + err.Error())
	}
	return results, nil
}

func (m RoomMembershipModel) GetOtherUserInMutualRoom(ctx context.Context, db *gorm.DB, timeout time.Duration) (RoomMembershipModel, error) {
	log.Println("RoomMembershipModel/GetOtherUserInMutualRoom/Query: ", m.UserId, m.RoomId)
	var result RoomMembershipModel
	err := db.Table(os.Getenv("POSTGRESQLMEMBERSHIPTABLE")).Where(`"room_id" = ? AND "user_id" != ?`, m.RoomId, m.UserId).Find(&result).Error
	if err != nil {
		return RoomMembershipModel{}, errors.New("Err in RoomMembershipModel/GetOtherUserInMutualRoom/Query: scan error")
	}
	return result, nil
}
