package session

import (
	"forum/internal/models"
	"time"

	"github.com/gofrs/uuid"
)

type SessionService struct {
	repo models.SessionRepo
}

func NewSessionService(repo models.SessionRepo) *SessionService {
	return &SessionService{repo}
}

func (s *SessionService) CreateSession(userId int) (*models.Session, error) {
	oldSession, _ := s.repo.GetSessionByUserID(userId)
	if oldSession != nil {
		err := s.repo.DeleteSessionByUUID(oldSession.UUID)
		if err != nil {
			return nil, err
		}
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	session := &models.Session{
		User_id:  userId,
		UUID:     uuid.String(),
		ExpireAt: time.Now().Add(time.Hour),
	}

	err = s.repo.CreateSession(session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *SessionService) DeleteSessionByUUID(uuid string) error {
	_, err := s.repo.GetSessionByUUID(uuid)
	if err != nil {
		return err
	}
	return s.repo.DeleteSessionByUUID(uuid)
}

func (u *SessionService) GetUserIdBySession(session *models.Session) (int, error) {
	user_id, err := u.repo.GetUserIdBySession(session)
	if err != nil {
		return 0, err
	}
	return user_id, nil
}

func (s *SessionService) GetSessionByUUID(uuid string) (*models.Session, error) {
	session, err := s.repo.GetSessionByUUID(uuid)

	switch err {
	case nil:
		if session.ExpireAt.Before(time.Now()) {
			return nil, models.ErrSessionExpired
		}
		return session, nil
	case models.ErrSqlNoRows:
		return nil, models.ErrSqlNoRows
	default:
		return nil, err
	}
}
