package web

import (
	"forum/internal/handlers/utils"
	"forum/internal/render"
	"forum/pkg/forms"
	"net/http"
)

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	h.templates.Render(w, r, "log.page.html", &render.PageData{
		Form:              forms.New(nil),
		AuthenticatedUser: utils.GetUserFromContext(r),
	})
}
