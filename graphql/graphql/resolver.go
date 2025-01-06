package graphql

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"

	"chat-golang-react/chat/graphql/auth"

	common "chat-golang-react/chat/common"
	graphqlmodel "chat-golang-react/chat/graphql/models"
	graphqlresources "chat-golang-react/chat/graphql/resources"
	models "chat-golang-react/chat/models"
)

func authHandler(ctx context.Context, r Resolver) (graphqlresources.Resources, error) {
	// get token from context provided by Auth Middleware
	resAuth := ctx.Value(auth.UserCtxKey)

	if resAuth == nil {
		// no token found
		msgErr := "permission denied"
		return graphqlresources.Resources{}, errors.New(msgErr)
	}

	// token `any` to map[string]
	resultAuth := make(map[string]string)
	mapstructure.Decode(resAuth, &resultAuth)

	if resultAuth["user_id"] == "" || resultAuth["username"] == "" {
		// no token found
		msgErr := "permission denied"
		return graphqlresources.Resources{}, errors.New(msgErr)
	}

	newResources := graphqlresources.Resources{
		ROOMSDB: r.ROOMSDB,
		USERDB:  r.USERDB,
		Timeout: r.Timeout,
		User: graphqlresources.User{
			UserId:   resultAuth["user_id"],
			Username: resultAuth["username"],
		},
	}
	return newResources, nil
}

type Resolver struct {
	// rooms and memberships postgresql database
	ROOMSDB *gorm.DB
	// client user mongodb database
	USERDB *mongo.Client
	// connection timeout
	Timeout time.Duration
}

func (r *mutationResolver) CreateRoom(ctx context.Context, input *graphqlmodel.CreateRoomInput) (*graphqlmodel.RoomResponse, error) {
	response := &graphqlmodel.RoomResponse{}

	newResources, errAuth := authHandler(ctx, *r.Resolver)
	if errAuth != nil {
		return nil, errAuth
	}

	switch input.RoomType.String() {
	case graphqlmodel.RoomTypeMutual.String():
		{
			if len(input.Users) != 1 {
				msgErr := "mutual room needs exactly one user"
				response.Error = &graphqlmodel.Error{
					Message: &msgErr,
				}
				return response, errors.New(msgErr)
			}
			newRoom, createErr := graphqlmodel.CreateMutualRoom(ctx, newResources, newResources.User.Username, newResources.User.UserId, input.Users[0].Username, input.Users[0].ID)
			if createErr != nil {
				msgErr := "couldn't create the room"
				response.Error = &graphqlmodel.Error{
					Message: &msgErr,
				}
				return response, errors.New(msgErr)
			}
			response.Data = convertToRoomGraphModel(newRoom, true, true)
		}
	case graphqlmodel.RoomTypeGroup.String():
		{
			// TODO: for now we let user create empty groups
			// if len(input.Users) == 0 {
			// 	msgErr := "Users list cant be empty"
			// 	response.Error = &graphqlmodel.Error{
			// 		Message: &msgErr,
			// 	}
			// 	return response, errors.New(msgErr)
			// }
			newRoom, createErr := graphqlmodel.CreateGroupRoom(ctx, newResources, newResources.User.Username, newResources.User.UserId, *input.Title, input.Users)
			if createErr != nil {
				msgErr := "couldn't create the room"
				response.Error = &graphqlmodel.Error{
					Message: &msgErr,
				}
				return response, errors.New(msgErr)
			}

			response.Data = convertToRoomGraphModel(newRoom, true, true)
		}
	case graphqlmodel.RoomTypeChannel.String():
		{
			// TODO : create channel
			msgErr := "channels are not implemented yet"
			response.Error = &graphqlmodel.Error{
				Message: &msgErr,
			}
			return response, errors.New(msgErr)
		}
	}

	return response, nil
}

func (r *mutationResolver) AddRoomMembers(ctx context.Context, input *graphqlmodel.AddRoomMembersInput) (*graphqlmodel.RoomMembershipResponse, error) {
	response := &graphqlmodel.RoomMembershipResponse{}

	newResources, errAuth := authHandler(ctx, *r.Resolver)
	if errAuth != nil {
		return nil, errAuth
	}

	if len(input.Users) == 0 {
		msgErr := "field users is required"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	isRoomExist, existErr := graphqlmodel.IsRoomExists(ctx, newResources, input.RoomID)
	if existErr != nil {
		log.Println("CreateMutualRoom/IsRoomExist error is: ", existErr.Error())
		var msgErr = "can't find the room"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}
	if !isRoomExist {
		var msgErr = "Room not found"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	// check if user is a room member
	isMember, err := graphqlmodel.IsRoomMember(ctx, newResources, newResources.User.Username, newResources.User.UserId, input.RoomID)
	if err != nil {
		log.Println("err in GetRoom error is: ", err.Error())
		return nil, err
	}
	if !isMember {
		msgErr := "permission denied"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	newRoomModel, err := graphqlmodel.AddToPublicGroupRoom(ctx, newResources, input.RoomID, input.Users)
	if err != nil {
		log.Println("err in MembershipToRoom error is: ", err.Error())
		return response, err
	}

	response.Data = convertToMembershipGraphModel(newRoomModel)
	return response, nil
}

func (r *mutationResolver) RemoveRoomMember(ctx context.Context, input *graphqlmodel.RemoveRoomMemberInput) (*graphqlmodel.RoomMembershipResponse, error) {
	response := &graphqlmodel.RoomMembershipResponse{}

	newResources, errAuth := authHandler(ctx, *r.Resolver)
	if errAuth != nil {
		return nil, errAuth
	}

	isRoomExist, existErr := graphqlmodel.IsRoomExists(ctx, newResources, input.RoomID)
	if existErr != nil {
		log.Println("CreateMutualRoom/IsRoomExist error is: ", existErr.Error())
		var msgErr = "can't find the room"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}
	if !isRoomExist {
		var msgErr = "Room not found"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	// check if user is a room member
	isMember, err := graphqlmodel.IsRoomMember(ctx, newResources, newResources.User.Username, newResources.User.UserId, input.RoomID)
	if err != nil {
		log.Println("err in GetRoom error is: ", err.Error())
		return nil, err
	}
	if !isMember {
		msgErr := "permission denied"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	// check if user is admin
	mmbership, err := graphqlmodel.GetRoomMembership(ctx, newResources, newResources.User.Username, newResources.User.UserId, input.RoomID)
	if err != nil {
		log.Println("err in GetRoomMembership error is: ", err.Error())
		return nil, err
	}
	if !mmbership.IsAdmin {
		msgErr := "permission denied, only an admin can do this"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	newRoomModel, err := graphqlmodel.RemoveMemberFromPublicGroupRoom(ctx, newResources, input.RoomID, newResources.User.Username, newResources.User.UserId, input.User.Username, input.User.ID)
	if err != nil {
		log.Println("err in MembershipToRoom error is: ", err.Error())
		return response, err
	}

	response.Data = convertToMembershipGraphModel([]models.RoomMembershipModel{newRoomModel})
	return response, nil
}

func (r *mutationResolver) JoinRoom(ctx context.Context, input *graphqlmodel.JoinRoomInput) (*graphqlmodel.RoomMembershipResponse, error) {
	response := &graphqlmodel.RoomMembershipResponse{}

	newResources, errAuth := authHandler(ctx, *r.Resolver)
	if errAuth != nil {
		return nil, errAuth
	}

	isRoomExist, existErr := graphqlmodel.IsRoomExists(ctx, newResources, input.RoomID)
	if existErr != nil {
		log.Println("CreateMutualRoom/IsRoomExist error is: ", existErr.Error())
		var msgErr = "can't find the room"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}
	if !isRoomExist {
		var msgErr = "Room not found"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	// TODO: check if room is mutual

	// check if user is a room member
	isMember, err := graphqlmodel.IsRoomMember(ctx, newResources, newResources.User.Username, newResources.User.UserId, input.RoomID)
	if err != nil {
		log.Println("err in GetRoom error is: ", err.Error())
		return nil, err
	}
	if isMember {
		msgErr := "already joined"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	newMembership, err := graphqlmodel.AddToPublicGroupRoom(ctx, newResources, input.RoomID, []*graphqlmodel.UserInput{
		{
			ID:       newResources.User.UserId,
			Username: newResources.User.Username,
		},
	})
	if err != nil {
		log.Println("err in MembershipToRoom error is: ", err.Error())
		return response, err
	}

	response.Data = convertToMembershipGraphModel(newMembership)
	return response, nil
}

func (r *mutationResolver) LeaveRoom(ctx context.Context, input *graphqlmodel.LeaveRoomInput) (*graphqlmodel.RoomMembershipResponse, error) {
	response := &graphqlmodel.RoomMembershipResponse{}

	newResources, errAuth := authHandler(ctx, *r.Resolver)
	if errAuth != nil {
		return nil, errAuth
	}

	isRoomExist, existErr := graphqlmodel.IsRoomExists(ctx, newResources, input.RoomID)
	if existErr != nil {
		log.Println("CreateMutualRoom/IsRoomExist error is: ", existErr.Error())
		var msgErr = "can't find the room"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}
	if !isRoomExist {
		var msgErr = "Room not found"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	// check if user is a room member
	isMember, err := graphqlmodel.IsRoomMember(ctx, newResources, newResources.User.Username, newResources.User.UserId, input.RoomID)
	if err != nil {
		log.Println("err in GetRoom error is: ", err.Error())
		return nil, err
	}
	if !isMember {
		msgErr := "not a member"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	mmbership, err := graphqlmodel.LeavePublicGroupRoom(ctx, newResources, newResources.User.Username, newResources.User.UserId, input.RoomID)
	if err != nil {
		log.Println("err in GetRoomMembership error is: ", err.Error())
		return nil, err
	}

	response.Data = convertToMembershipGraphModel([]models.RoomMembershipModel{mmbership})
	return response, nil
}

func (r *queryResolver) GetRooms(ctx context.Context, pagination *graphqlmodel.Pagination, room_type *graphqlmodel.RoomType) (*graphqlmodel.RoomSummeryResponse, error) {

	response := &graphqlmodel.RoomSummeryResponse{}

	newResources, errAuth := authHandler(ctx, *r.Resolver)
	if errAuth != nil {
		return nil, errAuth
	}

	var limit = 10
	var page = 0

	if pagination != nil {
		// response based on the Pagination
		if pagination.Limit != 0 && pagination.Limit < 50 {
			limit = pagination.Limit
		}
		if pagination.Page != 0 && pagination.Page < 100 {
			page = pagination.Page
		}
	}

	// return all rooms that user is in
	var memberships []models.RoomMembershipModel
	var err error
	if room_type != nil {
		if room_type.String() == graphqlmodel.RoomTypeGroup.String() {
			memberships, err = graphqlmodel.GetAllGroupRoomsByUserId(ctx, newResources, newResources.User.UserId, limit, page)

		}
	}
	if memberships == nil {
		memberships, err = graphqlmodel.GetAllRoomsByUserId(ctx, newResources, newResources.User.UserId, limit, page)
	}
	if err != nil {
		log.Println("err in GetAllRoomsByUserId error is: ", err.Error())
		msgErr := "no rooms found"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}
	rooms, err := membershipToRoom(ctx, newResources, memberships)
	if err != nil {
		log.Println("err in GetAllRoomsByUserId error is: ", err.Error())
		msgErr := "no rooms found"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	for roomIndex := range rooms {
		if rooms[roomIndex].RoomType.String() == strings.ToLower(graphqlmodel.RoomTypeMutual.String()) {
			otherMember, err := graphqlmodel.GetOtherUserInMutualRoom(ctx, newResources, newResources.User.UserId, rooms[roomIndex].ID)
			if err != nil {
				log.Println("err in GetOtherUserInMutualRoom error is: ", err.Error())
				return response, err
			}
			user := models.User{}
			err = user.GetByID(ctx, newResources.USERDB, otherMember.UserId)
			if err != nil {
				log.Println("err in GetRooms/USERDB/GetByID: ", err.Error())
				return response, err
			}
			rooms[roomIndex].Title = user.Fullname
			if len(user.Images) > 0 {
				rooms[roomIndex].Image = &user.Images[0]
			}
		}
	}

	response.Data = rooms
	return response, nil
}

func (r *queryResolver) GetRoom(ctx context.Context, id string) (*graphqlmodel.RoomResponse, error) {
	response := &graphqlmodel.RoomResponse{}

	newResources, errAuth := authHandler(ctx, *r.Resolver)
	if errAuth != nil {
		return nil, errAuth
	}

	if id == "" {
		msgErr := "id is required"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	var roomIds [2]string

	//  no need to validate the id anymore
	// validate room id format:
	roomId, invalidErr := common.GetAndValidateUUID(id)
	if invalidErr != nil {
		var msgErr = "Invalid room_id"
		log.Println(msgErr)
		containsDash := strings.Contains(id, "-")
		if !containsDash {
			roomIds[0] = id + "-" + newResources.User.UserId
			roomIds[1] = newResources.User.UserId + "-" + id
		} else if strings.HasPrefix(id, "group") {
			roomIds[0] = id
		}
	} else {
		roomIds[0] = roomId
	}

	for roomIndex := range roomIds {

		// check if user is a room member
		isMember, err := graphqlmodel.IsRoomMember(ctx, newResources, newResources.User.Username, newResources.User.UserId, roomIds[roomIndex])
		if err != nil {
			log.Println("err in GetRoom error is: ", err.Error())
			return response, err
		}
		if !isMember {
			msgErr := "permission denied"
			if len(roomIds) == roomIndex {
				response.Error = &graphqlmodel.Error{
					Message: &msgErr,
				}
				return response, errors.New(msgErr)
			} else {
				log.Println(msgErr)
				roomIndex = roomIndex + 1
			}
		}

		room, err := graphqlmodel.GetRoomById(ctx, newResources, roomIds[roomIndex])
		if err != nil {
			log.Println("err in GetRoom error is: ", err.Error())
			return response, err
		}

		// check if user is admin
		mmbership, err := graphqlmodel.GetRoomMembership(ctx, newResources, newResources.User.Username, newResources.User.UserId, room.Id)
		if err != nil {
			log.Println("err in GetRoomMembership error is: ", err.Error())
			return nil, err
		}

		if strings.ToTitle(room.Type) == graphqlmodel.RoomTypeMutual.String() {
			otherMember, err := graphqlmodel.GetOtherUserInMutualRoom(ctx, newResources, newResources.User.UserId, roomIds[roomIndex])
			if err != nil {
				log.Println("err in GetOtherUserInMutualRoom error is: ", err.Error())
				return response, err
			}
			user := models.User{}
			err = user.GetByID(ctx, newResources.USERDB, otherMember.UserId)
			if err != nil {
				log.Println("err in GetRoom/USERDB/GetByID: ", err.Error())
				return response, err
			}
			room.Title = user.Fullname
			if len(user.Images) > 0 {
				room.Image = user.Images[0]
			}
		}

		log.Println("response GetRoom/room: ", room)
		response.Data = convertToRoomGraphModel(room, isMember, mmbership.IsAdmin)
		return response, nil
	}
	return response, nil
}

func (r *queryResolver) GetRoomMembers(ctx context.Context, roomID string, pagination *graphqlmodel.Pagination) (*graphqlmodel.RoomMembershipResponse, error) {
	response := &graphqlmodel.RoomMembershipResponse{}

	newResources, errAuth := authHandler(ctx, *r.Resolver)
	if errAuth != nil {
		return nil, errAuth
	}

	var limit = 10
	var page = 0

	if pagination != nil {
		// response based on the Pagination
		if pagination.Limit != 0 && pagination.Limit < 50 {
			limit = pagination.Limit
		}
		if pagination.Page != 0 && pagination.Page < 100 {
			page = pagination.Page
		}
	}

	// check if user is a room member
	isMember, err := graphqlmodel.IsRoomMember(ctx, newResources, newResources.User.Username, newResources.User.UserId, roomID)
	if err != nil {
		log.Println("err in GetRoom error is: ", err.Error())
		return nil, err
	}
	if !isMember {
		msgErr := "permission denied"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	// return all members
	memberships, err := graphqlmodel.GetAllRoomMembersByRoomId(ctx, newResources, roomID, limit, page)
	if err != nil {
		log.Println("err in GetAllRoomsByUserId error is: ", err.Error())
		msgErr := "no rooms found"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	response.Data = convertToMembershipGraphModel(memberships)
	return response, nil
}

func (r *queryResolver) GetMessages(ctx context.Context, roomID string, lastID *string, pagination *graphqlmodel.Pagination) (*graphqlmodel.MessagesResponse, error) {
	response := &graphqlmodel.MessagesResponse{}

	newResources, errAuth := authHandler(ctx, *r.Resolver)
	if errAuth != nil {
		return response, errAuth
	}

	var limit = 50
	var page = 0

	if pagination != nil {
		// response based on the Pagination
		if pagination.Limit != 0 && pagination.Limit < 50 {
			limit = pagination.Limit
		}
		if pagination.Page != 0 && pagination.Page < 100 {
			page = pagination.Page
		}
	}

	// check if user is a room member
	isMember, err := graphqlmodel.IsRoomMember(ctx, newResources, newResources.User.Username, newResources.User.UserId, roomID)
	if err != nil {
		log.Println("err in GetRoom error is: ", err.Error())
		return response, err
	}
	if !isMember {
		msgErr := "permission denied"
		response.Error = &graphqlmodel.Error{
			Message: &msgErr,
		}
		return response, errors.New(msgErr)
	}

	messages, err := graphqlmodel.GetMessagesByRoomId(ctx, newResources, roomID, limit, page)
	if err != nil {
		log.Println("err in GetRoom error is: ", err.Error())
		return response, err
	}

	response.Data = convertToMessageGraphModel(messages, newResources.User.UserId)
	return response, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

func convertToMembershipGraphModel(mdlmsgs []models.RoomMembershipModel) []*graphqlmodel.RoomMembership {
	mmbrs := []*graphqlmodel.RoomMembership{}

	for i := 0; i < len(mdlmsgs); i++ {
		newMmbr := graphqlmodel.RoomMembership{
			ID:        &mdlmsgs[i].Id,
			RoomID:    &mdlmsgs[i].RoomId,
			JoinBy:    &mdlmsgs[i].JoinBy,
			Username:  &mdlmsgs[i].Username,
			UserID:    &mdlmsgs[i].UserId,
			IsAdmin:   &mdlmsgs[i].IsAdmin,
			MuteUntil: &mdlmsgs[i].MuteUntil,
			JoinAt:    &mdlmsgs[i].JoinAt,
			LeaveAt:   &mdlmsgs[i].LeaveAt,
		}
		mmbrs = append(mmbrs, &newMmbr)
	}

	return mmbrs
}

func convertToMessageGraphModel(mdlmsgs []models.MessageModel, userId string) []*graphqlmodel.Message {
	msgs := []*graphqlmodel.Message{}

	for i := 0; i < len(mdlmsgs); i++ {
		isOwner := mdlmsgs[i].UserId == userId
		newMsg := graphqlmodel.Message{
			ID:          &mdlmsgs[i].Id,
			GUID:        &mdlmsgs[i].GUID,
			Username:    &mdlmsgs[i].Username,
			UserID:      &mdlmsgs[i].UserId,
			RoomID:      &mdlmsgs[i].RoomId,
			MessageType: stringToMessageType(mdlmsgs[i].Type),
			TextContent: &mdlmsgs[i].TextContent,
			LinkURL:     &mdlmsgs[i].LinkUrl,
			CreatedAt:   &mdlmsgs[i].CreatedAt,
			UpdatedAt:   &mdlmsgs[i].UpdatedAt,
			DeletedAt:   &mdlmsgs[i].DeletedAt,
			DeletedFor:  &mdlmsgs[i].DeletedFor,
			IsOwner:     &isOwner,
		}
		msgs = append(msgs, &newMsg)
	}

	return msgs
}

func stringToMessageType(tp string) *graphqlmodel.MessageType {
	switch tp {
	case graphqlmodel.MessageTypeText.String():
		return &graphqlmodel.AllMessageType[0]
	case graphqlmodel.MessageTypeLink.String():
		return &graphqlmodel.AllMessageType[1]
	case graphqlmodel.MessageTypeLocation.String():
		return &graphqlmodel.AllMessageType[2]
	case graphqlmodel.MessageTypeImage.String():
		return &graphqlmodel.AllMessageType[3]
	case graphqlmodel.MessageTypeVoice.String():
		return &graphqlmodel.AllMessageType[4]
	case graphqlmodel.MessageTypeVideo.String():
		return &graphqlmodel.AllMessageType[5]
	case graphqlmodel.MessageTypePDF.String():
		return &graphqlmodel.AllMessageType[6]
	case graphqlmodel.MessageTypeCall.String():
		return &graphqlmodel.AllMessageType[7]
	case graphqlmodel.MessageTypeVideoCall.String():
		return &graphqlmodel.AllMessageType[8]
	case graphqlmodel.MessageTypeVoiceCall.String():
		return &graphqlmodel.AllMessageType[9]
	case graphqlmodel.MessageTypeOthers.String():
		return &graphqlmodel.AllMessageType[10]
	default:
		return &graphqlmodel.AllMessageType[10]
	}
}

func convertToRoomGraphModel(room models.RoomModel, isMember, isAdmin bool) *graphqlmodel.Room {
	return &graphqlmodel.Room{
		ID:         &room.Id,
		Title:      &room.Title,
		Image:      &room.Image,
		Biography:  &room.Biography,
		Creator:    &room.Creator,
		CreatorID:  &room.CreatorId,
		RoomType:   stringToRoomType(room.Type),
		IsReported: &room.IsReported,
		IsActive:   &room.IsActive,
		CreatedAt:  &room.CreatedAt,
		UpdatedAt:  &room.UpdatedAt,
		DeletedAt:  &room.DeletedAt,
		DeletedFor: &room.DeletedFor,
		IsMember:   &isMember,
		IsAdmin:    &isAdmin,
	}
}

func stringToRoomType(tp string) *graphqlmodel.RoomType {
	switch strings.ToTitle(tp) {
	case graphqlmodel.RoomTypeMutual.String():
		return &graphqlmodel.AllRoomType[0]
	case graphqlmodel.RoomTypeGroup.String():
		return &graphqlmodel.AllRoomType[1]
	case graphqlmodel.RoomTypeChannel.String():
		return &graphqlmodel.AllRoomType[2]
	case graphqlmodel.RoomTypeSupport.String():
		return &graphqlmodel.AllRoomType[3]
	default:
		return &graphqlmodel.AllRoomType[0]
	}
}

func membershipToRoom(ctx context.Context, res graphqlresources.Resources, roomMemberships []models.RoomMembershipModel) ([]*graphqlmodel.RoomSummery, error) {

	rooms := []*graphqlmodel.RoomSummery{}

	hasSupport := false
	for i := 0; i < len(roomMemberships); i++ {
		newRoomModel := models.RoomModel{
			Id: roomMemberships[i].RoomId,
		}
		newRoomModel, err := newRoomModel.GetItem(ctx, res.ROOMSDB, res.Timeout)
		if err != nil {
			log.Println("err in MembershipToRoom error is: ", err.Error())
			return nil, err
		}
		roomType := graphqlmodel.RoomType(newRoomModel.Type)
		if roomType == graphqlmodel.RoomTypeSupport {
			hasSupport = true
		}

		newRoom := graphqlmodel.RoomSummery{
			ID:       newRoomModel.Id,
			Title:    graphqlmodel.GenerateRoomTitle(ctx, res, newRoomModel),
			RoomType: roomType,
		}
		rooms = append(rooms, &newRoom)
	}

	sortedRooms := []*graphqlmodel.RoomSummery{}
	if hasSupport {
		for i := 0; i < len(rooms); i++ {
			if rooms[i].RoomType != graphqlmodel.RoomTypeSupport {
				sortedRooms = append(sortedRooms, rooms[i])
			}
		}
	} else {
		newRoom, createErr := graphqlmodel.CreateSupportRoom(ctx, res, res.User.Username, res.User.UserId)
		if createErr == nil {
			newRoomSummery := graphqlmodel.RoomSummery{
				ID:       newRoom.Id,
				Title:    graphqlmodel.GenerateRoomTitle(ctx, res, newRoom),
				RoomType: graphqlmodel.RoomType(newRoom.Type),
			}

			sortedRooms = append(sortedRooms, &newRoomSummery)
		}

		for i := 0; i < len(rooms); i++ {
			sortedRooms = append(sortedRooms, rooms[i])
		}

	}

	return sortedRooms, nil
}
