package posts

import (
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delete" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	urlId := r.URL.Query().Get("id")
	if urlId == "" {
		http.Error(w, "Invalid id", http.StatusInternalServerError)
		return
	}
	id, err := strconv.ParseInt(urlId, 10, 64)
	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, "Parse error", http.StatusInternalServerError)
		return
	}

	err = h.service.PostService.DeletePost(int(id))
	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/myposts"), http.StatusFound)
}
