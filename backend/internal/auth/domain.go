package auth

import (
	"errors"
	"time"
)

var (
	ErrWrongIdentificator  = errors.New("wrong identificator given")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserAlreadyExist    = errors.New("user with this username already exists")
	ErrPasswordMismatch    = errors.New("password mismatch")
	ErrTooLongPassword     = errors.New("too long password")
	ErrEmptyUsername       = errors.New("empty username given")
	ErrTooShortPassword    = errors.New("too short password")
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrInvalidTokenSub     = errors.New("invalid token subject")
	ErrInvalidBalance      = errors.New("balance amount must be positive")
	ErrInsufficientBalance = errors.New("insufficient balance")
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	Balance      int64
	AccountType  string
}

type jwtToken struct {
	AccessToken string
	ExpiresAt   time.Time
}
