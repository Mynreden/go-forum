package comments

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (h *Handler) Routes() http.Handler {
	router := mux.NewRouter()

	router.Handle("/create", http.HandlerFunc(h.createComment)).Methods(http.MethodPost)
	router.Handle("/reaction", http.HandlerFunc(h.reactionComment)).Methods(http.MethodPost)

	return router
}
