package handlers

import (
	"forum/configs"
	"forum/internal/render"
	"forum/internal/service"
)

type Handler struct {
	service      *service.Service
	templates    render.TemplatesHTML
	googleConfig configs.GoogleConfig
	githubConfig configs.GithubConfig
}

func NewHandler(service *service.Service, tmlp render.TemplatesHTML, googc configs.GoogleConfig, gitc configs.GithubConfig) *Handler {
	return &Handler{
		service:      service,
		templates:    tmlp,
		googleConfig: googc,
		githubConfig: gitc,
	}
}
