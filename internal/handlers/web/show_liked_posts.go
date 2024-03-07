package web

import (
	"forum/internal/handlers/utils"
	"forum/internal/render"
	"net/http"
	"strconv"
)

func (h *Handler) likedPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/likedposts" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := utils.GetUserFromContext(r)
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		h.service.Log.Println(err)

		limit = 10
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		h.service.Log.Println(err)

		offset = 0
	}
	posts, err := h.service.PostService.GetLikedPosts(user.ID, offset, limit)

	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = h.service.PostReactionService.GetAllPostReactionsByPostID(posts)
	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	categories, err := h.service.CategoryService.GetAllCategories()

	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	h.templates.Render(w, r, "home.page.html", &render.PageData{
		Topic:             "LikedPosts",
		Categories:        categories,
		Posts:             posts,
		AuthenticatedUser: user,
	})

}

func (h *Handler) getLikedPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/lp" {
		h.service.Log.Println(r.URL.Path)

		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	isUserGay := r.Header.Get("INFINITE-SCROLL")
	if len(isUserGay) == 0 {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		h.service.Log.Println(err)

		limit = 10
	}
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		h.service.Log.Println(err)

		offset = 0
	}
	id := utils.GetUserFromContext(r)
	posts, err := h.service.PostService.GetLikedPosts(id.ID, offset, limit)
	err = h.service.PostReactionService.GetAllPostReactionsByPostID(posts)
	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err != nil {
		h.service.Log.Println(err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	h.templates.Render(w, r, "posts.page.html", &render.PageData{Posts: posts})
}
