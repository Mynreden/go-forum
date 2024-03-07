package web

import (
	"forum/internal/handlers/utils"
	"net/http"
	"strconv"

	"forum/internal/render"
)

func (h *Handler) showPostsByCategory(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/postscat" {
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	category := r.URL.Query().Get("category")

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

	posts, err := h.service.PostService.GetPostsByCategory(category, offset, limit)
	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		Topic:             category,
		Posts:             posts,
		Categories:        categories,
		AuthenticatedUser: utils.GetUserFromContext(r),
	})
}

func (h *Handler) getPostsCat(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/pc" {
		h.service.Log.Println(r.URL.Path)

		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	isUserGay := r.Header.Get("INFINITE-SCROLL")
	if len(isUserGay) == 0 {
		http.Error(w, "", http.StatusNotFound)
		return
	}
	category := r.URL.Query().Get("category")

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
	posts, err := h.service.PostService.GetPostsByCategory(category, offset, limit)
	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	err = h.service.PostReactionService.GetAllPostReactionsByPostID(posts)
	if err != nil {
		h.service.Log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	h.templates.Render(w, r, "posts.page.html", &render.PageData{Posts: posts})
}
