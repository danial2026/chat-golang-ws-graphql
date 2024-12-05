package model

import (
	"context"

	"chat-golang-react/chat/graphql/resources"
	models "chat-golang-react/chat/models"

	"log"
)

func GetAllRoomMembersByRoomId(ctx context.Context, res resources.Resources, roomId string, limit, page int) ([]models.RoomMembershipModel, error) {
	member := models.RoomMembershipModel{
		RoomId: roomId,
	}

	result, err := member.GetByRoomIdWithPagination(ctx, res.ROOMSDB, res.Timeout, limit, page)
	if err != nil {
		log.Println("err in GetAllRoomsByUserId error is: ", err.Error())
		return []models.RoomMembershipModel{}, err
	}

	return result, nil
}


func GetAllGroupRoomsByUserId(ctx context.Context, res resources.Resources, userId string, limit, page int) ([]models.RoomMembershipModel, error) {
	member := models.RoomMembershipModel{
		UserId: userId,
	}

	result, err := member.GetGroupsByUserId(ctx, res.ROOMSDB, res.Timeout, limit, page)
	if err != nil {
		log.Println("err in GetAllRoomsByUserId error is: ", err.Error())
		return []models.RoomMembershipModel{}, err
	}

	return result, nil
}


func GetAllRoomsByUserId(ctx context.Context, res resources.Resources, userId string, limit, page int) ([]models.RoomMembershipModel, error) {
	member := models.RoomMembershipModel{
		UserId: userId,
	}

	result, err := member.GetByUserId(ctx, res.ROOMSDB, res.Timeout, limit, page)
	if err != nil {
		log.Println("err in GetAllRoomsByUserId error is: ", err.Error())
		return []models.RoomMembershipModel{}, err
	}

	return result, nil
}

func IsRoomMember(ctx context.Context, res resources.Resources, username, userId, roomId string) (bool, error) {
	member := models.RoomMembershipModel{
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

func GetRoomMembership(ctx context.Context, res resources.Resources, username, userId, roomId string) (models.RoomMembershipModel, error) {
	member := models.RoomMembershipModel{
		Username: username,
		UserId:   userId,
		RoomId:   roomId,
	}
	log.Println("GetRoomMembership/RoomMembershipModel: ", member)

	membership, err := member.GetItemByUserIdAndRoomId(ctx, res.ROOMSDB, res.Timeout)
	if err != nil {
		log.Println("err in GetRoomMembership/Exists: ", err.Error())
		return member, err
	}

	return membership, nil
}

func GetOtherUserInMutualRoom(ctx context.Context, res resources.Resources, userId, roomId string) (models.RoomMembershipModel, error) {
	member := models.RoomMembershipModel{
		UserId: userId,
		RoomId: roomId,
	}

	result, err := member.GetOtherUserInMutualRoom(ctx, res.ROOMSDB, res.Timeout)
	if err != nil {
		log.Println("err in GetOtherUserInMutualRoom error is: ", err.Error())
		return models.RoomMembershipModel{}, err
	}

	return result, nil
}
