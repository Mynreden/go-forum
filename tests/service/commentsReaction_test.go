package repository

import (
	"database/sql"
	"errors"
	"forum/internal/domain"
	"forum/internal/repository/commentsReaction"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCommentsReactionsStorage_CreateCommentsReactions(t *testing.T) {
	testCases := []struct {
		name          string
		reaction      *domain.CommentReactionDTO
		expectedError error
	}{
		{
			name: "Successful creation",
			reaction: &domain.CommentReactionDTO{
				ID:        1,
				UserID:    1,
				CommentID: 1,
				Status:    true,
			},
			expectedError: nil,
		},
		{
			name: "Empty status",
			reaction: &domain.CommentReactionDTO{
				ID:        1,
				UserID:    1,
				CommentID: 1,
				Status:    true,
			},
			expectedError: errors.New("reaction status is empty"),
		},
		{
			name: "Database error",
			reaction: &domain.CommentReactionDTO{
				ID:        1,
				UserID:    1,
				CommentID: 1,
				Status:    true,
			},
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

			if tc.expectedError == nil {
				mock.ExpectExec("INSERT INTO commentsReactions").
					WithArgs(tc.reaction.UserID, tc.reaction.CommentID, tc.reaction.Status).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("INSERT INTO commentsReactions").
					WithArgs(tc.reaction.UserID, tc.reaction.CommentID, tc.reaction.Status).
					WillReturnError(tc.expectedError)
			}

			repo := commentsReaction.NewCommentsReactionsStorage(db)
			err = repo.CreateCommentsReactions(tc.reaction)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCommentsReactionsStorage_DeleteCommentsReactions(t *testing.T) {
	testCases := []struct {
		name          string
		commentID     int
		expectedError error
	}{
		{
			name:          "Successful deletion",
			commentID:     1,
			expectedError: nil,
		},
		{
			name:          "Database error",
			commentID:     1,
			expectedError: errors.New("database error"),
		},
		{
			name:          "Nonexistent comment ID",
			commentID:     999,
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			if tc.expectedError == nil {
				mock.ExpectExec("DELETE FROM commentsReactions").
					WithArgs(tc.commentID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			} else {
				mock.ExpectExec("DELETE FROM commentsReactions").
					WithArgs(tc.commentID).
					WillReturnError(tc.expectedError)
			}

			repo := commentsReaction.NewCommentsReactionsStorage(db)
			err = repo.DeleteCommentsReactions(tc.commentID)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestCommentsReactionsStorage_GetReactionByUserIDAndCommentID(t *testing.T) {
	testCases := []struct {
		name              string
		userID, commentID int
		expectedRow       *sqlmock.Rows
		expectedError     error
	}{
		{
			name:          "Successful retrieval",
			userID:        1,
			commentID:     1,
			expectedRow:   sqlmock.NewRows([]string{"id", "reaction"}).AddRow(1, true),
			expectedError: nil,
		},
		{
			name:          "No reaction found",
			userID:        1,
			commentID:     1,
			expectedRow:   sqlmock.NewRows([]string{"id", "reaction"}),
			expectedError: sql.ErrNoRows,
		},
		{
			name:          "Database error",
			userID:        1,
			commentID:     1,
			expectedRow:   nil,
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

			if tc.expectedError == nil {
				mock.ExpectQuery("SELECT id, reaction FROM commentsReactions").
					WithArgs(tc.userID, tc.commentID).
					WillReturnRows(tc.expectedRow)
			} else {
				mock.ExpectQuery("SELECT id, reaction FROM commentsReactions").
					WithArgs(tc.userID, tc.commentID).
					WillReturnError(tc.expectedError)
			}

			repo := commentsReaction.NewCommentsReactionsStorage(db)
			reaction, err := repo.GetReactionByUserIDAndCommentID(tc.userID, tc.commentID)

			assert.Equal(t, tc.expectedError, err)
			if err == nil {
				assert.NotNil(t, reaction)
			} else {
				assert.Nil(t, reaction)
			}
		})
	}
}
