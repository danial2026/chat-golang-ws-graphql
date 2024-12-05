package main

import (
	"log"
	"net/http"
	"os"

	"chat-golang-react/chat/websocket/handlers"
	"chat-golang-react/chat/websocket/resources"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalln("Error loading .env file : ", err)
	}

	res, err := resources.ConstructResource()
	if err != nil {
		log.Println("Error during resource construction:", err)

		return
	}

	defer func() {
		// TODO : flush the redis database
		log.Println(" >>> flushing redis")
		res.CONNDB.FlushAll(res.CONNDB.Context())
		res.CONNDB.FlushDB(res.CONNDB.Context())
		res.CONNDB.Close()
	}()

	// TODO : use more specific options
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			// alow all clients
			return true
		},
	}

	log.Println("Starting server...")
	ws := handlers.WebSocketMiddleware{
		Res:      res,
		Upgrader: upgrader,
	}

	log.Println("websocket port:", os.Getenv("WSPORT"))
	http.HandleFunc("/ws", ws.AuthorizerMiddleware)
	log.Fatal(http.ListenAndServe(":"+os.Getenv("WSPORT"), nil))
}
