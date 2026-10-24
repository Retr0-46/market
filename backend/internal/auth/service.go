package auth

import (
	"context"
	"market/internal/platform"
	"strings"
)

type Service struct {
	repo         UserRepository
	tokenManager platform.TokenManagerInterface
}

// AuthService is the application boundary used by HTTP handlers.
type AuthService interface {
	Register(ctx context.Context, username, password, accountType string) (User, error)
	Login(ctx context.Context, username, password string) (LoginResult, error)
	GetUserInfo(ctx context.Context, userID int64) (User, error)
}

var _ AuthService = (*Service)(nil)

func NewService(repo UserRepository, tkManager platform.TokenManagerInterface) *Service {
	return &Service{repo: repo, tokenManager: tkManager}
}

func (s *Service) Register(ctx context.Context, username, password, accountType string) (User, error) {

	name := strings.TrimSpace(username)

	if len(name) < 1 {
		return User{}, ErrEmptyUsername
	}

	if len([]rune(password)) < 4 {
		return User{}, ErrTooShortPassword
	}

	if len([]byte(password)) > 72 {
		return User{}, ErrTooLongPassword
	}

	hashPassword, err := HashPassword(password)
	if err != nil {
		return User{}, err
	}

	user, err := s.repo.CreateUser(ctx, name, hashPassword, accountType)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, username, password string) (LoginResult, error) {
	user, err := s.repo.GetUser(ctx, username)

	if err != nil {
		return LoginResult{}, err
	}

	if ComparePassword(user.PasswordHash, password) != nil {
		return LoginResult{}, ErrPasswordMismatch
	}

	token, err := s.tokenManager.IssueAccessToken(user.ID)

	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{AccessToken: token.Value, ExpiresAt: token.ExpiresAt}, nil
}

func (s *Service) GetUserInfo(ctx context.Context, userID int64) (User, error) {
	user, err := s.repo.GetUser(ctx, userID)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
