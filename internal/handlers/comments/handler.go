package comments

import (
	"forum/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewCommentsHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}
