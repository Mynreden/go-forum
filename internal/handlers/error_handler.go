package handlers

import "net/http"

func (h *Handler) errorHandler(w http.ResponseWriter, r *http.Request, status int, err error) {
	w.WriteHeader(status)
	h.templates.Render(w, r, "error.page.html", nil)
}
