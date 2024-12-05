package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"chat-golang-react/chat/common/decoders"
	"chat-golang-react/chat/websocket/configs"
	"chat-golang-react/chat/websocket/pushers"
	"chat-golang-react/chat/websocket/resources"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
)

type WebSocketMiddleware struct {
	// websocket upgrader
	Upgrader websocket.Upgrader

	// resources
	Res resources.Resource
}

func (ws WebSocketMiddleware) AuthorizerMiddleware(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing authorizerMiddleware")
	var connectionsIsSuccessfully = false

	ctx := context.Background()
	token := r.URL.Query().Get("token")

	// Upgrade initial GET request to a websocket connection
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during connection upgradation:", err)
		return
	}

	if token == "" {
		log.Println("Token is empty")
		msgErr := "unauthorized"
		msg := ForbiddenErrorResponse(msgErr)
		conn.WriteJSON(msg)
		conn.Close()
		return
	}

	tokenClaims, err := decoders.ParseToken(ctx, token)
	if err != nil {
		log.Println("Error during checking token:", err)
		conn.Close()
		return
	}
	tokenClaims = decoders.GenerateContextAuth(tokenClaims, token)

	queryParameters := QueryStringParameters{
		Operation: r.URL.Query().Get("operation"),
		RoomId:    r.URL.Query().Get("room_id"),
	}

	newRes, err := resources.ConstructUserResource(tokenClaims, ws.Res, queryParameters.RoomId)
	if err != nil {
		log.Println("Error during resource construction:", err)
		conn.Close()
		return
	}

	defer func() {
		log.Println("Executing defer")
		// var resmsg MessageResponse
		if connectionsIsSuccessfully {
			DisconnectHandler(ctx, newRes, queryParameters.RoomId)

			// remove connection from the list of connections
			pushers.ConnectionDel(newRes.User.UserId, newRes.User.ConnectionId)
		}
		conn.Close()
	}()

	resmsg := ConnectHandler(ctx, queryParameters, &newRes)
	conn.WriteJSON(resmsg)
	// Room not found
	if resmsg.StatusCode == 404 {
		conn.Close()
		return
	}

	// add connection to a list of connections contain the connetion and connectionId and username
	pushers.ConnectionAdd(conn, newRes.User)
	connectionsIsSuccessfully = true

	switch err {
	case nil:
		ws.actionHandler(ctx, conn, tokenClaims, newRes)
		conn.Close()
		return
	case decoders.ErrTokenExpire:
		conn.WriteMessage(websocket.TextMessage, []byte("token expired"))
		conn.Close()
		return
	default:
		conn.WriteMessage(websocket.TextMessage, []byte("unrecognized WebSocket action"))
		conn.Close()
		return
	}
}

func (ws WebSocketMiddleware) actionHandler(ctx context.Context, conn *websocket.Conn, tokenClaims jwt.MapClaims, newRes resources.Resource) {
	log.Println("Executing websocketFunc")
	// The event loop
	for {
		var inputMessage InputMessageBody
		messageType, inputMessageByte, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			if inputMessageByte == nil {
				break
			}
			continue
		}
		if messageType != 1 {
			log.Println("Error messageType: messageType wasnt string")
			continue
		}
		if err := json.Unmarshal(inputMessageByte, &inputMessage); err != nil {
			log.Println("Error during message reading:", err)
			log.Println("message :", inputMessageByte)
			continue
		}

		log.Printf("Received: %s", inputMessage)

		switch inputMessage.Action {
		case configs.SendMessageRoute:
			resmsg := SendMessageHandler(ctx, inputMessage, newRes)
			conn.WriteJSON(resmsg)
			if err != nil {
				log.Println("Error during message writing:", err)
				return
			}
		default:
			resmsg := DefaultHandler()
			conn.WriteJSON(resmsg)
			if err != nil {
				log.Println("Error during message writing:", err)
				continue
			}
		}
	}
}
