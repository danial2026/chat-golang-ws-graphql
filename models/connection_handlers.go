package models

import (
	"context"
	"log"

	"chat-golang-react/chat/websocket/resources"
)

// ----------------------------- ConnectionsTable's related functions

func CreateConnection(ctx context.Context, res resources.Resource, connId, username, userId, roomId string) error {
	conn := ConnectionModel{
		ConnectionID: connId,
		Username:     username,
		UserId:       userId,
		RoomId:       roomId,
	}

	_, err := conn.Create(ctx, res.CONNDB, res.Timeout)
	if err != nil {
		log.Println("err in CreateConnection/Create " + err.Error())
		return err
	}
	return nil
}

func GetConnectionsList(ctx context.Context, res resources.Resource, roomId string) ([]ConnectionModel, error) {
	conn := ConnectionModel{
		RoomId: roomId,
	}
	return conn.ListByRoom(ctx, res.CONNDB, res.Timeout)
}

func SplitToEqualChunks(conns []ConnectionModel, size int) [][]ConnectionModel {
	var results [][]ConnectionModel

	var j int
	for i := 0; i < len(conns); i += size {
		j += size
		if j > len(conns) {
			j = len(conns)
		}
		results = append(results, conns[i:j])
	}

	return results
}

func removeConnectionsList(ctx context.Context, res resources.Resource, connList []ConnectionModel, table string) {
	go func() {
		for _, conn := range connList {

			err := DisconnectConnection(ctx, res, conn.RoomId)
			if err != nil {
				log.Println("err in Connection_Handler/removeConnectionsList/DisconnectConnection:  " + err.Error())
			}
		}
	}()
}

func RemoveConnectionsList(ctx context.Context, res resources.Resource, connList []ConnectionModel, table string) error {
	// TODO : for connections more than 25, make concurrent request
	var results []error

	if len(connList) > 0 {
		connsChunks := SplitToEqualChunks(connList, 25)
		for _, item := range connsChunks {
			removeConnectionsList(ctx, res, item, table)
		}
	}

	if len(results) > 0 {
		return results[0]
	}
	return nil

}

// -----------------------------  ConversationListTable's functions:

func DisconnectConnection(ctx context.Context, res resources.Resource, roomId string) error {
	conn := ConnectionModel{
		ConnectionID: res.User.ConnectionId,
		Username:     res.User.Username,
		UserId:       res.User.UserId,
		RoomId:       roomId,
	}

	isDeleted, delErr := conn.Delete(ctx, res.CONNDB, res.Timeout)
	if delErr != nil {
		return delErr
	}

	if !isDeleted {
		delErr2 := conn.DeleteConversation(ctx, res.CONNDB, res.User.ConnectionId, res.Timeout)
		if delErr2 != nil {
			return delErr2
		}
		return nil
	}
	return nil
}

func getConversationConnectionList(ctx context.Context, res resources.Resource, connList []ConnectionModel) ([]ConnectionModel, error) {
	// inputsRequest := map[string]*dynamodb.KeysAndAttributes{}

	// TCtx, cancel := context.WithTimeout(ctx, res.Timeout)
	// defer cancel()

	// for _, conn := range connList {
	// 	pk, _ := conn.GetConversationPK()
	// 	sk, _ := conn.GetConversationSK()
	// 	inputsRequest[ConversationListTable] = &dynamodb.KeysAndAttributes{
	// 		Keys: []map[string]*dynamodb.AttributeValue{
	// 			{

	// 				"PK": {S: aws.String(pk)},
	// 				"SK": {S: aws.String(sk)},
	// 			},
	// 		},
	// 	}
	// }

	// input := &dynamodb.BatchGetItemInput{
	// 	RequestItems: inputsRequest,
	// }

	// results, err := res.RDB.BatchQuery(TCtx, input)
	// if err != nil {
	// 	log.Println("err in ConnectionModel/GetConversationConnectionList/BatchQuery:  " + err.Error())
	// 	return []ConnectionModel{}, errors.New("err in ConnectionModel/GetConversationConnectionList/BatchQuery: " + err.Error())
	// }

	// items, ok := results.Responses[ConversationListTable]
	// if !ok {
	// 	return []ConnectionModel{}, nil
	// }
	// var conns []ConnectionModel
	// if marshalErr := dynamodbattribute.UnmarshalListOfMaps(items, &conns); marshalErr != nil {
	// 	return []ConnectionModel{}, marshalErr
	// }

	// return conns, nil
	return []ConnectionModel{}, nil
}

func GetConversationConnectionList(ctx context.Context, res resources.Resource, connList []ConnectionModel) ([]ConnectionModel, error) {
	var conns []ConnectionModel

	if len(connList) > 0 {
		connsChunks := SplitToEqualChunks(connList, 100)
		for _, item := range connsChunks {
			results, err := getConversationConnectionList(ctx, res, item)
			if err != nil {
				return conns, err
			}
			conns = append(results, conns...)
		}
	}

	return conns, nil
}
