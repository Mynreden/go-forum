package repository

import (
	"forum/internal/domain"
	"forum/internal/repository/postReaction"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestGetPostReactionsByPostID(t *testing.T) {
	testCases := []struct {
		name              string
		postID            int
		expectedReactions []*domain.PostReaction
		expectedError     error
	}{
		{
			name:   "Valid post ID with reactions",
			postID: 1,
			expectedReactions: []*domain.PostReaction{
				{UserID: 0, PostID: 0, Status: true},
				{UserID: 0, PostID: 0, Status: false},
			},
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

			repo := postReaction.NewPostReactionStorage(db)

			rows := sqlmock.NewRows([]string{"reaction"}).
				AddRow(1).
				AddRow(0)

			mock.ExpectQuery("SELECT reaction FROM postsReactions WHERE post_id = ?").
				WithArgs(tc.postID).
				WillReturnRows(rows)

			reactions, err := repo.GetPostReactionsByPostID(tc.postID)
			for i, react := range reactions {
				assert.Equal(t, tc.expectedReactions[i], react)

			}
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetPostsReactionsByUserID(t *testing.T) {
	testCases := []struct {
		name              string
		userID            int
		expectedReactions []*domain.PostReaction
		expectedError     error
	}{
		{
			name:   "Valid user ID with reactions",
			userID: 1,
			expectedReactions: []*domain.PostReaction{
				{UserID: 1, PostID: 1, Status: true},
				{UserID: 1, PostID: 2, Status: false},
			},
			expectedError: nil,
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			repo := postReaction.NewPostReactionStorage(db)

			rows := sqlmock.NewRows([]string{"post_id", "reaction"})
			for _, react := range tc.expectedReactions {
				rows.AddRow(react.PostID, react.Status)
			}

			mock.ExpectQuery("SELECT post_id, reaction FROM postsReactions WHERE user_id = ?").
				WithArgs(tc.userID).
				WillReturnRows(rows)

			reactions, err := repo.GetPostsReactionsByUserID(tc.userID)
			for i, react := range reactions {
				assert.Equal(t, tc.expectedReactions[i], react)

			}
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetAllPostReactions(t *testing.T) {
	testCases := []struct {
		name              string
		expectedReactions []*domain.PostReaction
		expectedError     error
	}{
		{
			name: "Valid reactions",
			expectedReactions: []*domain.PostReaction{
				{UserID: 1, PostID: 1, Status: true},
				{UserID: 2, PostID: 1, Status: false},
				{UserID: 1, PostID: 2, Status: true},
			},
			expectedError: nil,
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			repo := postReaction.NewPostReactionStorage(db)

			rows := sqlmock.NewRows([]string{"user_id", "post_id", "reaction"}).
				AddRow(1, 1, 1).
				AddRow(2, 1, 0).
				AddRow(1, 2, 1)

			mock.ExpectQuery("SELECT user_id, post_id, reaction FROM postsReactions").
				WillReturnRows(rows)

			reactions, err := repo.GetAllPostReactions()

			assert.Equal(t, tc.expectedReactions, reactions)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetReactionByUserIDAndPostID(t *testing.T) {
	testCases := []struct {
		name             string
		userID           int
		postID           int
		expectedReaction *domain.PostReaction
		expectedError    error
	}{
		{
			name:   "Valid user and post ID with reaction",
			userID: 1,
			postID: 1,
			expectedReaction: &domain.PostReaction{
				UserID: 1,
				PostID: 1,
				Status: true,
			},
			expectedError: nil,
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			repo := postReaction.NewPostReactionStorage(db)

			rows := sqlmock.NewRows([]string{"id", "reaction"}).
				AddRow(tc.expectedReaction.ID, tc.expectedReaction.Status)

			mock.ExpectQuery(regexp.QuoteMeta("SELECT id, reaction FROM postsReactions WHERE user_id = ? AND post_id = ?")).
				WithArgs(tc.userID, tc.postID).
				WillReturnRows(rows)

			reaction, err := repo.GetReactionByUserIDAndPostID(tc.userID, tc.postID)

			assert.Equal(t, tc.expectedReaction, reaction)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestDeletePostReactionByID(t *testing.T) {
	testCases := []struct {
		name          string
		reactionID    int
		expectedError error
	}{
		{
			name:          "Valid reaction ID",
			reactionID:    1,
			expectedError: nil,
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			repo := postReaction.NewPostReactionStorage(db)

			mock.ExpectExec("DELETE FROM postsReactions WHERE id = ?").
				WithArgs(tc.reactionID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			err = repo.DeletePostReactionByID(tc.reactionID)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}
