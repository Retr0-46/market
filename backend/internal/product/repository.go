package product

import (
	"context"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	pool     *pgxpool.Pool
	rdb      *redis.Client
	cacheTTL time.Duration
}

// ProductRepository is the persistence boundary used by the product service.
type ProductRepository interface {
	CreateProduct(ctx context.Context, productType string, sellerID int64, name string, price, stockQuantity int64) (Product, error)
	GetProductByID(ctx context.Context, id int64) (Product, error)
	GetListProduct(ctx context.Context, number int64, pTypes ...string) ([]Product, error)
}

var _ ProductRepository = (*Repository)(nil)

func NewRepository(pool *pgxpool.Pool, rdb *redis.Client) *Repository {
	return &Repository{pool: pool, rdb: rdb}
}

func (r *Repository) CreateProduct(ctx context.Context, productType string, sellerID int64, name string, price, stockQuantity int64) (Product, error) {
	var prod Product

	row := r.pool.QueryRow(ctx, `INSERT INTO products (p_type, seller_id, name, price, stock_quantity) VALUES ($1, $2, $3, $4, $5) RETURNING id, p_type, seller_id, name, price, stock_quantity`,
		productType, sellerID, name, price, stockQuantity)
	err := row.Scan(&prod.ID, &prod.Type, &prod.SellerID, &prod.Name, &prod.Price, &prod.StockQuantity)

	if err != nil {
		return Product{}, err
	}

	return prod, nil
}

func (r *Repository) GetProductByID(ctx context.Context, id int64) (Product, error) {
	var prod Product

	cmd := r.rdb.HGetAll(ctx, r.key(id))

	if len(cmd.Val()) != 0 {
		if err := cmd.Scan(&prod); err != nil {
			return prod, nil
		}
	}

	row := r.pool.QueryRow(ctx, `SELECT id, p_type, seller_id, name, price, stock_quantity FROM products WHERE id = $1`, id)
	err := row.Scan(&prod.ID, &prod.Type, &prod.SellerID, &prod.Name, &prod.Price, &prod.StockQuantity)

	if err != nil {
		return Product{}, err
	}

	_ = r.rdb.HSet(ctx, r.key(prod.ID), prod).Err()
	_ = r.rdb.Expire(ctx, r.key(prod.ID), r.cacheTTL).Err()

	return prod, nil
}

func (r *Repository) key(id int64) string {
	return "product:" + strconv.FormatInt(id, 10)
}

// GetListProduct возвращает список продуктов заданных типов с лимитом number.
// Если pTypes пуст — фильтр по типу не применяется.
//
// Схема:
//  1. Из БД берём только id (index-only scan, дёшево).
//  2. Пайплайном читаем продукты из Redis по id.
//  3. Промахи одним батчем добираем из БД и прогреваем кэш.
//  4. Собираем результат в исходном порядке.
func (r *Repository) GetListProduct(ctx context.Context, number int64, pTypes ...string) ([]Product, error) {
	if number <= 0 {
		return nil, ErrNegativeNumber
	}

	// 1. Получаем только id
	var (
		rows pgx.Rows
		err  error
	)
	if len(pTypes) == 0 {
		rows, err = r.pool.Query(ctx,
			`SELECT id FROM products LIMIT $1`, number)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT id FROM products WHERE p_type = ANY($1) LIMIT $2`, pTypes, number)
	}
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, number)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, err
		}
		ids = append(ids, id)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []Product{}, nil
	}

	// 2. Пайплайном читаем из Redis
	pipe := r.rdb.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(ids))
	for i, id := range ids {
		cmds[i] = pipe.HGetAll(ctx, r.key(id))
	}
	// Общая ошибка Exec интересует мало — проверяем каждую команду отдельно,
	// чтобы корректно отработать частичный промах.
	_, _ = pipe.Exec(ctx)

	found := make(map[int64]Product, len(ids))
	var missing []int64
	for i, id := range ids {
		cmd := cmds[i]
		if cmd.Err() != nil || len(cmd.Val()) == 0 {
			missing = append(missing, id)
			continue
		}
		var p Product
		if err := cmd.Scan(&p); err != nil {
			missing = append(missing, id)
			continue
		}
		found[id] = p
	}

	// 3. Промахи — одним запросом в БД
	if len(missing) > 0 {
		rows, err := r.pool.Query(ctx,
			`SELECT id, p_type, seller_id, name, price, stock_quantity
			 FROM products WHERE id = ANY($1)`, missing)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		pipe := r.rdb.Pipeline()
		for rows.Next() {
			var p Product
			if err := rows.Scan(&p.ID, &p.Type, &p.SellerID, &p.Name, &p.Price, &p.StockQuantity); err != nil {
				return nil, err
			}
			found[p.ID] = p
			pipe.HSet(ctx, r.key(p.ID), p)
			pipe.Expire(ctx, r.key(p.ID), r.cacheTTL)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
		_, _ = pipe.Exec(ctx)
	}

	// 4. Собираем в исходном порядке id
	products := make([]Product, 0, len(ids))
	for _, id := range ids {
		if p, ok := found[id]; ok {
			products = append(products, p)
		}
	}

	return products, nil
}
