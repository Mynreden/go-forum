package handlers

import (
	"forum/internal/helpers/cookies"
	"forum/internal/models"
	"forum/internal/render"
	"forum/pkg/forms"
	"net/http"
	"time"
)

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.templates.Render(w, r, "log.page.html", &render.PageData{
			Form:              forms.New(nil),
			AuthenticatedUser: h.getUserFromContext(r),
		})
	case http.MethodPost:

		err := r.ParseForm()
		if err != nil {
			h.service.Log.Println(err)

			http.Error(w, "Invalid POST request", http.StatusInternalServerError)
			return
		}

		form := forms.New(r.PostForm)
		form.Required("email", "password")
		form.MaxLength("email", 50)
		form.MatchesPattern("email", forms.EmailRX)
		form.MaxLength("password", 50)
		form.MinLength("password", 8)

		if !form.Valid() {
			w.WriteHeader(http.StatusBadRequest)
			h.templates.Render(w, r, "log.page.html", &render.PageData{
				Form: form,
			})
			return
		}

		req := &models.LoginUserDTO{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		user_id, err := h.service.UserService.LoginUser(req)
		if err != nil {
			h.service.Log.Println(err)
			if err == models.ErrInvalidCredentials {

				form.Errors.Add("generic", "Email or password is incorrect")
				h.templates.Render(w, r, "log.page.html", &render.PageData{
					Form: form,
				})
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := h.service.SessionService.CreateSession(user_id)
		if err != nil {
			h.service.Log.Println(err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cookies.SetCookie(w, session.UUID, int(time.Until(session.ExpireAt).Seconds()))

		http.Redirect(w, r, "/", http.StatusFound)

	default:
		http.Error(w, "incorrect Method", http.StatusMethodNotAllowed)
	}
	return
}
