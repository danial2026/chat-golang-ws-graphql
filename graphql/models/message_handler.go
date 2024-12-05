package model

import (
	"context"

	"chat-golang-react/chat/graphql/resources"
	models "chat-golang-react/chat/models"

	"log"
)

func GetMessagesByRoomId(ctx context.Context, res resources.Resources, roomId string, limit, page int) ([]models.MessageModel, error) {
	msg := models.MessageModel{
		RoomId: roomId,
	}

	result, err := msg.GetAllItems(ctx, res.ROOMSDB, res.Timeout, limit, page)
	if err != nil {
		log.Println("err in GetMessagesByRoomId error is: ", err.Error())
		return []models.MessageModel{}, err
	}

	return result, nil
}
