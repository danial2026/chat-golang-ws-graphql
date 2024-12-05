package handlers

import (
	"context"
	"log"

	"chat-golang-react/chat/models"
	"chat-golang-react/chat/websocket/pushers"
	"chat-golang-react/chat/websocket/resources"
)

func ConnectHandler(ctx context.Context, request QueryStringParameters, res *resources.Resource) MessageResponse {
	log.Println("ConnectHandler username is: ", res.User.Username)

	// validate user inputs:
	connOut, validErr := ValidateConnect(request, res.User.Username, res.User.UserId)
	log.Println("ConnectHandler/ValidateConnect output is: ", connOut)
	if validErr != nil {
		log.Println("ConnectHandler/ValidateConnect error is: ", validErr.Error())
		return BadRequestErrorResponse(validErr.Error())
	}

	conn := connOut.Connection
	operation := connOut.Operation

	// set connectionID in User
	res.User.ConnectionId = conn.ConnectionID

	// handle new connection to the room
	if operation == ConnectToRoom {

		if conn.RoomId != "" {
			// check if room exists:
			roomOk, existErr := models.IsRoomExists(ctx, *res, conn.RoomId)
			if existErr != nil {
				log.Println("ConnectHandler/IsRoomExists error is: ", existErr.Error())
				return InternalServerErrorResponse(existErr.Error())
			}

			if !roomOk {
				msgErr := "room not found"
				return NotFoundErrorResponse(msgErr)
			}
		}

		// create connection:
		createErr := models.CreateConnection(ctx, *res, conn.ConnectionID, res.User.Username, res.User.UserId, conn.RoomId)
		if createErr != nil {
			log.Println("ConnectHandler/CreateConnection error is: ", createErr.Error())
			return InternalServerErrorResponse(createErr.Error())
		}
	} else {
		// TODO : uncomment me after making sure i work
		// handle conversion list connection:
		// creationErr := models.CreateOrUpdateConversationConnection(ctx, res, conn.ConnectionID, res.User.Username, res.User.UserId)
		// if creationErr != nil {
		// 	return InternalServerErrorResponse(creationErr.Error())
		// }
		log.Println("[!] ConnectHandler/else")
	}

	// TODO update user online status
	return OkResponse("successfully connected to the room")
}

func DisconnectHandler(ctx context.Context, res resources.Resource, roomId string) MessageResponse {
	err := models.DisconnectConnection(ctx, res, roomId)
	if err != nil {
		return InternalServerErrorResponse(err.Error())
	}

	// TODO update user online status
	return OkResponse("connectionId is successfully disconnected")
}

func SendMessageHandler(ctx context.Context, request InputMessageBody, res resources.Resource) MessageResponse {
	msgOut, validErr := ValidateSendMessage(request, res.User.Username, res.User.UserId, res.RoomID)
	if validErr != nil {
		return BadRequestErrorResponse(validErr.Error())
	}
	log.Println("SendMessageHandler/ValidateSendMessage msgOut: ", msgOut)
	msg := msgOut.Message

	// check existence of the room:
	isExist, existErr := models.IsRoomExists(ctx, res, msg.RoomId)
	if existErr != nil {
		return InternalServerErrorResponse(existErr.Error())
	}
	if !isExist {
		msgErr := "room not found"
		return NotFoundErrorResponse(msgErr)
	}

	// check user membership:
	isMember, memberErr := models.IsRoomMember(ctx, res, res.User.Username, res.User.UserId, msg.RoomId)
	log.Println("SendMessageHandler/IsRoomMember: requested user: ", res.User.Username, " roomId: ", msg.RoomId)
	if memberErr != nil {
		log.Println("SendMessageHandler/IsRoomMember/memberErr: ", memberErr)
		return InternalServerErrorResponse(memberErr.Error())
	}
	if !isMember {
		log.Println("SendMessageHandler/IsRoomMember/!isMember")
		msgErr := "user is not a member"
		return ForbiddenErrorResponse(msgErr)
	}

	// create message:
	if msgOut.Operation == CreateMessageOperation {
		newMsg, creationErr := models.CreateMessage(ctx, res, res.User.Username, res.User.UserId, msg)
		if creationErr != nil {
			log.Println("SendMessageHandler/CreateMessage/creationErr: ", creationErr)
			return InternalServerErrorResponse(creationErr.Error())
		}
		// send message to the online members:
		_, sentErr := pushers.PushNewMessage(
			ctx,
			res,
			pushers.PushInput{
				Action:  pushers.SendToRoomPush,
				RoomId:  msg.RoomId,
				Message: newMsg,
			})
		if sentErr != nil {
			log.Println("SendMessageHandler/PushNewMessage/sentErr: ", sentErr)
			msgErr := "new message could not be pushed to the client"
			return InternalServerErrorResponse(msgErr)
		}
		// send message to room members conversion list room
		//allMembers, getErr := models.GetMembersByRoomId(ctx, res, msg.RoomId)
		//if getErr != nil {
		//	return InternalServerErrorResponse(getErr.Error())
		//}
		//conversationListConns, queryErr := GetNotSentUsers(ctx, res, sentConns, allMembers)
		//if queryErr != nil {
		//	return InternalServerErrorResponse(queryErr.Error())
		//}
		//_, pushErr := pushers.PushNewMessage(
		//	ctx,
		//	res,
		//	pushers.PushInput{
		//		Action:         pushers.SendToConversionListPush,
		//		ConnectionList: conversationListConns,
		//		Message:        newMsg,
		//	})
		//if pushErr != nil {
		//	log.Println("SendMessageHandler/PushNewMessage/sentErr: ", pushErr)
		//	msgErr := "new message could not be pushed to the client"
		//	return InternalServerErrorResponse(msgErr)
		//}
		// _, pushErr := notifs.PushNewNotification(
		// 	ctx,
		// 	res,
		// 	notifs.PushInput{
		// 		Action:         pushers.SendToConversionListPush,
		// 		ConnectionList: conversationListConns,
		// 		Message:        newMsg,
		// 	})
		// if pushErr != nil {
		// 	log.Println("SendMessageHandler/PushNewNotification/sentErr: ", pushErr)
		// 	msgErr := "new message could not be pushed to the client"
		// 	return InternalServerErrorResponse(msgErr)
		// }
		// delete the message:
	} else if msgOut.Operation == DeleteMessageOperation {
		// check message ownership"
		_, hasPerm, permErr := models.IsMessageOwner(ctx, res, msg.Id, msg.RoomId, res.User.Username, res.User.UserId)
		if permErr != nil {
			return InternalServerErrorResponse(permErr.Error())
		}
		if !hasPerm {
			msgErr := "permission denied"
			return ForbiddenErrorResponse(msgErr)
		}
		// delete the message:
		newMsg, deleteErr := models.DeleteMessageForAll(ctx, res, msg.Id, msg.RoomId, res.User.Username, res.User.UserId)
		if deleteErr != nil {
			return InternalServerErrorResponse(deleteErr.Error())
		}
		// send message to the online members:
		sentConns, sentErr := pushers.PushNewMessage(
			ctx,
			res,
			pushers.PushInput{
				Action:  pushers.SendToRoomPush,
				RoomId:  msg.RoomId,
				Message: newMsg,
			})
		if sentErr != nil {
			log.Println("SendMessageHandler/PushNewMessage/sentErr: ", sentErr)
			msgErr := "new message could not be pushed to the client"
			return InternalServerErrorResponse(msgErr)
		}
		// send message to room members conversion list room
		allMembers, getErr := models.GetMembersByRoomId(ctx, res, msg.RoomId)
		if getErr != nil {
			return InternalServerErrorResponse(getErr.Error())
		}
		conversationListConns, queryErr := GetNotSentUsers(ctx, res, sentConns, allMembers)
		if queryErr != nil {
			return InternalServerErrorResponse(queryErr.Error())
		}
		_, pushErr := pushers.PushNewMessage(
			ctx,
			res,
			pushers.PushInput{
				Action:         pushers.SendToConversionListPush,
				ConnectionList: conversationListConns,
				Message:        newMsg,
			},
		)
		if pushErr != nil {
			log.Println("SendMessageHandler/PushNewMessage/sentErr: ", pushErr)
			msgErr := "new message could not be pushed to the client"
			return InternalServerErrorResponse(msgErr)
		}
	} else if msgOut.Operation == DeleteMessageForAllOperation {
		// check message ownership:
		_, hasPerm, permErr := models.IsMessageOwner(ctx, res, msg.Id, msg.RoomId, res.User.Username, res.User.UserId)
		if permErr != nil {
			return InternalServerErrorResponse(permErr.Error())
		}
		if !hasPerm {
			msgErr := "permission denied"
			return ForbiddenErrorResponse(msgErr)
		}

		// delete message:
		newMsg, deleteErr := models.DeleteMessageForAll(ctx, res, msg.Id, msg.RoomId, res.User.Username, res.User.UserId)
		if deleteErr != nil {
			return InternalServerErrorResponse(deleteErr.Error())
		}

		// send message to the online members:
		sentConns, sentErr := pushers.PushNewMessage(
			ctx,
			res,
			pushers.PushInput{
				Action:  pushers.SendToRoomPush,
				RoomId:  msg.RoomId,
				Message: newMsg,
			})
		if sentErr != nil {
			log.Println("SendMessageHandler/PushNewMessage/sentErr: ", sentErr)
			msgErr := "new message could not be pushed to the client"
			return InternalServerErrorResponse(msgErr)
		}

		// send message to room members conversion list room
		allMembers, getErr := models.GetMembersByRoomId(ctx, res, msg.RoomId)
		if getErr != nil {
			return InternalServerErrorResponse(getErr.Error())
		}
		conversationListConns, queryErr := GetNotSentUsers(ctx, res, sentConns, allMembers)
		if queryErr != nil {
			return InternalServerErrorResponse(queryErr.Error())
		}
		_, pushErr := pushers.PushNewMessage(
			ctx,
			res,
			pushers.PushInput{
				Action:         pushers.SendToConversionListPush,
				ConnectionList: conversationListConns,
				Message:        newMsg,
			})
		if pushErr != nil {
			log.Println("SendMessageHandler/PushNewMessage/sentErr: ", pushErr)
			msgErr := "new message could not be pushed to the client"
			return InternalServerErrorResponse(msgErr)
		}
	} else if msgOut.Operation == DeleteMessageForMeOperation {
		// check message ownership:
		_, hasPerm, permErr := models.IsMessageOwner(ctx, res, msg.Id, msg.RoomId, res.User.Username, res.User.UserId)
		if permErr != nil {
			return InternalServerErrorResponse(permErr.Error())
		}
		if !hasPerm {
			msgErr := "permission denied"
			return ForbiddenErrorResponse(msgErr)
		}

		// delete message:
		newMsg, deleteErr := models.DeleteMessageForMe(ctx, res, msg.Id, msg.RoomId, res.User.Username, res.User.UserId)
		if deleteErr != nil {
			return InternalServerErrorResponse(deleteErr.Error())
		}

		// send message to the online members:
		sentConns, sentErr := pushers.PushNewMessage(
			ctx,
			res,
			pushers.PushInput{
				Action:  pushers.SendToRoomPush,
				RoomId:  msg.RoomId,
				Message: newMsg,
			})
		if sentErr != nil {
			log.Println("SendMessageHandler/PushNewMessage/sentErr: ", sentErr)
			msgErr := "new message could not be pushed to the client"
			return InternalServerErrorResponse(msgErr)
		}

		// send message to room members conversion list room
		allMembers, getErr := models.GetMembersByRoomId(ctx, res, msg.RoomId)
		if getErr != nil {
			return InternalServerErrorResponse(getErr.Error())
		}
		conversationListConns, queryErr := GetNotSentUsers(ctx, res, sentConns, allMembers)
		if queryErr != nil {
			return InternalServerErrorResponse(queryErr.Error())
		}
		_, pushErr := pushers.PushNewMessage(
			ctx,
			res,
			pushers.PushInput{
				Action:         pushers.SendToConversionListPush,
				ConnectionList: conversationListConns,
				Message:        newMsg,
			})
		if pushErr != nil {
			log.Println("SendMessageHandler/PushNewMessage/sentErr: ", pushErr)
			msgErr := "new message could not be pushed to the client"
			return InternalServerErrorResponse(msgErr)
		}
	}

	return OkResponse("message successfully handeled")
}

func DefaultHandler() MessageResponse {
	msgBody := "unrecognized WebSocket action"
	return BadRequestErrorResponse(msgBody)
}
