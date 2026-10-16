package cart

import (
	"market/internal/platform"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type CartModule struct {
	handler      *Handler
	tokenManager *platform.TokenManager
}

func NewCartModule(rdb *redis.Client, secret []byte) *CartModule {
	repo := NewRepository(rdb)
	serv := NewService(repo)
	return &CartModule{
		handler:      NewHandler(serv),
		tokenManager: platform.NewTokenManager(secret, time.Minute*3),
	}
}

func (m *CartModule) RequireAuth(next http.Handler) http.Handler {
	return platform.RequireAuth(m.tokenManager)(next)
}

func (m *CartModule) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /cart", m.RequireAuth(http.HandlerFunc(m.handler.GetCart)))
	mux.Handle("POST /cart/items", m.RequireAuth(http.HandlerFunc(m.handler.AddItem)))
	mux.Handle("PUT /cart/items/{product_id}", m.RequireAuth(http.HandlerFunc(m.handler.SetItem)))
	mux.Handle("DELETE /cart/items/{product_id}", m.RequireAuth(http.HandlerFunc(m.handler.DeleteItem)))
	mux.Handle("DELETE /cart", m.RequireAuth(http.HandlerFunc(m.handler.ClearCart)))
}
