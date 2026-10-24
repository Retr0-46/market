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

// LoginResult contains the access token returned after a successful login.
// It is exported so callers can implement AuthService in tests without
// depending on an unexported return type.
type LoginResult struct {
	AccessToken string
	ExpiresAt   time.Time
}

// jwtToken is kept as an alias for code that used the old internal name.
type jwtToken = LoginResult
