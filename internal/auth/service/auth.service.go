package service

import (
	"github.com/ezep02/microservicios/internal/auth/repository"
	"github.com/ezep02/microservicios/internal/auth/types"
)

type AuthService struct {
	repo *repository.AuthRepository
}

func NewAuthService(authRepo *repository.AuthRepository) *AuthService {
	return &AuthService{
		repo: authRepo,
	}
}

func (Auth_s *AuthService) UserSignUp(user *types.User) (*types.UserResponse, error) {
	return Auth_s.repo.UserSignUp(user)
} //Register user service init

func (Auth_s *AuthService) UserSignIn(user *types.UserRequest) (*types.UserSignInResponse, error) {
	return Auth_s.repo.UserSignIn(user)
} //Login user service init

// Session
func (Auth_s *AuthService) CreateSession(s *types.Session) (*types.Session, error) {
	return Auth_s.repo.CreateSession(s)
} //Register user service init

func (Auth_s *AuthService) GetSession(id string) (*types.Session, error) {
	return Auth_s.repo.GetSession(id)
} //

func (Auth_s *AuthService) RevokeSession(id string) error {
	return Auth_s.repo.RevokeSession(id)
} //Register user service init

func (Auth_s *AuthService) DeleteSession(id string) error {
	return Auth_s.repo.DeleteSession(id)
} //Delete session
