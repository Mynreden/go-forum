package repository

import (
	"database/sql"
	"forum/internal/domain"
	"forum/internal/repository/user"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name           string
		user           *domain.User
		expectedError  error
		expectedErrMsg string
	}{
		{
			name: "Valid user",
			user: &domain.User{
				Username:  "testuser",
				HashedPW:  "hashedpassword",
				Email:     "test@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: nil,
		},
		{
			name: "Duplicate email",
			user: &domain.User{
				Username:  "testuser2",
				HashedPW:  "hashedpassword2",
				Email:     "duplicate@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError:  domain.ErrDuplicateEmail,
			expectedErrMsg: "UNIQUE constraint failed: users.email",
		},
		{
			name: "Duplicate username",
			user: &domain.User{
				Username:  "duplicateuser",
				HashedPW:  "hashedpassword3",
				Email:     "test3@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError:  domain.ErrDuplicateUsername,
			expectedErrMsg: "UNIQUE constraint failed: users.username",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %s", err)
			}
			defer db.Close()
			if tc.expectedError == nil {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(tc.user.Username, tc.user.HashedPW, tc.user.Email, tc.user.CreatedAt, tc.user.UpdatedAt).
					WillReturnResult(sqlmock.NewResult(1, 1))
			} else {
				mock.ExpectExec("INSERT INTO users").
					WithArgs(tc.user.Username, tc.user.HashedPW, tc.user.Email, tc.user.CreatedAt, tc.user.UpdatedAt).
					WillReturnError(tc.expectedError)
			}

			userRepo := user.NewUserStorage(db)

			err = userRepo.CreateUser(tc.user)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	testCases := []struct {
		name           string
		username       string
		mockRows       *sqlmock.Rows
		expectedUser   *domain.User
		expectedError  error
		expectedErrMsg string
	}{
		{
			name:     "Valid username",
			username: "testuser",
			mockRows: sqlmock.NewRows([]string{"id", "username", "hashed_pw", "email"}).AddRow(1, "testuser", "hashedpassword", "test@example.com"),
			expectedUser: &domain.User{
				ID:       1,
				Username: "testuser",
				HashedPW: "hashedpassword",
				Email:    "test@example.com",
			},
			expectedError: nil,
		},
		{
			name:           "User not found",
			username:       "nonexistentuser",
			mockRows:       sqlmock.NewRows([]string{}),
			expectedUser:   nil,
			expectedError:  sql.ErrNoRows,
			expectedErrMsg: "no rows in result set",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %s", err)
			}
			defer db.Close()

			mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE username = $1")).
				WithArgs(tc.username).
				WillReturnRows(tc.mockRows)

			userRepo := user.NewUserStorage(db)

			user, err := userRepo.GetUserByUsername(tc.username)

			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedUser, user)
			if err != nil {
				assert.Contains(t, err.Error(), tc.expectedErrMsg)
			}
		})
	}
}
