package auth

import (
	"market/internal/platform"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthModule struct {
	handler      *Handler
	tokenManager *platform.TokenManager
}

func NewAuthModule(pool *pgxpool.Pool) *AuthModule {
	repo := NewRepository(pool)
	tkManager := platform.NewTokenManager([]byte("secret"), time.Minute*3)
	serv := NewService(repo, tkManager)
	hand := NewHandler(serv)

	return &AuthModule{
		handler:      hand,
		tokenManager: tkManager,
	}
}

func (a *AuthModule) RequireAuth(next http.Handler) http.Handler {
	return platform.RequireAuth(a.tokenManager)(next)
}

func (a *AuthModule) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", a.handler.Register)
	mux.HandleFunc("POST /auth/login", a.handler.Login)
	mux.Handle("GET /auth/me", a.RequireAuth(http.HandlerFunc(a.handler.Me)))
}
