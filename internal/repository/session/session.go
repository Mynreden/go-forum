package session

import (
	"database/sql"
	"forum/internal/models"
)

type SessionStorage struct {
	db *sql.DB
}

func NewSessionStorage(db *sql.DB) *SessionStorage {
	return &SessionStorage{db}
}

func (s *SessionStorage) CreateSession(session *models.Session) error {
	_, err := s.db.Exec("INSERT INTO sessions (uuid, user_id, expire_at) VALUES ($1, $2, $3)", session.UUID, session.User_id, session.ExpireAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionStorage) GetSessionByUserID(userID int) (*models.Session, error) {
	var session models.Session
	err := s.db.QueryRow("SELECT * FROM sessions WHERE user_ID = $1", userID).Scan(&session.UUID, &session.User_id, &session.ExpireAt)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionStorage) GetSessionByUUID(sessionID string) (*models.Session, error) {
	var session models.Session
	err := s.db.QueryRow("SELECT * FROM sessions WHERE uuid = $1", sessionID).Scan(&session.UUID, &session.User_id, &session.ExpireAt)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// func (s *SessionStorage)

func (s *SessionStorage) DeleteSessionByUUID(sessionID string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE uuid = $1", sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (s *SessionStorage) GetUserIdBySession(session *models.Session) (user_id int, err error) {
	err = s.db.QueryRow("SELECT user_id FROM sessions WHERE uuid = $1", session.UUID).Scan(&user_id)
	if err != nil {
		return 0, err
	}

	return user_id, nil
}
