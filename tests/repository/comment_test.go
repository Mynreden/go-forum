package repository

import (
	"database/sql"
	"errors"
	"forum/internal/domain"
	"forum/internal/repository/comment"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCommentStorage_CreateComment(t *testing.T) {
	testCases := []struct {
		name           string
		comment        *domain.Comment
		expectedError  error
		expectedExecID int64
	}{
		{
			name: "Successful creation",
			comment: &domain.Comment{
				Content:    "Test comment",
				PostID:     1,
				AuthorID:   1,
				AuthorName: "testuser",
				CreatedAt:  time.Now(),
			},
			expectedError:  nil,
			expectedExecID: 1,
		},
		{
			name: "Empty content",
			comment: &domain.Comment{
				Content:    "",
				PostID:     1,
				AuthorID:   1,
				AuthorName: "testuser",
				CreatedAt:  time.Now(),
			},
			expectedError:  errors.New("comment content is empty"),
			expectedExecID: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			if tc.expectedError == nil {
				mock.ExpectExec("INSERT INTO comments").
					WithArgs(tc.comment.Content, tc.comment.PostID, tc.comment.AuthorID, tc.comment.AuthorName, tc.comment.CreatedAt).
					WillReturnResult(sqlmock.NewResult(tc.expectedExecID, 1))
			} else {
				mock.ExpectExec("INSERT INTO comments").
					WithArgs(tc.comment.Content, tc.comment.PostID, tc.comment.AuthorID, tc.comment.AuthorName, tc.comment.CreatedAt).
					WillReturnError(tc.expectedError)
			}

			storage := comment.NewCommentStorage(db)
			err = storage.CreateComment(tc.comment)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCommentStorage_GetAllByPostID(t *testing.T) {
	testCases := []struct {
		name           string
		postID         int
		expectedRows   *sqlmock.Rows
		expectedError  error
		expectedLength int
	}{
		{
			name:   "Successful retrieval",
			postID: 1,
			expectedRows: sqlmock.NewRows([]string{"id", "comment", "post_id", "user_id", "userName", "created_at"}).
				AddRow(1, "Test comment 1", 1, 1, "testuser", time.Now()).
				AddRow(2, "Test comment 2", 1, 2, "anotheruser", time.Now()),
			expectedError:  nil,
			expectedLength: 2,
		},
		{
			name:           "No comments found",
			postID:         2,
			expectedRows:   sqlmock.NewRows([]string{"id", "comment", "post_id", "user_id", "userName", "created_at"}),
			expectedError:  nil,
			expectedLength: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)

			mock.ExpectQuery("SELECT id, comment, post_id, user_id, userName, created_at FROM comments WHERE post_id = ?").
				WithArgs(tc.postID).
				WillReturnRows(tc.expectedRows)

			storage := comment.NewCommentStorage(db)
			comments, err := storage.GetAllByPostID(tc.postID)

			assert.Equal(t, tc.expectedError, err)
			assert.Len(t, comments, tc.expectedLength)
		})
	}
}

func TestCommentStorage_GetCommentByID(t *testing.T) {
	testCases := []struct {
		name           string
		commentID      int
		expectedRow    *sqlmock.Rows
		expectedError  error
		expectedLength int
	}{
		{
			name:      "Successful retrieval",
			commentID: 1,
			expectedRow: sqlmock.NewRows([]string{"id", "comment", "post_id", "user_id", "userName", "created_at"}).
				AddRow(1, "Test comment", 1, 1, "testuser", time.Now()),
			expectedError:  nil,
			expectedLength: 1,
		},
		{
			name:           "Comment not found",
			commentID:      2,
			expectedRow:    sqlmock.NewRows([]string{"id", "comment", "post_id", "user_id", "userName", "created_at"}),
			expectedError:  sql.ErrNoRows,
			expectedLength: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func(db *sql.DB) {
				err := db.Close()
				if err != nil {

				}
			}(db)
			if tc.expectedError == nil {
				mock.ExpectQuery("SELECT id, comment, post_id, user_id, userName, created_at FROM comments WHERE id = ?").
					WithArgs(tc.commentID).
					WillReturnRows(tc.expectedRow)
			} else {
				mock.ExpectQuery("SELECT id, comment, post_id, user_id, userName, created_at FROM comments WHERE id = ?").
					WithArgs(tc.commentID).
					WillReturnError(tc.expectedError)
			}

			storage := comment.NewCommentStorage(db)
			com, err := storage.GetCommentByID(tc.commentID)

			assert.Equal(t, tc.expectedError, err)
			if tc.expectedError == nil {
				assert.NotNil(t, com)
			} else {
				assert.Nil(t, com)
			}
		})
	}
}
