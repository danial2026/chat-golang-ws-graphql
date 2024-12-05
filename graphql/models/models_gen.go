// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Error struct {
	Message *string `json:"message,omitempty"`
	Code    *int    `json:"code,omitempty"`
}

type Message struct {
	ID          *string      `json:"id,omitempty"`
	GUID        *string      `json:"guid,omitempty"`
	Username    *string      `json:"username,omitempty"`
	UserID      *string      `json:"user_id,omitempty"`
	RoomID      *string      `json:"room_id,omitempty"`
	MessageType *MessageType `json:"message_type,omitempty"`
	TextContent *string      `json:"text_content,omitempty"`
	LinkURL     *string      `json:"link_url,omitempty"`
	CreatedAt   *int         `json:"created_at,omitempty"`
	UpdatedAt   *int         `json:"updated_at,omitempty"`
	DeletedAt   *int         `json:"deleted_at,omitempty"`
	DeletedFor  *string      `json:"deleted_for,omitempty"`
	IsOwner     *bool        `json:"is_owner,omitempty"`
}

type MessagesResponse struct {
	Data  []*Message `json:"data,omitempty"`
	Error *Error     `json:"error,omitempty"`
}

type Mutation struct {
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type Query struct {
}

type Room struct {
	ID         *string   `json:"id,omitempty"`
	Title      *string   `json:"title,omitempty"`
	Image      *string   `json:"image,omitempty"`
	Biography  *string   `json:"biography,omitempty"`
	Creator    *string   `json:"creator,omitempty"`
	CreatorID  *string   `json:"creator_id,omitempty"`
	RoomType   *RoomType `json:"room_type,omitempty"`
	IsReported *bool     `json:"is_reported,omitempty"`
	IsActive   *bool     `json:"is_active,omitempty"`
	CreatedAt  *int      `json:"created_at,omitempty"`
	UpdatedAt  *int      `json:"updated_at,omitempty"`
	DeletedAt  *int      `json:"deleted_at,omitempty"`
	DeletedFor *string   `json:"deleted_for,omitempty"`
	IsMember   *bool     `json:"is_member,omitempty"`
	IsAdmin    *bool     `json:"is_admin,omitempty"`
}

type RoomMembership struct {
	ID        *string `json:"id,omitempty"`
	RoomID    *string `json:"room_id,omitempty"`
	JoinBy    *string `json:"join_by,omitempty"`
	Username  *string `json:"username,omitempty"`
	UserID    *string `json:"user_id,omitempty"`
	IsAdmin   *bool   `json:"is_admin,omitempty"`
	MuteUntil *int    `json:"mute_until,omitempty"`
	JoinAt    *int    `json:"join_at,omitempty"`
	LeaveAt   *int    `json:"leave_at,omitempty"`
}

type RoomMembershipResponse struct {
	Data  []*RoomMembership `json:"data,omitempty"`
	Error *Error            `json:"error,omitempty"`
}

type RoomResponse struct {
	Data  *Room  `json:"data,omitempty"`
	Error *Error `json:"error,omitempty"`
}

type RoomSummery struct {
	ID       string   `json:"id"`
	RoomType RoomType `json:"room_type"`
	Title    string   `json:"title"`
	Image    *string  `json:"image,omitempty"`
}

type RoomSummeryResponse struct {
	Data  []*RoomSummery `json:"data,omitempty"`
	Error *Error         `json:"error,omitempty"`
}

type UserBlock struct {
	ID          *string `json:"id,omitempty"`
	BlockerUser *string `json:"blocker_user,omitempty"`
	BlockedUser *string `json:"blocked_user,omitempty"`
	Description *string `json:"description,omitempty"`
	CreatedAt   *int    `json:"created_at,omitempty"`
}

type UserInput struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type AddRoomMembersInput struct {
	RoomID string       `json:"room_id"`
	Users  []*UserInput `json:"users"`
}

type CreateRoomInput struct {
	Title     *string      `json:"title,omitempty"`
	RoomType  RoomType     `json:"room_type"`
	Users     []*UserInput `json:"users,omitempty"`
	Biography *string      `json:"biography,omitempty"`
}

type JoinRoomInput struct {
	RoomID string `json:"room_id"`
}

type LeaveRoomInput struct {
	RoomID string `json:"room_id"`
}

type RemoveRoomMemberInput struct {
	RoomID string     `json:"room_id"`
	User   *UserInput `json:"user"`
}

type MessageType string

const (
	MessageTypeText      MessageType = "text"
	MessageTypeLink      MessageType = "link"
	MessageTypeLocation  MessageType = "location"
	MessageTypeImage     MessageType = "image"
	MessageTypeVoice     MessageType = "voice"
	MessageTypeVideo     MessageType = "video"
	MessageTypePDF       MessageType = "pdf"
	MessageTypeCall      MessageType = "call"
	MessageTypeVideoCall MessageType = "video_call"
	MessageTypeVoiceCall MessageType = "voice_call"
	MessageTypeOthers    MessageType = "others"
)

var AllMessageType = []MessageType{
	MessageTypeText,
	MessageTypeLink,
	MessageTypeLocation,
	MessageTypeImage,
	MessageTypeVoice,
	MessageTypeVideo,
	MessageTypePDF,
	MessageTypeCall,
	MessageTypeVideoCall,
	MessageTypeVoiceCall,
	MessageTypeOthers,
}

func (e MessageType) IsValid() bool {
	switch e {
	case MessageTypeText, MessageTypeLink, MessageTypeLocation, MessageTypeImage, MessageTypeVoice, MessageTypeVideo, MessageTypePDF, MessageTypeCall, MessageTypeVideoCall, MessageTypeVoiceCall, MessageTypeOthers:
		return true
	}
	return false
}

func (e MessageType) String() string {
	return string(e)
}

func (e *MessageType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = MessageType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid MessageType", str)
	}
	return nil
}

func (e MessageType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type RoomType string

const (
	RoomTypeMutual  RoomType = "MUTUAL"
	RoomTypeGroup   RoomType = "GROUP"
	RoomTypeChannel RoomType = "CHANNEL"
)

var AllRoomType = []RoomType{
	RoomTypeMutual,
	RoomTypeGroup,
	RoomTypeChannel,
}

func (e RoomType) IsValid() bool {
	switch e {
	case RoomTypeMutual, RoomTypeGroup, RoomTypeChannel:
		return true
	}
	return false
}

func (e RoomType) String() string {
	return string(e)
}

func (e *RoomType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = RoomType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid RoomType", str)
	}
	return nil
}

func (e RoomType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}