package handlers

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/forms"
	"net/http"
	"strconv"
)

func (h *Handler) reactionComment(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/comment/reaction" {
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

	form.Required("comment_id", "status", "post_id")
	postID := form.IsInt("post_id")
	id := form.IsInt("comment_id")
	status, err := strconv.Atoi(r.FormValue("status"))
	if err != nil {
		h.service.Log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !form.Valid() {
		http.Redirect(w, r, fmt.Sprintf("/post/?id=%d", postID), http.StatusSeeOther)
		return
	}

	author := h.getUserFromContext(r)

	if err != nil {
		h.service.Log.Println("Error converting status: %v", err)

		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	switch status {
	case 1:
		// Status is Like
	case 0:
		// Status is DisLike
	default:
		h.service.Log.Println("Invalid status value")
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}
	_, err = h.service.PostService.GetPostByID(postID)
	if err != nil {
		h.service.Log.Println(err)
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}
	_, err = h.service.CommentService.GetCommentByID(id)
	if err != nil {
		h.service.Log.Println(err)
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	vote := &models.CommentReactionDTO{
		CommentID: id,
		Status:    status == 1,
		UserID:    author.ID,
	}

	if err := h.service.CommentReactionService.CreateCommentsReactions(vote); err != nil {
		h.service.Log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/?id=%d", postID), http.StatusSeeOther)
}
