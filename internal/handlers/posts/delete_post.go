package posts

import (
	"fmt"
	"forum/internal/handlers/utils"
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

	user := utils.GetUserFromContext(r)
	post, err := h.service.PostService.GetPostByID(int(id))

	if err != nil {
		h.service.Log.Println(err)
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if user.ID != post.AuthorID {
		h.service.Log.Println(err)
		http.Error(w, "You cannot delete foreign posts", http.StatusBadRequest)
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
