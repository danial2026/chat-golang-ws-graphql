package pushers

import (
	"context"
	"encoding/json"
	"errors"
	"sync"

	"chat-golang-react/chat/models"
	"chat-golang-react/chat/websocket/resources"
)

const (
	SendToRoomPush           = "push-to-room"
	SendToConversionListPush = "push-to-conversation-list"
)

const (
	NewMessage    = "new_message"
	DeleteMessage = "delete_message"
)

type PushInput struct {
	Action         string
	RoomId         string
	ConnectionList []models.ConnectionModel
	Message        models.MessageModel
}

func getConnectionList(ctx context.Context, res resources.Resource, in PushInput) ([]models.ConnectionModel, error) {
	switch in.Action {
	case SendToRoomPush:
		conns, getErr := models.GetConnectionsList(ctx, res, in.RoomId)
		if getErr != nil {
			return []models.ConnectionModel{}, errors.New("err in PushNewMessage/GetConnectionsList: " + getErr.Error())
		}
		return conns, nil
	case SendToConversionListPush:
		return in.ConnectionList, nil
	default:
		return []models.ConnectionModel{}, errors.New("err in PushNewMessage/getConnectionList invalid action")
	}
}

func PushNewMessage(ctx context.Context, res resources.Resource, in PushInput) ([]models.ConnectionModel, error) {
	/*
		Logic:
		1- jsonify the message
		2 - get connections list by roomId
		3- send message by the connectionId
		4- if connection is gone in sending process delete the connection
		5- push message for not sent persons in the room in the ConversationsListRoom
	*/

	// jsonify the message
	body, jsonErr := json.Marshal(&in.Message)
	if jsonErr != nil {
		return []models.ConnectionModel{}, errors.New("err in PushNewMessage/json.Marshal: " + jsonErr.Error())
	}
	newMsg := map[string]interface{}{
		"type": NewMessage,
		"body": string(body),
	}

	// get connections list:
	conns, getErr := getConnectionList(ctx, res, in)
	if getErr != nil {
		return []models.ConnectionModel{}, getErr
	}

	connsLen := len(conns)
	if connsLen == 0 {
		return []models.ConnectionModel{}, nil
	}

	// send message:
	var wg sync.WaitGroup

	for _, conn := range conns {
		wg.Add(1)
		isOwner := in.Message.UserId == conn.UserId
		AsyncPostToConnectionBasic(conn, newMsg, isOwner, func() { wg.Done() })
	}

	go func() {
		wg.Wait()
	}()

	return conns, nil
}
