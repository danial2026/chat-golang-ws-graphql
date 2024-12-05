package model

import (
	"context"
	"errors"

	"chat-golang-react/chat/graphql/resources"
	models "chat-golang-react/chat/models"

	"log"
)

func GetRoomById(ctx context.Context, res resources.Resources, roomId string) (models.RoomModel, error) {
	room := models.RoomModel{
		Id: roomId,
	}

	result, err := room.GetItem(ctx, res.ROOMSDB, res.Timeout)
	if err != nil {
		log.Println("err in IsRoomExists error is: ", err.Error())
		return models.RoomModel{}, err
	}

	return result, nil
}

func IsRoomExists(ctx context.Context, res resources.Resources, roomId string) (bool, error) {
	room := models.RoomModel{
		Id: roomId,
	}

	_, ok, err := room.Exists(ctx, res.ROOMSDB, res.Timeout)
	if err != nil {
		log.Println("err in IsRoomExists error is: ", err.Error())
		return false, err
	}

	return ok, nil
}

func CreateMutualRoom(ctx context.Context, res resources.Resources, RequestedUser, RequestedUserId, TargetUser, TargetUserId string) (models.RoomModel, error) {
	room := models.RoomModel{
		Title: RequestedUser + ", " + TargetUser,
		Type:  models.RoomMutualType,
		Id:    RequestedUserId + "-" + TargetUserId,
	}

	newRoom, roomErr := room.Create(ctx, res.ROOMSDB, res.Timeout)
	if roomErr != nil {
		log.Println("err in CreateMutualRoom/room.ToNewItem: " + roomErr.Error())
		return room, roomErr
	}

	memberOne := models.RoomMembershipModel{
		RoomId:   newRoom.Id,
		JoinBy:   RequestedUser,
		JoinById: RequestedUserId,
		Username: RequestedUser,
		UserId:   RequestedUserId,
		IsAdmin:  false,
	}
	memberTwo := models.RoomMembershipModel{
		RoomId:   newRoom.Id,
		JoinBy:   RequestedUser,
		JoinById: RequestedUserId,
		Username: TargetUser,
		UserId:   TargetUserId,
		IsAdmin:  false,
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

	isRoomMembershipExist, existErr := IsRoomMember(ctx, res, memberOneModel.Username, memberOneModel.UserId, newRoom.Id)
	if existErr != nil {
		log.Println("CreateMutualRoom/isRoomMembershipExist error is: ", existErr.Error())
		return room, existErr
	}
	if !isRoomMembershipExist {
		_, creationErr := memberOneModel.Create(ctx, res.ROOMSDB, res.Timeout)
		if creationErr != nil {
			return newRoom, errors.New("Err in CreateMutualRoom/CreateMembership 1: " + creationErr.Error())
		}
	}

	isRoomMembershipExist, existErr = IsRoomMember(ctx, res, memberTwoModel.Username, memberTwoModel.UserId, newRoom.Id)
	if existErr != nil {
		log.Println("CreateMutualRoom/isRoomMembershipExist error is: ", existErr.Error())
		return room, existErr
	}
	if !isRoomMembershipExist {
		_, creationErr := memberTwoModel.Create(ctx, res.ROOMSDB, res.Timeout)
		if creationErr != nil {
			return newRoom, errors.New("Err in CreateMutualRoom/CreateMembership 2: " + creationErr.Error())
		}
	}

	return newRoom, nil
}

// 1. create the room
// 2. loop thro the list of users and create a membership for them
func CreateGroupRoom(ctx context.Context, res resources.Resources, RequestedUser, RequestedUserId, roomTitle string, TargetUsers []*UserInput) (models.RoomModel, error) {

	var newRoomTitle = ""
	if roomTitle != "" {
		newRoomTitle = roomTitle
	} else {
		if len(TargetUsers) == 1 {
			newRoomTitle = RequestedUser + ", " + TargetUsers[0].Username
		} else {
			newRoomTitle = RequestedUser + ", " + TargetUsers[0].Username + " and others"
		}
	}

	room := models.RoomModel{
		Title:    newRoomTitle,
		Type:     models.RoomGroupType,
		IsPublic: true,
	}

	newRoom, roomErr := room.Create(ctx, res.ROOMSDB, res.Timeout)
	if roomErr != nil {
		log.Println("err in CreateGroupRoom/room.ToNewItem: " + roomErr.Error())
		return room, roomErr
	}

	memberOne := models.RoomMembershipModel{
		RoomId:   newRoom.Id,
		JoinBy:   RequestedUser,
		JoinById: RequestedUserId,
		Username: RequestedUser,
		UserId:   RequestedUserId,
		IsAdmin:  true,
	}

	memberOneModel, memberOneErr := memberOne.ToNewStruct()
	if memberOneErr != nil {
		log.Println("err in CreateGroupRoom/memberOne.ToNewItem: " + memberOneErr.Error())
		return room, memberOneErr
	}

	isRoomMembershipExist, existErr := IsRoomMember(ctx, res, memberOneModel.Username, memberOneModel.UserId, newRoom.Id)
	if existErr != nil {
		log.Println("CreateGroupRoom/isRoomMembershipExist error is: ", existErr.Error())
		return room, existErr
	}
	if !isRoomMembershipExist {
		_, creationErr := memberOneModel.Create(ctx, res.ROOMSDB, res.Timeout)
		if creationErr != nil {
			return newRoom, errors.New("Err in CreateGroupRoom/CreateMembership 1: " + creationErr.Error())
		}
	}

	for i := 0; i < len(TargetUsers); i++ {
		targetUser := TargetUsers[i]

		memberTwo := models.RoomMembershipModel{
			RoomId:   newRoom.Id,
			JoinBy:   RequestedUser,
			JoinById: RequestedUserId,
			Username: targetUser.Username,
			UserId:   targetUser.ID,
			IsAdmin:  false,
		}

		memberTwoModel, memberTwoErr := memberTwo.ToNewStruct()
		if memberTwoErr != nil {
			log.Println("err in CreateGroupRoom/memberTwo.ToNewItem: " + memberTwoErr.Error())
			return room, memberTwoErr
		}

		isRoomMembershipExist, existErr = IsRoomMember(ctx, res, memberTwoModel.Username, memberTwoModel.UserId, newRoom.Id)
		if existErr != nil {
			log.Println("CreateGroupRoom/isRoomMembershipExist error is: ", existErr.Error())
			return room, existErr
		}
		if !isRoomMembershipExist {
			_, creationErr := memberTwoModel.Create(ctx, res.ROOMSDB, res.Timeout)
			if creationErr != nil {
				return newRoom, errors.New("Err in CreateMutualRoom/CreateMembership : " + creationErr.Error())
			}
		}
	}

	return newRoom, nil
}

func LeavePublicGroupRoom(ctx context.Context, res resources.Resources, RoomId, TargetUser, TargetUserId string) (models.RoomMembershipModel, error) {
	response := models.RoomMembershipModel{}

	log.Println("not implemented")

	return response, nil
}

func RemoveMemberFromPublicGroupRoom(ctx context.Context, res resources.Resources, RoomId, RequestedUser, RequestedUserId, TargetUser, TargetUserId string) (models.RoomMembershipModel, error) {
	response := models.RoomMembershipModel{}

	log.Println("not implemented")

	return response, nil
}

func AddToPublicGroupRoom(ctx context.Context, res resources.Resources, RoomId string, TargetUsers []*UserInput) ([]models.RoomMembershipModel, error) {
	response := []models.RoomMembershipModel{}

	isRoomExist, existErr := IsRoomExists(ctx, res, RoomId)
	if existErr != nil {
		log.Println("AddToPublicGroupRoom/IsRoomExist error is: ", existErr.Error())
		return response, existErr
	}
	if !isRoomExist {
		var errMsg = "Room not found"
		return response, errors.New(errMsg)
	}

	room := models.RoomModel{
		Id: RoomId,
	}

	roomModel, roomErr := room.GetItem(ctx, res.ROOMSDB, res.Timeout)
	if roomErr != nil {
		log.Println("err in JoinPublicGroupRoom/room.GetItem: " + roomErr.Error())
		return response, roomErr
	}

	/*
		check if room is Active , Public and type Group
	*/
	if roomModel.Type != models.RoomGroupType {
		var errMsg = "Room isn't a Group"
		return response, errors.New(errMsg)
	}
	if !roomModel.IsPublic {
		var errMsg = "Room isn't Public"
		return response, errors.New(errMsg)
	}
	if !roomModel.IsActive {
		var errMsg = "Room isn't Active"
		return response, errors.New(errMsg)
	}

	for i := 0; i < len(TargetUsers); i++ {
		targetUser := TargetUsers[i]

		memberTwo := models.RoomMembershipModel{
			RoomId:   RoomId,
			JoinBy:   res.User.Username,
			JoinById: res.User.UserId,
			Username: targetUser.Username,
			UserId:   targetUser.ID,
			IsAdmin:  false,
		}

		memberTwoModel, memberTwoErr := memberTwo.ToNewStruct()
		if memberTwoErr != nil {
			log.Println("err in AddToPublicGroupRoom/memberTwo.ToNewItem: " + memberTwoErr.Error())
			return response, memberTwoErr
		}

		isRoomMembershipExist, existErr := IsRoomMember(ctx, res, memberTwoModel.Username, memberTwoModel.UserId, RoomId)
		if existErr != nil {
			log.Println("AddToPublicGroupRoom/isRoomMembershipExist error is: ", existErr.Error())
			return response, existErr
		}
		if !isRoomMembershipExist {
			newMembership, creationErr := memberTwoModel.Create(ctx, res.ROOMSDB, res.Timeout)
			if creationErr != nil {
				return response, errors.New("Err in AddToPublicGroupRoom/CreateMembership: " + creationErr.Error())
			}
			response = append(response, newMembership)
		}
	}
	return response, nil
}

func JoinPublicGroupRoom(ctx context.Context, res resources.Resources, RoomId, RequestedUser, RequestedUserId string) (bool, error) {

	isRoomExist, existErr := IsRoomExists(ctx, res, RoomId)
	if existErr != nil {
		log.Println("JoinPublicGroupRoom/IsRoomExist error is: ", existErr.Error())
		return false, existErr
	}
	if !isRoomExist {
		var errMsg = "Room not found"
		return false, errors.New(errMsg)
	}

	room := models.RoomModel{
		Id: RoomId,
	}

	roomModel, roomErr := room.GetItem(ctx, res.ROOMSDB, res.Timeout)
	if roomErr != nil {
		log.Println("err in JoinPublicGroupRoom/room.GetItem: " + roomErr.Error())
		return false, roomErr
	}

	/*
		check if room is Active , Public and type Group
	*/
	if roomModel.Type != models.RoomGroupType {
		var errMsg = "Room isn't a Group"
		return false, errors.New(errMsg)
	}
	if !roomModel.IsPublic {
		var errMsg = "Room isn't Public"
		return false, errors.New(errMsg)
	}
	if !roomModel.IsActive {
		var errMsg = "Room isn't Active"
		return false, errors.New(errMsg)
	}

	member := models.RoomMembershipModel{
		RoomId:   RoomId,
		JoinBy:   RequestedUser,
		JoinById: RequestedUserId,
		Username: RequestedUser,
		UserId:   RequestedUserId,
		IsAdmin:  false,
	}

	memberModel, memberErr := member.ToNewStruct()
	if memberErr != nil {
		log.Println("err in JoinPublicGroupRoom/member.ToNewItem: " + memberErr.Error())
		return false, memberErr
	}

	isRoomMembershipExist, existErr := IsRoomMember(ctx, res, memberModel.Username, memberModel.UserId, RoomId)
	if existErr != nil {
		log.Println("JoinPublicGroupRoom/isRoomMembershipExist error is: ", existErr.Error())
		return false, existErr
	}
	if !isRoomMembershipExist {
		var errMsg = "Room not found"
		return false, errors.New(errMsg)
	}

	return true, nil
}

func GenerateRoomTitle(ctx context.Context, res resources.Resources, RoomModel models.RoomModel) string {

	var newRoomTitle = ""
	if RoomModel.Type != models.RoomMutualType {
		newRoomTitle = RoomModel.Title
	} else {
		newMemberships := models.RoomMembershipModel{
			RoomId: RoomModel.Id,
		}
		result, existErr := newMemberships.GetByRoomId(ctx, res.ROOMSDB, resources.DefaultTimeout)
		if existErr != nil {
			log.Println("GenerateRoomTitle error is: ", existErr.Error())
			return newRoomTitle
		}
		for i := 0; i < len(result); i++ {
			if result[i].UserId != res.User.UserId {
				newRoomTitle = result[i].Username
				break
			}
		}
	}

	return newRoomTitle
}
