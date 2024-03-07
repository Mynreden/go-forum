package render

import (
	"forum/internal/domain"
	"forum/pkg/forms"
)

type PageData struct {
	Topic             string
	Form              *forms.Form
	AuthenticatedUser *domain.User
	Post              *domain.Post
	Posts             []*domain.Post
	Categories        []*domain.Category
	Comments          []*domain.Comment
}
