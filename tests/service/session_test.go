package repository

import (
	"database/sql"
	"errors"
	"forum/internal/domain"
	"forum/internal/repository/session"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestGetSessionByUserID(t *testing.T) {
	testCases := []struct {
		name              string
		userID            int
		expectedSession   *domain.Session
		expectedErrorType error
	}{
		{
			name:              "Valid user ID",
			userID:            1,
			expectedSession:   &domain.Session{UUID: "123", User_id: 1, ExpireAt: time.Now()},
			expectedErrorType: nil,
		},
		{
			name:              "Non-existent user ID",
			userID:            999,
			expectedSession:   nil,
			expectedErrorType: sql.ErrNoRows,
		},
		{
			name:              "Negative user ID",
			userID:            -1,
			expectedSession:   nil,
			expectedErrorType: errors.New("negative user ID"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			repo := session.NewSessionStorage(db)

			rows := sqlmock.NewRows([]string{"uuid", "user_id", "expire_at"}).
				AddRow("123", 1, time.Now())
			if tc.expectedErrorType == nil {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM sessions WHERE user_ID = $1")).
					WithArgs(tc.userID).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM sessions WHERE user_ID = $1")).
					WithArgs(tc.userID).
					WillReturnError(tc.expectedErrorType)
			}

			session, err := repo.GetSessionByUserID(tc.userID)

			assert.Equal(t, tc.expectedErrorType, err)
			assert.Equal(t, tc.expectedSession, session)
		})
	}
}

func TestGetSessionByUUID(t *testing.T) {
	testCases := []struct {
		name              string
		sessionID         string
		expectedSession   *domain.Session
		expectedErrorType error
	}{
		{
			name:              "Valid session ID",
			sessionID:         "123",
			expectedSession:   &domain.Session{UUID: "123", User_id: 1, ExpireAt: time.Now()},
			expectedErrorType: nil,
		},
		{
			name:              "Non-existent session ID",
			sessionID:         "999",
			expectedSession:   nil,
			expectedErrorType: sql.ErrNoRows,
		},
		{
			name:              "Empty session ID",
			sessionID:         "",
			expectedSession:   nil,
			expectedErrorType: errors.New("empty session ID"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			repo := session.NewSessionStorage(db)

			if tc.expectedErrorType == nil {
				rows := sqlmock.NewRows([]string{"uuid", "user_id", "expire_at"}).
					AddRow(tc.expectedSession.UUID, tc.expectedSession.User_id, tc.expectedSession.ExpireAt)

				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM sessions WHERE uuid = $1")).
					WithArgs(tc.sessionID).
					WillReturnRows(rows)
			} else {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM sessions WHERE uuid = $1")).
					WithArgs(tc.sessionID).
					WillReturnError(tc.expectedErrorType)
			}

			session, err := repo.GetSessionByUUID(tc.sessionID)
			assert.Equal(t, tc.expectedErrorType, err)

			assert.Equal(t, tc.expectedSession, session)
		})
	}
}

func TestDeleteSessionByUUID(t *testing.T) {
	testCases := []struct {
		name              string
		sessionID         string
		expectedErrorType error
	}{
		{
			name:              "Valid session ID",
			sessionID:         "123",
			expectedErrorType: nil,
		},
		{
			name:              "Non-existent session ID",
			sessionID:         "999",
			expectedErrorType: sql.ErrNoRows,
		},
		{
			name:              "Empty session ID",
			sessionID:         "",
			expectedErrorType: errors.New("empty session ID"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			repo := session.NewSessionStorage(db)
			if tc.expectedErrorType == nil {
				mock.ExpectExec("DELETE FROM sessions WHERE uuid = ?").
					WithArgs(tc.sessionID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("DELETE FROM sessions WHERE uuid = ?").
					WithArgs(tc.sessionID).
					WillReturnError(tc.expectedErrorType)
			}

			err = repo.DeleteSessionByUUID(tc.sessionID)

			assert.IsType(t, tc.expectedErrorType, err)
		})
	}
}
