package web

import (
	"forum/internal/handlers/utils"
	"forum/internal/render"
	"forum/pkg/forms"
	"net/http"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	h.templates.Render(w, r, "reg.page.html", &render.PageData{
		Form:              forms.New(nil),
		AuthenticatedUser: utils.GetUserFromContext(r),
	})
	return
}
