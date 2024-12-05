# chat-service GraphQL

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

> This runs the graphQL application at port `4000`. Try accessing API at `http://localhost:4000/query`

* run this to regenerate files

```bash
go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen generate
```

### A list of requests bodies communicated as post request:

* create mutual room

```graphql
mutation {
	createRoom (input: {
		title: "test",
		type: MUTUAL
	}) {
		id
		title
	}
}
```

## TODO : 

[x] create room
```
	CreateRoom
```
[] create membership / join the room or add memeber to room
```
	AddMemberToRoom(roomId, memberId)
		- if not member then create new membership
```
[?] delete membership / leave the room
```
	LeaveRoom(roomId)
		- if member then kick (mark membership as deleted)
```
[?] delete membership by admin/creator of the room
```
	KickMember(roomId, memberId)
		- if member then kick (mark membership as deleted and set deleted by to admin's Id)
```
