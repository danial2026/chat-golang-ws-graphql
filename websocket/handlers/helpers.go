package handlers

import (
	"context"
	"log"

	"chat-golang-react/chat/models"
	"chat-golang-react/chat/websocket/resources"
)

func GetNotSentUsers(ctx context.Context, res resources.Resource, sentConns []models.ConnectionModel, allMembers []models.RoomMembershipModel) ([]models.ConnectionModel, error) {
	var notSentUsers []models.ConnectionModel
	for _, member := range allMembers {
		found := false
		for _, conn := range sentConns {
			if member.Username == conn.Username {
				found = true
			}
		}

		if !found {
			notSentUsers = append(notSentUsers, models.ConnectionModel{
				Username: member.Username,
				RoomId:   member.RoomId,
			})
		}

	}
	// get conversation list room connections:
	connList, err := models.GetConversationConnectionList(ctx, res, notSentUsers)
	if err != nil {
		log.Println("err in GetNotSentUsers/GetConversationConnectionList: ", err.Error())
		return connList, err
	}

	return connList, nil
}
