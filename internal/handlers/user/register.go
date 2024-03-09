package user

import (
	"forum/internal/domain"
	"forum/internal/render"
	"forum/pkg/forms"
	"net/http"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	err := r.ParseForm()
	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("username", "email", "password", "rpass")
	form.MaxLength("username", 50)
	form.MaxLength("email", 50)
	form.MatchesPattern("email", forms.EmailRX)
	form.MaxLength("password", 50)
	form.MinLength("password", 8)
	if r.FormValue("rpass") != r.FormValue("password") {
		form.Errors.Add("password", "Pas credentials")
	}
	if !form.Valid() {
		form.Errors.Add("generic", "Passwords don't match")
		w.WriteHeader(http.StatusBadRequest)
		h.templates.Render(w, r, "reg.page.html", &render.PageData{
			Form: form,
		})
		return
	}

	req := &domain.CreateUserDTO{
		Email:    r.FormValue("email"),
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}

	err = h.service.UserService.CreateUser(req)

	if err != nil {
		h.service.Log.Println(err)
		switch err {
		case domain.ErrDuplicateEmail:
			form.Errors.Add("email", "Email already in use")
			w.WriteHeader(http.StatusBadRequest)
			h.templates.Render(w, r, "reg.page.html", &render.PageData{
				Form: form,
			})
			return
		case domain.ErrDuplicateUsername:
			form.Errors.Add("username", "Username already in use")
			w.WriteHeader(http.StatusBadRequest)
			h.templates.Render(w, r, "reg.page.html", &render.PageData{
				Form: form,
			})
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)

	return
}
