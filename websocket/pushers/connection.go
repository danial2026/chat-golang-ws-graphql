package pushers

import (
	"log"

	"chat-golang-react/chat/models"
	"chat-golang-react/chat/websocket/resources"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/syncmap"
)

type Session struct {
	Conn          *websocket.Conn
	UserSessionID string
}

type Connection struct {
	UserSessionIds []string
}

type Result struct {
	Connection models.ConnectionModel
	Error      error
}

var Connections = syncmap.Map{}
var Sessions = syncmap.Map{}

func ConnectionAdd(conn *websocket.Conn, user resources.User) {
	var newUserSessionIds []string
	val, ok := Connections.Load(user.UserId)
	if ok {
		otherConnections, ok := val.(Connection)
		if !ok {
			// this will break iteration
			log.Fatalln("ConnectionAdd: error in map load")
		}
		newUserSessionIds = append(otherConnections.UserSessionIds, user.ConnectionId)
	} else {
		newUserSessionIds = []string{user.ConnectionId}
	}

	Connections.Store(user.UserId, Connection{
		UserSessionIds: newUserSessionIds,
	})

	Sessions.Store(user.ConnectionId, Session{
		Conn:          conn,
		UserSessionID: user.ConnectionId,
	})
}

func ConnectionDel(userId, connectionId string) {
	Sessions.Delete(connectionId)
	val, ok := Connections.Load(userId)
	if !ok {
		// this will break iteration
		log.Fatalln("ConnectionDel: error in map load")
	}

	otherConnections, ok := val.(Connection)
	if !ok {
		// this will break iteration
		log.Fatalln("ConnectionDel: error in map load")
	}

	if len(otherConnections.UserSessionIds) == 0 {
		Connections.Delete(connectionId)
	}
}

func AsyncPostToConnection(conn models.ConnectionModel, data map[string]interface{}, c chan<- Result, onExit func()) {
	// send messages to each active session
	// go func() {
	defer onExit()

	val, ok := Connections.Load(conn.UserId)
	if ok {
		otherConnections, ok := val.(Connection)
		if !ok {
			// this will break iteration
			log.Fatalln("AsyncPostToConnection: error in map load")
		}

		for _, item := range otherConnections.UserSessionIds {
			val, ok := Sessions.Load(item)
			if ok {

				otherSessions, ok := val.(Session)
				if !ok {
					// this will break iteration
					log.Fatalln("AsyncPostToConnection/for: error in map load")
				}

				err := otherSessions.Conn.WriteJSON(data)
				c <- Result{Connection: conn, Error: err}
			}
		}
	}
	// }()
}

func AsyncPostToConnectionBasic(conn models.ConnectionModel, data map[string]interface{}, isOwner bool, onExit func()) {
	// send messages to each active session
	defer onExit()

	// go func() {
	val, ok := Connections.Load(conn.UserId)
	data["is_owner"] = isOwner
	if ok {
		otherConnections, ok := val.(Connection)
		if !ok {
			// this will break iteration
			log.Fatalln("AsyncPostToConnectionBasic: error in map load")
		}

		for _, item := range otherConnections.UserSessionIds {
			val, ok := Sessions.Load(item)
			if ok {
				otherSessions, ok := val.(Session)
				if !ok {
					// this will break iteration
					log.Fatalln("AsyncPostToConnectionBasic/for: error in map load")
				}

				err := otherSessions.Conn.WriteJSON(data)
				if err != nil {
					log.Println("AsyncPostToConnectionBasic/for/WriteJSON:", err)
				}
			}
		}
	}
	// }()
}
