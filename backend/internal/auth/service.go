package auth

import (
	"context"
	"market/internal/platform"
	"strings"
)

type Service struct {
	repo         *Repository
	tokenManager *platform.TokenManager
}

func NewService(repo *Repository, tkManager *platform.TokenManager) *Service {
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

func (s *Service) Login(ctx context.Context, username, password string) (jwtToken, error) {
	user, err := s.repo.GetUser(ctx, username)

	if err != nil {
		return jwtToken{}, err
	}

	if ComparePassword(user.PasswordHash, password) != nil {
		return jwtToken{}, ErrPasswordMismatch
	}

	token, err := s.tokenManager.IssueAccessToken(user.ID)

	if err != nil {
		return jwtToken{}, err
	}

	return jwtToken{token.Value, token.ExpiresAt}, nil
}

func (s *Service) GetUserInfo(ctx context.Context, userID int64) (User, error) {
	user, err := s.repo.GetUser(ctx, userID)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
