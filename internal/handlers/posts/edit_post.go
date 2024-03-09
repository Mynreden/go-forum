package posts

import (
	"fmt"
	"forum/internal/domain"
	"forum/internal/handlers/utils"
	"forum/internal/render"
	"forum/pkg/forms"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) edit(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/edit" {
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

		http.Error(w, "Pasre error", http.StatusInternalServerError)
		return
	}

	if err := r.ParseMultipartForm(20 << 20); err != nil {
		h.service.Log.Println(err)

		http.Error(w, "Invalid POST request", http.StatusInternalServerError)
		return
	}

	form := forms.New(r.PostForm)

	form.Required("title", "content")
	form.MaxLength("title", 100)
	form.MaxLength("content", 500)

	if !form.Valid() {

		categories, err := h.service.CategoryService.GetAllCategories()
		if err != nil {
			h.service.Log.Println(err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		h.templates.Render(w, r, "error.page.html", &render.PageData{
			Form:              form,
			Categories:        categories,
			AuthenticatedUser: utils.GetUserFromContext(r),
		})

		return
	}

	autor := utils.GetUserFromContext(r)

	post := &domain.CreatePostDTO{
		Title:      r.PostFormValue("title"),
		Content:    r.PostFormValue("content"),
		Author:     autor.ID,
		AuthorName: autor.Username,
		// Categories: categories,
	}

	file, fileHeader, err := r.FormFile("image")

	if err != nil {
		h.service.Log.Println(err)

		if err != http.ErrMissingFile {

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		post.ImageFile = nil
	} else {
		post.ImageFile = file
		defer file.Close()

		fileType := fileHeader.Header.Get("Content-Type")
		if !form.IsImg(fileType) {
			categories, err := h.service.CategoryService.GetAllCategories()
			if err != nil {
				log.Println(err)

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			form.Errors.Add("image", "File is not an image")
			form.Categories = append(form.Categories, categories...)

			w.WriteHeader(http.StatusBadRequest)
			h.templates.Render(w, r, "error.page.html", &render.PageData{
				Form:              form,
				Categories:        categories,
				AuthenticatedUser: utils.GetUserFromContext(r),
			})
			return
		}

		if fileHeader.Size > 5*1024*1024 {
			categories, err := h.service.CategoryService.GetAllCategories()
			if err != nil {
				h.service.Log.Println(err)

				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			form.Categories = append(form.Categories, categories...)
			form.Errors.Add("image", "File is too big")

			w.WriteHeader(http.StatusBadRequest)
			h.templates.Render(w, r, "error.page.html", &render.PageData{
				Form:              form,
				Categories:        categories,
				AuthenticatedUser: utils.GetUserFromContext(r),
			})
			return
		}
	}

	categoriesS := r.PostFormValue("category")
	if len(categoriesS) == 0 {

		categories, err := h.service.CategoryService.GetAllCategories()
		if err != nil {
			h.service.Log.Println(err)

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		form.Errors.Add("generic", "You must select at least one category")
		form.Categories = append(form.Categories, categories...)

		w.WriteHeader(http.StatusBadRequest)
		h.templates.Render(w, r, "error.page.html", &render.PageData{
			Form:              form,
			Categories:        categories,
			AuthenticatedUser: utils.GetUserFromContext(r),
		})
		return
	}

	tempD := strings.Split(categoriesS, ",")
	for i, v := range tempD {
		tempD[i] = strings.TrimSpace(v)
	}
	categories := make([]*domain.Category, 0, len(tempD))
	for _, name := range tempD {
		c, err := h.service.CategoryService.GetCategoryByName(name)
		if err != nil {
			h.service.Log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			categories = append(categories, c)
		}
	}

	post.Categories = append(post.Categories, categories...)

	post_id, err := h.service.PostService.EditPostWithImage(post, int(id))
	if err != nil {
		h.service.Log.Println(err)

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/post?id=%d", post_id), http.StatusFound)

}
