package models

import (
	"context"
	"errors"

	"chat-golang-react/chat/websocket/resources"

	"log"
)

func IsRoomExists(ctx context.Context, res resources.Resource, roomId string) (bool, error) {
	room := RoomModel{
		Id: roomId,
	}

	_, ok, err := room.Exists(ctx, res.ROOMSDB, res.Timeout)
	if err != nil {
		log.Println("err in IsRoomExists error is: ", err.Error())
		return false, err
	}

	return ok, nil
}

// TODO : use to create room later
func SyncMutualRoom(ctx context.Context, res resources.Resource, roomId string) error {
	// room, getErr := syncers.GetRoomInfo(res.User.AccessToken, roomId)
	// if getErr != nil {
	// 	log.Println("err in SyncRoom/syncers.GetRoomInfo: " + getErr.Error())
	// 	return getErr
	// }

	_, createErr := CreateMutualRoom(ctx, res, roomId, res.User.Username, res.User.UserId, "room.UserName", "room.UserID")
	if createErr != nil {
		return createErr
	}

	return nil
}

func CreateMutualRoom(ctx context.Context, res resources.Resource, roomId, RequestedUser, RequestedUserId, TargetUser, TargetUserId string) (RoomModel, error) {
	room := RoomModel{
		Id:   roomId,
		Type: RoomMutualType,
	}

	memberOne := RoomMembershipModel{
		RoomId:   roomId,
		JoinBy:   RequestedUser,
		JoinById: RequestedUserId,
		Username: RequestedUser,
		UserId:   RequestedUserId,
		IsAdmin:  false,
	}
	memberTwo := RoomMembershipModel{
		RoomId:   roomId,
		JoinBy:   TargetUser,
		JoinById: TargetUserId,
		Username: TargetUser,
		UserId:   TargetUserId,
		IsAdmin:  false,
	}

	newRoom, roomErr := room.Create(ctx, res.ROOMSDB, res.Timeout)
	if roomErr != nil {
		log.Println("err in CreateMutualRoom/room.ToNewItem: " + roomErr.Error())
		return room, roomErr
	}
	memberOneModel, memberOneErr := memberOne.ToNewStruct()
	if memberOneErr != nil {
		log.Println("err in CreateMutualRoom/memberOne.ToNewItem: " + memberOneErr.Error())
		return room, memberOneErr
	}
	memberTwoModel, memberTwoErr := memberTwo.ToNewStruct()
	if memberTwoErr != nil {
		log.Println("err in CreateMutualRoom/memberTwo.ToNewItem: " + memberTwoErr.Error())
		return room, memberTwoErr
	}

	isRoomExist, existErr := IsRoomExists(ctx, res, roomId)
	if existErr != nil {
		log.Println("CreateMutualRoom/IsRoomExists error is: ", existErr.Error())
		return room, existErr
	}
	if !isRoomExist {
		_, creationErr := newRoom.Create(ctx, res.ROOMSDB, res.Timeout)
		if creationErr != nil {
			return room, errors.New("Err in CreateMutualRoom/CreateRoom: " + creationErr.Error())
		}
	}

	isRoomExist, existErr = IsRoomMember(ctx, res, memberOneModel.Username, memberOneModel.UserId, roomId)
	if existErr != nil {
		log.Println("CreateMutualRoom/IsRoomExists error is: ", existErr.Error())
		return room, existErr
	}
	if !isRoomExist {
		_, creationErr := memberOneModel.Create(ctx, res.ROOMSDB, res.Timeout)
		if creationErr != nil {
			return newRoom, errors.New("Err in CreateMutualRoom/CreateMembership 1: " + creationErr.Error())
		}
	}

	isRoomExist, existErr = IsRoomMember(ctx, res, memberTwoModel.Username, memberTwoModel.UserId, roomId)
	if existErr != nil {
		log.Println("CreateMutualRoom/IsRoomExists error is: ", existErr.Error())
		return room, existErr
	}
	if !isRoomExist {
		_, creationErr := memberTwoModel.Create(ctx, res.ROOMSDB, res.Timeout)
		if creationErr != nil {
			return newRoom, errors.New("Err in CreateMutualRoom/CreateMembership 2: " + creationErr.Error())
		}
	}

	return newRoom, nil
}

func DeleteRoomForMe(ctx context.Context, res resources.Resource, messageId, roomId, username, userId string) (RoomModel, error) {
	room := RoomModel{
		Id: messageId,
	}
	err := room.DeleteItemFor(ctx, res.ROOMSDB, res.Timeout, DeletedForMe)
	if err != nil {
		return room, err
	}
	return room, nil
}

func DeleteRoomForAll(ctx context.Context, res resources.Resource, messageId, roomId, username, userId string) (RoomModel, error) {
	room := RoomModel{
		Id: messageId,
	}
	err := room.DeleteItemFor(ctx, res.ROOMSDB, res.Timeout, DeletedForAll)
	if err != nil {
		return room, err
	}
	return room, nil
}
