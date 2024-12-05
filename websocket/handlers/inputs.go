package handlers

// connection related's:
const (
	ConnectToRoom              = "connect-to-room"
	ConnectToConversationsList = "connect-to-room-conversations-list"
)

const (
	ConversationContacts = "conversation-contacts-list"
)

// message related's:
const (
	CreateMessageOperation       = "create"
	EditMessageOperation         = "edit"
	DeleteMessageOperation       = "delete"
	DeleteMessageForMeOperation  = "deleteForMe"
	DeleteMessageForAllOperation = "deleteForAll"
)

type InputMessageBody struct {
	//Username    string `json:"username"`
	Action      string `json:"action"`
	Operation   string `json:"operation"`
	MessageId   string `json:"message_id"`
	RoomId      string `json:"room_id"`
	Type        string `json:"type"`
	TextContent string `json:"text_content"`
	LinkUrl     string `json:"link_url"`
}

type QueryStringParameters struct {
	Operation string
	RoomId    string
	Type      string
}
