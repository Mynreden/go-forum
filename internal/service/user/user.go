package user

import (
	"forum/internal/domain"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo domain.UserRepo
}

func NewUserService(repo domain.UserRepo) *UserService {
	return &UserService{repo}
}

func (u *UserService) CreateUser(userDTO *domain.CreateUserDTO) error {
	hashedPW, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Username:  userDTO.Username,
		Email:     userDTO.Email,
		HashedPW:  string(hashedPW),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err = u.repo.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserService) LoginUser(userDTO *domain.LoginUserDTO) (int, error) {
	user, err := u.repo.GetUserByEmail(userDTO.Email)
	if err != nil {
		switch err {
		case domain.ErrSqlNoRows:
			return 0, domain.ErrInvalidCredentials
		default:
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPW), []byte(userDTO.Password))
	if err != nil {
		return 0, domain.ErrInvalidCredentials
	}

	return user.ID, nil
}

func (u *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return u.repo.GetUserByEmail(email)
}

// func (u *UserService) GetUserBySession(session *domain.Session) (user *domain.User, err error) {
// 	user_id, err := u.repo.GetUserIDySession(session)
// 	if err != nil {
// 		return nil, err
// 	}
// 	user, err := u.repo.GetUserByID(user_id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return user, nil
// }

// if needed=================
func (u *UserService) UpdateUser(user *domain.User) error {
	return nil
}

// func (s *SessionService) GetUserBySession(session *domain.Session) (user *domain.User, err error) {
// 	user, err = s.repo.GetUserByID(user.ID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return user, nil
// }

func (u *UserService) GetUserByID(id int) (user *domain.User, err error) {
	return u.repo.GetUserByID(id)
}
