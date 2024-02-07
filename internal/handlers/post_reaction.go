package handlers

import (
	"fmt"
	"forum/internal/models"
	"net/http"
	"strconv"
)

func (h *Handler) reactionPost(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/post/reaction" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Incorrect Method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		h.service.Log.Printf("Error parsing form: %v\n", err)
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	userID := h.getUserFromContext(r)

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil {
		h.service.Log.Printf("Error converting post_id: %v\n", err)
		http.Error(w, "Invalid post_id", http.StatusBadRequest)
		return
	}

	_, err = h.service.PostService.GetPostByID(postID)
	if err != nil {
		h.service.Log.Println(err)
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}
	status, err := strconv.Atoi(r.FormValue("status"))
	if err != nil {
		h.service.Log.Printf("Error converting status: %v\n", err)
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	switch status {
	case 1:
		// Status is true
	case 0:
		// Status is false
	default:
		h.service.Log.Printf("Invalid status value")
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	_, err = h.service.PostService.GetPostByID(postID)
	if err != nil {
		h.service.Log.Println(err)
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	reaction := &models.PostReactionDTO{
		UserID: userID.ID,
		PostID: postID,
		Status: status == 1,
	}

	err = h.service.PostReactionService.CreatePostReaction(reaction)
	if err != nil {

		h.service.Log.Printf("Error creating post reaction: %v\n", err)
		http.Error(w, "Error creating post reaction", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/post/?id=%d", postID), http.StatusFound)
}
