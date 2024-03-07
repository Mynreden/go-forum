package web

import (
	"forum/internal/render"
	"forum/internal/service"
)

type Handler struct {
	service   *service.Service
	templates render.TemplatesHTML
}

func NewWebHandler(service *service.Service, tmlp render.TemplatesHTML) *Handler {
	return &Handler{
		service:   service,
		templates: tmlp,
	}
}
