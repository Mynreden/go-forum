package repository

import (
	"errors"
	"forum/internal/domain"
	"forum/internal/repository/post"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestPostStorage_CreatePost(t *testing.T) {
	testCases := []struct {
		name          string
		post          *domain.Post
		expectedID    int
		expectedError error
	}{
		{
			name: "Successful creation",
			post: &domain.Post{
				Title:      "Test Post",
				Content:    "Test Content",
				AuthorID:   1,
				AuthorName: "Test Author",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				Categories: []*domain.Category{
					{Name: "Category1"},
					{Name: "Category2"},
				},
			},
			expectedID:    0,
			expectedError: nil,
		},
		{
			name: "Error on category insertion",
			post: &domain.Post{
				Title:      "Test Post",
				Content:    "Test Content",
				AuthorID:   1,
				AuthorName: "Test Author",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				Categories: []*domain.Category{
					{Name: "Category1"},
					{Name: "Category2"},
				},
			},
			expectedID:    0,
			expectedError: errors.New("category insertion failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			mock.ExpectQuery("INSERT INTO posts").
				WithArgs(tc.post.Title, tc.post.Content, tc.post.AuthorID, tc.post.AuthorName, tc.post.CreatedAt, tc.post.UpdatedAt).
				WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(tc.expectedID, tc.post.CreatedAt, tc.post.UpdatedAt))

			for _, category := range tc.post.Categories {
				if tc.expectedError == nil {
					mock.ExpectExec("INSERT INTO PostCategories").
						WithArgs(tc.expectedID, category.Name).
						WillReturnResult(sqlmock.NewResult(0, 1))
				} else {
					mock.ExpectExec("INSERT INTO PostCategories").
						WithArgs(tc.expectedID, category.Name).
						WillReturnError(tc.expectedError)
				}
			}

			repo := post.NewPostStorage(db)
			id, err := repo.CreatePost(tc.post)

			assert.Equal(t, tc.expectedID, id)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCreatePostWithImage(t *testing.T) {
	testCases := []struct {
		name          string
		post          *domain.Post
		expectedID    int
		expectedError error
	}{
		{
			name: "Successful insertion with categories and image",
			post: &domain.Post{
				Title:      "Test Title",
				Content:    "Test Content",
				AuthorID:   100,
				AuthorName: "Test Author",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
				ImagePath:  "test/image/path.jpg",
				Categories: []*domain.Category{
					{Name: "Category1"},
					{Name: "Category2"},
				},
			},
			expectedID:    0,
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			repo := post.NewPostStorage(db)

			mock.ExpectQuery("INSERT INTO posts").
				WithArgs(tc.post.Title, tc.post.Content, tc.post.AuthorID, tc.post.AuthorName, tc.post.CreatedAt, tc.post.UpdatedAt, tc.post.ImagePath).
				WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(tc.expectedID, tc.post.CreatedAt, tc.post.UpdatedAt))

			for _, category := range tc.post.Categories {
				mock.ExpectExec("INSERT INTO PostCategories").
					WithArgs(tc.expectedID, category.Name).
					WillReturnResult(sqlmock.NewResult(0, 1))
			}

			if tc.expectedError == nil {
				mock.ExpectExec("INSERT INTO images").
					WithArgs(tc.expectedID, tc.post.ImagePath).
					WillReturnResult(sqlmock.NewResult(0, 1))
			}

			id, err := repo.CreatePostWithImage(tc.post)

			assert.Equal(t, tc.expectedID, id)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetAllPosts(t *testing.T) {
	testCases := []struct {
		name          string
		offset        int
		limit         int
		expectedPosts []*domain.Post
		expectedError error
	}{
		{
			name:   "Success",
			offset: 0,
			limit:  10,
			expectedPosts: []*domain.Post{
				{
					ID:         1,
					Title:      "Post 1",
					Content:    "Content 1",
					AuthorID:   1,
					AuthorName: "Author 1",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
				{
					ID:         2,
					Title:      "Post 2",
					Content:    "Content 2",
					AuthorID:   2,
					AuthorName: "Author 2",
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				},
			},
			expectedError: nil,
		},
		{
			name:          "No Posts",
			offset:        0,
			limit:         10,
			expectedPosts: nil,
			expectedError: nil,
		},
		{
			name:          "Error",
			offset:        0,
			limit:         10,
			expectedPosts: nil,
			expectedError: errors.New("database error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			repo := post.NewPostStorage(db)

			rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "authorname", "created_at", "updated_at"})
			for _, post := range tc.expectedPosts {
				rows.AddRow(post.ID, post.Title, post.Content, post.AuthorID, post.AuthorName, post.CreatedAt, post.UpdatedAt)
			}
			if tc.expectedError != nil {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM posts ORDER BY id DESC LIMIT $1 OFFSET $2`)).WithArgs(tc.limit, tc.offset).WillReturnError(tc.expectedError)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM posts ORDER BY id DESC LIMIT $1 OFFSET $2`)).WithArgs(tc.limit, tc.offset).WillReturnRows(rows)
			}

			for _, post := range tc.expectedPosts {
				rows := sqlmock.NewRows([]string{"category_name"}).
					AddRow("Category1").
					AddRow("Category2")
				mock.ExpectQuery("^SELECT c.category_name FROM categories c JOIN PostCategories pc ON c.category_name = pc.category_name WHERE pc.post_id = \\$1$").
					WithArgs(post.ID).
					WillReturnRows(rows)

				rows = sqlmock.NewRows([]string{"image_path"}).
					AddRow("image_path1").
					AddRow("image_path2")
				mock.ExpectQuery("^SELECT image_path FROM images WHERE post_id = \\$1$").
					WithArgs(post.ID).
					WillReturnRows(rows)
			}

			posts, err := repo.GetAllPosts(tc.offset, tc.limit)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Nil(t, posts)
				assert.True(t, errors.Is(err, tc.expectedError), "got unexpected error: %v", err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tc.expectedPosts), len(posts))
				for i, expectedPost := range tc.expectedPosts {
					assert.Equal(t, expectedPost.ID, posts[i].ID)
					assert.Equal(t, expectedPost.Title, posts[i].Title)
					assert.Equal(t, expectedPost.Content, posts[i].Content)
					assert.Equal(t, expectedPost.AuthorID, posts[i].AuthorID)
					assert.Equal(t, expectedPost.AuthorName, posts[i].AuthorName)
					assert.Equal(t, expectedPost.CreatedAt, posts[i].CreatedAt)
					assert.Equal(t, expectedPost.UpdatedAt, posts[i].UpdatedAt)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
