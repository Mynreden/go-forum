package comments

import (
	"bytes"
	"context"
	"fmt"
	"forum/internal/domain"
	"forum/internal/handlers/utils"
	"forum/internal/repository"
	"forum/internal/service"
	"forum/pkg/forms"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestReactionCommentHandler(t *testing.T) {
	// Создаем мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' occurred while opening a stub database connection", err)
	}
	defer db.Close()

	addMockData(mock)

	// Создаем мок репозитория и сервиса
	mockRepository := repository.NewRepository(db)

	logger := log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	mockService := service.NewService(mockRepository, logger)

	handler := NewCommentsHandler(mockService)

	tests := []struct {
		name           string
		method         string
		postID         string
		status         string
		commentID      string
		expectedStatus int
	}{
		{"ValidCredentials", http.MethodPost, "1", "1", "1", http.StatusSeeOther},
		{"InvalidCredentials", http.MethodPost, "100", "1", "1", http.StatusBadRequest},
		{"IncorrectMethod", http.MethodGet, "1", "1", "1", http.StatusMethodNotAllowed},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			formData := fmt.Sprintf("post_id=%s&comment_id=%s&status=%s", tc.postID, tc.commentID, tc.status)
			req, err := http.NewRequest(tc.method, "/comment/reaction", bytes.NewBufferString(formData))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			req.ParseForm()

			user := &domain.User{1,
				"mynreden",
				"password",
				"email@gmail.com",
				time.Now(),
				time.Now()}

			ctx := context.WithValue(req.Context(), utils.ContextKeyUser, user)
			req = req.WithContext(ctx)

			form := forms.New(req.PostForm)

			form.Required("comment_id", "status", "post_id")
			postID := form.IsInt("post_id")
			id := form.IsInt("comment_id")
			t.Log(postID, id)
			rr := httptest.NewRecorder()

			handler.reactionComment(rr, req)

			// Check status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tc.expectedStatus, rr.Code)
			}

			// Check cookie is set on successful login
			if tc.expectedStatus == http.StatusFound {
				cookies := rr.Result().Cookies()
				if len(cookies) != 1 || cookies[0].Name != "session_token" || cookies[0].Value != "test-uuid" {
					t.Errorf("Expected session cookie 'session_token=test-uuid', got %v", cookies)
				}
			}
		})
	}
}

func addMockData(mock sqlmock.Sqlmock) {
	createdAt, _ := time.Parse(time.RFC3339Nano, "2024-02-01T09:20:52.575372164+06:00")
	updatedAt, _ := time.Parse(time.RFC3339Nano, "2024-02-01 09:20:52.575372345+06:00")
	postRows := mock.NewRows([]string{"id", "title", "content", "author_id", "author_name", "created_at", "updated_at"}).
		AddRow(1, "title", "bla bla bla", 2, "styan", createdAt, updatedAt)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM posts WHERE id = $1`)).WithArgs(1).WillReturnRows(postRows)

	categoryRows := mock.NewRows([]string{"category_name"}).
		AddRow("Dota").
		AddRow("Books")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT category_name FROM PostCategories WHERE post_id = $1`)).WithArgs(1).WillReturnRows(categoryRows)

	imageRows := mock.NewRows([]string{"image_path"}).
		AddRow("image.png")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT image_path FROM images WHERE post_id = $1`)).WithArgs(1).WillReturnRows(imageRows)

	// get comment
	commentRows := mock.NewRows([]string{"id", "comment", "post_id", "user_id", "userName", "created_at"}).
		AddRow(1, "comment text", 1, 1, "sultan", createdAt)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, comment, post_id, user_id, userName, created_at FROM comments WHERE id = ?`)).WithArgs(1).WillReturnRows(commentRows)

}
