package user

import (
	"database/sql"
	"forum/internal/domain"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{db: db}
}

// registration user ===================================================
func (s *UserStorage) CreateUser(user *domain.User) error {
	_, err := s.db.Exec("INSERT INTO users (username, hashed_pw, email, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)",
		user.Username,
		user.HashedPW,
		user.Email,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		switch err.Error() {
		case "UNIQUE constraint failed: users.email":
			return domain.ErrDuplicateEmail
		case "UNIQUE constraint failed: users.username":
			return domain.ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

// for authentification user=============================================
func (s *UserStorage) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := s.db.QueryRow("SELECT * FROM users WHERE username = $1", username).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPW,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStorage) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := s.db.QueryRow("SELECT * FROM users WHERE email = $1", email).Scan(&user.ID,
		&user.Username,
		&user.HashedPW,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// for update user datas=========================================
func (s *UserStorage) UpdateUser(user *domain.User) error {
	_, err := s.db.Exec("UPDATE users SET username = $1, hashed_pw = $2, email = $3 WHERE id = $4",
		user.Username,
		user.HashedPW,
		user.Email,
		user.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) GetAllUsers() ([]*domain.User, error) {
	return nil, nil
}

// administration.UsersFuncs =====================================
func (s *UserStorage) DeleteUser(user *domain.User) error {
	return nil
}

func (s *UserStorage) GetUserByID(id int) (*domain.User, error) {
	var user domain.User
	err := s.db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&user.ID,
		&user.Username,
		&user.HashedPW,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
