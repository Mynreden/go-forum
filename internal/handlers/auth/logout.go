package auth

import (
	"forum/internal/helpers/cookies"
	"net/http"
)

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/logout" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := cookies.GetCookie(r)
	if err != nil {
		h.service.Log.Println(err)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.SessionService.DeleteSessionByUUID(cookie.Value)
	if err != nil {
		h.service.Log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cookies.DeleteCookie(w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
