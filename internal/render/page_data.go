package render

import (
	"forum/internal/models"
	"forum/pkg/forms"
)

type PageData struct {
	Topic             string
	Form              *forms.Form
	AuthenticatedUser *models.User
	Post              *models.Post
	Posts             []*models.Post
	Categories        []*models.Category
	Comments          []*models.Comment
}
