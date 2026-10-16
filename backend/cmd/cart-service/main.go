package main

import (
	"log"
	"market/internal/cart"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	defer rdb.Close()

	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		secret = "secret"
	}

	mux := http.NewServeMux()
	cart.NewCartModule(rdb, []byte(secret)).RegisterRoutes(mux)

	log.Println("Started cart-service on :8082")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Fatal(err)
	}
}
