package web

import (
	"forum/internal/handlers/middlewares"
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) Routes() http.Handler {
	router := mux.NewRouter()
	middleware := middlewares.NewMiddleware(h.service)

	router.HandleFunc("/login", h.login).Methods(http.MethodGet)
	router.HandleFunc("/register", h.register).Methods(http.MethodGet)
	router.HandleFunc("/post", h.showPost).Methods(http.MethodGet)
	router.HandleFunc("/posts", h.getPosts).Methods(http.MethodGet)
	router.Handle("/lp", middleware.RequireAuthentication(http.HandlerFunc(h.getLikedPosts))).Methods(http.MethodGet)
	router.HandleFunc("/postscat", h.showPostsByCategory).Methods(http.MethodGet)
	router.HandleFunc("/pc", h.getPostsCat).Methods(http.MethodGet)
	router.Handle("/myposts", middleware.RequireAuthentication(http.HandlerFunc(h.myposts))).Methods(http.MethodGet)
	router.Handle("/mp", middleware.RequireAuthentication(http.HandlerFunc(h.getMyPosts))).Methods(http.MethodGet)
	router.Handle("/likedposts", middleware.RequireAuthentication(http.HandlerFunc(h.likedPosts))).Methods(http.MethodGet)
	router.HandleFunc("/", h.home).Methods(http.MethodGet)

	return router
}
