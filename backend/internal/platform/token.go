package platform

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrInvalidTokenSub    = errors.New("invalid token subject")
)

type TokenManagerInterface interface {
	IssueAccessToken(userID int64) (AccessToken, error)
	VerifyAccessToken(raw string) (int64, error)
}

type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

type AccessToken struct {
	Value     string
	ExpiresAt time.Time
}

func NewTokenManager(secret []byte, ttl time.Duration) *TokenManager {
	return &TokenManager{secret: secret, ttl: ttl}
}

func (tm *TokenManager) IssueAccessToken(userID int64) (AccessToken, error) {
	now := time.Now()
	expiresAt := now.Add(tm.ttl)

	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		Issuer:    "market",
		Audience:  jwt.ClaimStrings{"market"},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(tm.secret)

	if err != nil {
		return AccessToken{}, err
	}

	return AccessToken{signedToken, expiresAt}, nil
}

func (tm *TokenManager) VerifyAccessToken(raw string) (int64, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (any, error) { return tm.secret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer("market"),
		jwt.WithAudience("market"),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, ErrInvalidAccessToken
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 {
		return 0, ErrInvalidTokenSub
	}

	return userID, nil
}
