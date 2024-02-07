package handlers

import (
	"forum/internal/render"
	"forum/pkg/forms"
	"net/http"
	"strconv"
)

func (h *Handler) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	form := forms.New(r.PostForm)
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

	// limit := 10
	// offset := 0
	posts, err := h.service.PostService.GetAllPosts(offset, limit)

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
	categories, err := h.service.CategoryService.GetAllCategories()

	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	h.templates.Render(w, r, "home.page.html", &render.PageData{
		Topic:             "News",
		Categories:        categories,
		Form:              form,
		Posts:             posts,
		AuthenticatedUser: h.getUserFromContext(r),
	})
}



func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/posts" {
		h.service.Log.Println(r.URL.Path)

		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Incorrect Method", http.StatusMethodNotAllowed)
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
	posts, err := h.service.PostService.GetAllPosts(offset, limit)
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
