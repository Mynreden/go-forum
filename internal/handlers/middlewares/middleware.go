package middlewares

import (
	"context"
	"forum/internal/handlers/utils"
	"forum/internal/helpers/cookies"
	"forum/internal/service"
	"net/http"
	"time"
)

type middleware struct {
	service *service.Service
}

func NewMiddleware(service *service.Service) *middleware {
	return &middleware{service}
}

func (h *middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := cookies.GetCookie(r)
		if err != nil {

			h.service.Log.Println(err)
			next.ServeHTTP(w, r)
			return
		}

		session, err := h.service.SessionService.GetSessionByUUID(cookie.Value)
		if err != nil {
			h.service.Log.Println(err)

			cookies.DeleteCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		if session.ExpireAt.Before(time.Now()) {
			cookies.DeleteCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		user, err := h.service.UserService.GetUserByID(session.User_id)
		if err != nil {
			h.service.Log.Println(err)

			cookies.DeleteCookie(w)
			h.service.SessionService.DeleteSessionByUUID(cookie.Value)
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), utils.ContextKeyUser, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *middleware) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := utils.GetUserFromContext(r)
		if user == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}
