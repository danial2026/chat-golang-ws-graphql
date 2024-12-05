package models

import (
	"context"
	"log"

	"chat-golang-react/chat/websocket/resources"
)

func GetMessage(ctx context.Context, res resources.Resource, messageId, roomId, username, userId string) (MessageModel, error) {
	msg := MessageModel{
		Id:       messageId,
		RoomId:   roomId,
		Username: username,
		UserId:   userId,
	}
	newMsg, err := msg.GetItem(ctx, res.MESSAGEDB, res.Timeout)
	if err != nil {
		return msg, err
	}
	return newMsg, nil
}

func CreateMessage(ctx context.Context, res resources.Resource, username, userId string, msg MessageModel) (MessageModel, error) {
	newMsg := MessageModel{
		Username:    username,
		UserId:      userId,
		RoomId:      msg.RoomId,
		Type:        msg.Type,
		TextContent: msg.TextContent,
		LinkUrl:     msg.LinkUrl,
	}

	newMsg, err := msg.Create(ctx, res.MESSAGEDB, res.Timeout)
	if err != nil {
		log.Println("err in CreateTextMessage/Create: " + err.Error())
		return MessageModel{}, err
	}

	return newMsg, nil
}

func IsMessageOwner(ctx context.Context, res resources.Resource, messageId, roomId, username, userId string) (MessageModel, bool, error) {
	msg := MessageModel{
		Id:       messageId,
		RoomId:   roomId,
		Username: username,
		UserId:   userId,
	}

	outMsg, isExists, err := msg.Exists(ctx, res.MESSAGEDB, res.Timeout)
	if err != nil {
		return msg, false, err
	}
	if !isExists {
		return outMsg, false, nil
	}

	if outMsg.Username != msg.Username {
		return outMsg, false, nil
	}
	return outMsg, true, nil
}

func DeleteMessageForMe(ctx context.Context, res resources.Resource, messageId, roomId, username, userId string) (MessageModel, error) {
	msg := MessageModel{
		Id:       messageId,
		RoomId:   roomId,
		Username: username,
		UserId:   userId,
	}
	err := msg.DeleteItemFor(ctx, res.MESSAGEDB, res.Timeout, DeletedForMe)
	if err != nil {
		return msg, err
	}
	return msg, nil
}

func DeleteMessageForAll(ctx context.Context, res resources.Resource, messageId, roomId, username, userId string) (MessageModel, error) {
	msg := MessageModel{
		Id:       messageId,
		RoomId:   roomId,
		Username: username,
		UserId:   userId,
	}
	err := msg.DeleteItemFor(ctx, res.MESSAGEDB, res.Timeout, DeletedForAll)
	if err != nil {
		return msg, err
	}
	return msg, nil
}
