package handlers

import (
	"fmt"
	"forum/configs"
	"forum/internal/models"
	"forum/internal/render"
	"forum/internal/service"
	"io/ioutil"
	"net/http"
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

type contextKey string

var contextKeyUser = contextKey("user")

func (h *Handler) getUserFromContext(r *http.Request) *models.User {
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}

func (h *Handler) getUserInfo(accessToken string, userInfoURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", userInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
