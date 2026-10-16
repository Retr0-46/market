package product

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ProductModule struct {
	handler *Handler
}

func NewProductModule(pool *pgxpool.Pool, rdb *redis.Client) *ProductModule {
	repo := NewRepository(pool, rdb)
	repo.cacheTTL = 5 * time.Minute
	serv := NewService(repo)

	return &ProductModule{handler: NewHandler(serv)}
}

func (m *ProductModule) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /products/list", m.handler.GetHomePage)
	mux.HandleFunc("GET /products", m.handler.GetProduct)
}
