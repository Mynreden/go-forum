package session

import (
	"forum/internal/domain"
	"time"

	"github.com/gofrs/uuid"
)

type SessionService struct {
	repo domain.SessionRepo
}

func NewSessionService(repo domain.SessionRepo) *SessionService {
	return &SessionService{repo}
}

func (s *SessionService) CreateSession(userId int) (*domain.Session, error) {
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
	session := &domain.Session{
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

func (u *SessionService) GetUserIdBySession(session *domain.Session) (int, error) {
	user_id, err := u.repo.GetUserIdBySession(session)
	if err != nil {
		return 0, err
	}
	return user_id, nil
}

func (s *SessionService) GetSessionByUUID(uuid string) (*domain.Session, error) {
	session, err := s.repo.GetSessionByUUID(uuid)

	switch err {
	case nil:
		if session.ExpireAt.Before(time.Now()) {
			return nil, domain.ErrSessionExpired
		}
		return session, nil
	case domain.ErrSqlNoRows:
		return nil, domain.ErrSqlNoRows
	default:
		return nil, err
	}
}
