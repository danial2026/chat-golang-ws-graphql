package handlers

import (
	"errors"
	"log"

	"chat-golang-react/chat/common"
	"chat-golang-react/chat/models"
)

type ConnectionOutput struct {
	Connection models.ConnectionModel
	Operation  string
}

type MessageOutput struct {
	Message   models.MessageModel
	Operation string
}

func ValidateConnect(request QueryStringParameters, userId, username string) (ConnectionOutput, error) {
	roomId := request.RoomId
	operation := request.Operation

	if operation == "" {
		operation = ConnectToRoom
	} else if operation == ConnectToConversationsList {
		roomId = models.ConversationsListRoom
		//roomOk = true
	}

	if !common.ContainsInStringArray([]string{ConnectToRoom, ConnectToConversationsList}, operation) {
		return ConnectionOutput{}, errors.New("operation is invalid")
	}

	//if !roomOk {
	//	targetUser, UserOk := request.QueryStringParameters["username"]
	//	if !UserOk {
	//		return ConnectionOutput{}, errors.New("username or room_id required")
	//	}
	//	if targetUser == "" {
	//		return ConnectionOutput{}, errors.New("username is empty")
	//	}
	//
	//	return ConnectionOutput{
	//		Room: RequestedRoom{
	//			IsMutual:      true,
	//			RequestedUser: username,
	//			TargetUser:    targetUser,
	//		},
	//		Connection: models.ConnectionModel{
	//			RoomId:       "",
	//			Username:     username,
	//			ConnectionID: connId,
	//		},
	//		Operation: operation,
	//	}, nil
	//}

	if roomId == "" {
		return ConnectionOutput{}, errors.New("room_id is empty")
	}

	connModel := models.ConnectionModel{
		Username: username,
		UserId:   userId,
		RoomId:   roomId,
	}

	// TODO : make sure this works
	connModel.ConnectionID = connModel.NewId()

	return ConnectionOutput{
		Connection: connModel,
		Operation:  operation,
	}, nil
}

func validateMessageInput(request InputMessageBody) (InputMessageBody, error) {

	// validate action
	validActions := []string{CreateMessageOperation, EditMessageOperation, DeleteMessageOperation}
	if !common.ContainsInStringArray(validActions, request.Operation) {
		return request, errors.New("action is invalid")
	}

	// validate room id:
	if request.RoomId == "" {
		return request, errors.New("room_id is required")
	}

	// TODO : validate room id :
	//reqRoom := RequestedRoom{}
	//if msg.RoomId == "" {
	//	if msg.Username == "" {
	//		return msg, reqRoom, errors.New("room_id or username is required")
	//	}
	//	reqRoom.RequestedUser = username
	//	reqRoom.TargetUser = msg.Username
	//	reqRoom.IsMutual = true
	//} else {
	//	reqRoom.RoomId = msg.RoomId
	//
	//	if isMutual(msg.RoomId) {
	//		reqRoom.IsMutual = true
	//	}
	//
	//}

	return request, nil
}

func validateTextMessage(msg InputMessageBody, username, userId, roomId string) (MessageOutput, error) {
	// validate message type:
	if msg.Type != models.MessageTextType {
		return MessageOutput{}, errors.New("type is invalid")
	}

	// validate and remove null chars from the text_content
	validMsg := common.RemoveNullChars(msg.TextContent)

	if common.IsSpaceString(validMsg) {
		return MessageOutput{}, errors.New("text_content is empty")
	}

	// construct message
	newMsg := models.MessageModel{
		RoomId:      roomId,
		TextContent: validMsg,
		Username:    username,
		UserId:      userId,
		Type:        msg.Type,
	}

	return MessageOutput{
		Message:   newMsg,
		Operation: msg.Operation,
	}, nil
}

func validateVideoCallMessage(msg InputMessageBody, username, userId, roomId string) (MessageOutput, error) {
	// validate message type:
	if msg.Type != models.MessageVideoCallType {
		return MessageOutput{}, errors.New("type is invalid")
	}

	// validate and remove null chars from the text_content
	validMsg := common.RemoveNullChars(msg.TextContent)

	if common.IsSpaceString(validMsg) {
		return MessageOutput{}, errors.New("text_content is empty")
	}

	if msg.LinkUrl == "" {
		return MessageOutput{}, errors.New("link_url is empty")
	}

	// construct message
	newMsg := models.MessageModel{
		RoomId:      roomId,
		TextContent: validMsg,
		Username:    username,
		UserId:      userId,
		LinkUrl:     msg.LinkUrl,
		Type:        msg.Type,
	}

	return MessageOutput{
		Message:   newMsg,
		Operation: msg.Operation,
	}, nil
}

func validateVoiceCallMessage(msg InputMessageBody, username, userId, roomId string) (MessageOutput, error) {
	// validate message type:
	if msg.Type != models.MessageVoiceCallType {
		return MessageOutput{}, errors.New("type is invalid")
	}

	// validate and remove null chars from the text_content
	validMsg := common.RemoveNullChars(msg.TextContent)

	if common.IsSpaceString(validMsg) {
		return MessageOutput{}, errors.New("text_content is empty")
	}

	if msg.LinkUrl == "" {
		return MessageOutput{}, errors.New("link_url is empty")
	}

	// construct message
	newMsg := models.MessageModel{
		RoomId:      roomId,
		TextContent: validMsg,
		Username:    username,
		UserId:      userId,
		LinkUrl:     msg.LinkUrl,
		Type:        msg.Type,
	}

	return MessageOutput{
		Message:   newMsg,
		Operation: msg.Operation,
	}, nil
}

func deleteMessage(msg InputMessageBody, username, userId, roomId string) (MessageOutput, error) {
	// validate message id
	if msg.MessageId == "" {
		return MessageOutput{}, errors.New("message_id is empty")
	}
	_, IsValidId := common.GetAndValidateULID(msg.MessageId)
	if IsValidId != nil {
		return MessageOutput{}, errors.New("message_id is invalid")
	}

	// construct the message for delete
	targetMsg := models.MessageModel{
		Id:       msg.MessageId,
		RoomId:   msg.RoomId,
		Username: username,
		UserId:   userId,
	}

	return MessageOutput{
		Message:   targetMsg,
		Operation: msg.Operation,
	}, nil

}

func ValidateSendMessage(request InputMessageBody, username, userId, roomId string) (MessageOutput, error) {
	// general validation:
	msg, inputErr := validateMessageInput(request)
	if inputErr != nil {
		log.Println(inputErr)
		return MessageOutput{}, inputErr
	}

	switch msg.Operation {
	case CreateMessageOperation:
		return ValidateMessageType(request, username, userId, roomId)
	case DeleteMessageOperation:
		return deleteMessage(msg, username, userId, roomId)
	default:
		return MessageOutput{}, errors.New("operation is invalid")
	}
}

func ValidateMessageType(request InputMessageBody, username, userId, roomId string) (MessageOutput, error) {
	// general validation:
	msg, inputErr := validateMessageInput(request)
	if inputErr != nil {
		log.Println(inputErr)
		return MessageOutput{}, inputErr
	}

	switch msg.Type {
	case models.MessageTextType:
		return validateTextMessage(msg, username, userId, roomId)
	case models.MessageVideoCallType:
		return validateVideoCallMessage(msg, username, userId, roomId)
	case models.MessageVoiceCallType:
		return validateVoiceCallMessage(msg, username, userId, roomId)
	default:
		return MessageOutput{}, errors.New("type is invalid")
	}
}
