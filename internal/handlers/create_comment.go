package handlers

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/forms"
	"net/http"
)

func (h *Handler) createComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/create" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("content", "post_id")
	form.MaxLength("content", 280)
	postID := form.IsInt("post_id")

	_, err = h.service.PostService.GetPostByID(postID)
	if err != nil {
		h.service.Log.Println(err)
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	if !form.Valid() {
		http.Redirect(w, r, fmt.Sprintf("/post/?id=%d", postID), http.StatusSeeOther)
		return
	}

	author := h.getUserFromContext(r)

	comment := &models.CreateCommentDTO{
		PostID:     postID,
		Content:    form.Get("content"),
		AuthorID:   author.ID,
		AuthorName: author.Username,
	}
	err = h.service.CommentService.CreateComment(comment)
	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/?id=%d", postID), http.StatusSeeOther)
}
