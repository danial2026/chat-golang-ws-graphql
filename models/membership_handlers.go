package models

import (
	"context"
	"errors"
	"log"

	"chat-golang-react/chat/websocket/resources"
)

func IsRoomMember(ctx context.Context, res resources.Resource, username, userId, roomId string) (bool, error) {
	member := RoomMembershipModel{
		Username: username,
		UserId:   userId,
		RoomId:   roomId,
	}
	log.Println("IsRoomMember/RoomMembershipModel: ", member)

	ok, err := member.Exists(ctx, res.ROOMSDB, res.Timeout)
	if err != nil {
		log.Println("err in IsRoomMember/Exists: ", err.Error())
		return false, err
	}

	return ok, nil
}

func GetMembersByRoomId(ctx context.Context, res resources.Resource, roomId string) ([]RoomMembershipModel, error) {
	member := RoomMembershipModel{RoomId: roomId}
	members, queryErr := member.GetByRoomId(ctx, res.ROOMSDB, res.Timeout)
	if queryErr != nil {
		return []RoomMembershipModel{}, queryErr
	}

	return members, nil
}

func GetMutualRoomContactUsers(ctx context.Context, res resources.Resource, roomId string) ([]RoomMembershipModel, error) {

	membership := RoomMembershipModel{
		RoomId: roomId,
	}

	members, memberErr := membership.GetByRoomId(ctx, res.ROOMSDB, res.Timeout)
	if memberErr != nil {
		return []RoomMembershipModel{}, errors.New("err in GetMutualRoomContactUser/GetByRoomId: " + memberErr.Error())
	}

	if len(members) == 0 {
		return []RoomMembershipModel{}, errors.New("ContactUser membership not found")
	}

	// for _, m := range members {
	// 	if m.Username != res.User.Username {
	// 		return m, nil
	// 	}
	// }

	return members, nil
}
