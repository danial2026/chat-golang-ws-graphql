package models

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"chat-golang-react/chat/common"

	"github.com/go-redis/redis/v8"
)

// usages:
// create connection by: username and room_id
// remove connection by: username and connection_id

type ConnectionModel struct {
	ConnectionID string `json:"connection_id"`
	Username     string `json:"username"`
	UserId       string `json:"user_id"`
	RoomId       string `json:"room_id"`
	JoinAt       int    `json:"join_at"`
}

func (c ConnectionModel) NewId() string {
	return common.GetULID()
}

func (c ConnectionModel) ItemToStruct(item []byte) (ConnectionModel, error) {
	if item == nil {
		return ConnectionModel{}, nil
	}
	var conn ConnectionModel
	if err := json.Unmarshal(item, &conn); err != nil {
		return ConnectionModel{}, err
	}
	return conn, nil
}

func (c ConnectionModel) ToNewStruct() (ConnectionModel, error) {
	// set inputs
	conn := ConnectionModel{
		ConnectionID: c.ConnectionID,
		Username:     c.Username,
		UserId:       c.UserId,
		RoomId:       c.RoomId,
		JoinAt:       common.Now(),
	}

	return conn, nil
}

func (c ConnectionModel) ToNewConversationStruct() (ConnectionModel, error) {
	// set inputs
	conn := ConnectionModel{
		ConnectionID: c.ConnectionID,
		Username:     c.Username,
		UserId:       c.UserId,
		RoomId:       c.RoomId,
		JoinAt:       common.Now(),
	}

	return conn, nil
}

func (c ConnectionModel) ToNewItem() (ConnectionModel, []byte, error) {
	conn, creationErr := c.ToNewStruct()
	if creationErr != nil {
		return conn, nil, creationErr
	}
	data, marshalErr := json.Marshal(conn)
	if marshalErr != nil {
		return conn, nil, errors.New("err in ConnectionModel/ToNewConversationItem/Marshal: " + marshalErr.Error())
	}
	return conn, data, nil
}

func (c ConnectionModel) ToNewConversationItem() (ConnectionModel, []byte, error) {
	conn, creationErr := c.ToNewConversationStruct()
	if creationErr != nil {
		return conn, nil, creationErr
	}
	data, marshalErr := json.Marshal(conn)
	if marshalErr != nil {
		return conn, nil, errors.New("err in ConnectionModel/ToNewConversationItem/Marshal: " + marshalErr.Error())
	}
	return conn, data, nil
}

func (c ConnectionModel) Create(ctx context.Context, db *redis.Client, timeout time.Duration) (ConnectionModel, error) {
	conn, serializedData, itemErr := c.ToNewItem()
	if itemErr != nil {
		return conn, itemErr
	}

	log.Println("ConnectionModel/Create: ", c.RoomId, " : ", c.ConnectionID)
	log.Println("ConnectionModel/Create/timeout: ", timeout)
	TCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	log.Println("ConnectionModel/Create/WithTimeout: ", TCtx)

	// creationErr := db.Set(TCtx, string(c.ConnectionID), string(serializedData), timeout)
	creationErr := db.HSet(TCtx, string(c.RoomId), string(c.ConnectionID), string(serializedData))

	if creationErr.Err() != nil {
		return conn, errors.New("err in ConnectionModel/Create: " + creationErr.Err().Error())
	}

	return conn, nil
}

func (c ConnectionModel) Delete(ctx context.Context, db *redis.Client, timeout time.Duration) (bool, error) {
	// delete item and if deletion is affected any row return true
	log.Println("ConnectionModel/Delete/timeout: ", timeout)
	TCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	delItem, deleteErr := db.HDel(TCtx, c.RoomId, c.ConnectionID).Result()
	if deleteErr != nil {
		return false, errors.New("err in ConnectionModel/Delete/DeleteItemWithContext: " + deleteErr.Error())
	}
	log.Println("ConnectionModel/Delete/delItem: ", delItem)

	return true, nil
}

func (c ConnectionModel) ListByRoom(ctx context.Context, db *redis.Client, timeout time.Duration) ([]ConnectionModel, error) {
	if c.RoomId == "" {
		return []ConnectionModel{}, errors.New("err in ConnectionModel/ListByRoom: can not be empty")
	}

	TCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	keys := db.HGetAll(TCtx, c.RoomId).Val()

	var conns []ConnectionModel
	for _, key := range keys {
		var conn ConnectionModel
		if marshalListErr := json.Unmarshal([]byte(key), &conn); marshalListErr != nil {
			log.Println("err in ConnectionModel/ListByRoom: " + marshalListErr.Error())
			return []ConnectionModel{}, errors.New("err in ConnectionModel/ListByRoom/Unmarshal: " + marshalListErr.Error())
		}
		conns = append(conns, conn)
	}

	return conns, nil
}

func (c ConnectionModel) CreateConversation(ctx context.Context, db *redis.Client, timeout time.Duration) (ConnectionModel, error) {
	conn, data, itemErr := c.ToNewConversationItem()
	if itemErr != nil {
		return conn, itemErr
	}

	log.Println("ConnectionModel/CreateConversation/ToNewConversationItem: ConnectionModel: ", conn, " Attribute: ", data)

	log.Println("ConnectionModel/CreateConversation/timeout: ", timeout)
	TCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	log.Println("ConnectionModel/CreateConversation/WithTimeout: ", TCtx)

	creationErr := db.HSet(TCtx, string(c.RoomId), string(c.ConnectionID), string(data))

	if creationErr != nil {
		return conn, errors.New("err in ConnectionModel/CreateConversation: " + creationErr.Err().Error())
	}

	return conn, nil
}

// func (c ConnectionModel) GetConversationItem(ctx context.Context, db *redis.Client, timeout time.Duration) (ConnectionModel, error) {
// 	keys, keysErr := c.GetConversationKeys()
// 	log.Println("ConnectionModel/GetConversationItem/GetConversationKeys: ", keys)
// 	if keysErr != nil {
// 		log.Println("ConnectionModel/GetConversationItem/keysErr: ", keysErr)
// 		return ConnectionModel{}, keysErr
// 	}
// 	input := &dynamodb.GetItemInput{
// 		Key: keys,
// 		// table name
// 		TableName: c.GetConversationTableName(),
// 	}

// 	TCtx, cancel := context.WithTimeout(ctx, timeout)
// 	defer cancel()

// 	result, readErr := db.Query(TCtx, input)
// 	if readErr != nil {
// 		log.Println("ConnectionModel/GetConversationItem/Query: ", readErr)
// 		return ConnectionModel{}, errors.New("Err in ConnectionModel/GetConversationItem/Query: " + readErr.Error())
// 	}

// 	conn, unmarshalErr := c.ItemToStruct(result.Item)
// 	if unmarshalErr != nil {
// 		log.Println("ConnectionModel/GetConversationItem/ItemToStruct: ", unmarshalErr)
// 		return conn, unmarshalErr
// 	}

// 	return conn, nil
// }

func (c ConnectionModel) DeleteConversation(ctx context.Context, db *redis.Client, connectionId string, timeout time.Duration) error {
	log.Println("ConnectionModel/DeleteConversation: ")
	TCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, deleteErr := db.Del(TCtx, string(connectionId)).Result()
	if deleteErr != nil {
		return errors.New("err in ConnectionModel/DeleteConversation/DeleteItemWithContext: " + deleteErr.Error())
	}

	return nil
}
