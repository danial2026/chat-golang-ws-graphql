package handlers

import (
	"net/http"
)

type MessageResponse struct {
	StatusCode int
	Body       string
}

func BadRequestErrorResponse(msg string) MessageResponse {
	return MessageResponse{
		StatusCode: http.StatusBadRequest,
		Body:       msg,
	}
}

func ForbiddenErrorResponse(msg string) MessageResponse {
	return MessageResponse{
		StatusCode: http.StatusForbidden,
		Body:       msg,
	}
}

func NotFoundErrorResponse(msg string) MessageResponse {
	return MessageResponse{
		StatusCode: http.StatusNotFound,
		Body:       msg,
	}
}

func OkResponse(msg string) MessageResponse {
	return MessageResponse{
		StatusCode: http.StatusOK,
		Body:       msg,
	}
}

func InternalServerErrorResponse(msg string) MessageResponse {
	return MessageResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       msg,
	}
}
