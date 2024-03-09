package user

import (
	"forum/internal/handlers/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) Routes() http.Handler {
	router := mux.NewRouter()
	middleware := middlewares.NewMiddleware(h.service)

	router.Handle("/logout", middleware.RequireAuthentication(http.HandlerFunc(h.logout))).Methods(http.MethodGet)
	router.HandleFunc("/login", h.login).Methods(http.MethodPost)
	router.HandleFunc("/register", h.register).Methods(http.MethodPost)
	//router.Handle("/edit", middleware.RequireAuthentication(http.HandlerFunc(h.edit))).Methods(http.MethodPut)

	return router
}
