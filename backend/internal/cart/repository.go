package cart

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	rdb *redis.Client
}

// CartRepository is the persistence boundary used by the cart service.
type CartRepository interface {
	GetCart(ctx context.Context, userID int64) ([]Item, error)
	AddItem(ctx context.Context, userID, productID, quantity int64) (Item, error)
	SetItem(ctx context.Context, userID, productID, quantity int64) (Item, error)
	DeleteItem(ctx context.Context, userID, productID int64) error
	ClearCart(ctx context.Context, userID int64) error
}

var _ CartRepository = (*Repository)(nil)

func NewRepository(rdb *redis.Client) *Repository {
	return &Repository{rdb: rdb}
}

func (r *Repository) GetCart(ctx context.Context, userID int64) ([]Item, error) {
	values, err := r.rdb.HGetAll(ctx, r.key(userID)).Result()
	if err != nil {
		return nil, err
	}

	items := make([]Item, 0, len(values))
	for field, value := range values {
		productID, err := strconv.ParseInt(field, 10, 64)
		if err != nil {
			return nil, err
		}
		quantity, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return nil, err
		}
		items = append(items, Item{ProductID: productID, Quantity: quantity})
	}
	return items, nil
}

func (r *Repository) AddItem(ctx context.Context, userID, productID, quantity int64) (Item, error) {
	total, err := r.rdb.HIncrBy(ctx, r.key(userID), strconv.FormatInt(productID, 10), quantity).Result()
	if err != nil {
		return Item{}, err
	}
	return Item{ProductID: productID, Quantity: total}, nil
}

func (r *Repository) SetItem(ctx context.Context, userID, productID, quantity int64) (Item, error) {
	err := r.rdb.HSet(ctx, r.key(userID), strconv.FormatInt(productID, 10), quantity).Err()
	if err != nil {
		return Item{}, err
	}
	return Item{ProductID: productID, Quantity: quantity}, nil
}

func (r *Repository) DeleteItem(ctx context.Context, userID, productID int64) error {
	return r.rdb.HDel(ctx, r.key(userID), strconv.FormatInt(productID, 10)).Err()
}

func (r *Repository) ClearCart(ctx context.Context, userID int64) error {
	return r.rdb.Del(ctx, r.key(userID)).Err()
}

func (r *Repository) key(userID int64) string {
	return "cart:" + strconv.FormatInt(userID, 10)
}
