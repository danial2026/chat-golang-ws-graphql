# chat-service Websockets

* build golang

```bash
go mod tidy
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go
```

* build docker 

```bash
docker-compose -f docker-compose-develop.yml up -d  --build
docker-compose -f docker-compose-production.yml up -d  --build
```

* run 

```bash
go mod tidy
go run main.go
```

> This runs the websocket application at port `8066`. Try accessing API at `http://localhost:8066/ws`

### A list of response and request communicated through websocket:

* client seccessfully connected to the room

```
{
  "body": "successfully connected to the room",
  "code": 200,
}
```

* client seccessfully disconnected from the room

```
{
  "body": "connectionId is successfully disconnected",
  "code": 200,
}
```

* create/delete message response

```
{
  "body": "message successfully handeled",
  "code": 200,
}
```

* create message request

```
{
  "action": "sendMessage",
  "room_id": "roomID",
  "type": "text",
  "text_content": "msg as text",
  "operation": "create",
  "message_id": "",
}
```

* delete message request

    operations: `delete`, `deleteForMe`, `deleteForAll`

```
{
  "action": "sendMessage",
  "room_id": "roomID",
  "type": "text",
  "text_content": "msg as text",
  "operation": "delete",
  "message_id": "",
}
```

### Tests:

* client nodejs

```bash
node ../clients/client.js
```
