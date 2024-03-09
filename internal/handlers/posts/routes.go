package posts

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) Routes() http.Handler {
	router := mux.NewRouter()

	router.Handle("/create", http.HandlerFunc(h.createPost)).Methods(http.MethodPost)
	router.Handle("/reaction", http.HandlerFunc(h.reactionPost)).Methods(http.MethodPost)
	router.Handle("/edit", http.HandlerFunc(h.edit)).Methods(http.MethodPut)
	router.Handle("/delete", http.HandlerFunc(h.delete)).Methods(http.MethodDelete)

	return router
}
