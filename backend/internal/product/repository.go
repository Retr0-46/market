package product

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// WithTx runs fn in a database transaction. The transaction is committed
// when fn returns nil and rolled back when fn returns an error.
func (r *Repository) WithTx(ctx context.Context, fn func(pgx.Tx) error) error {
	return pgx.BeginFunc(ctx, r.pool, fn)
}

// CreateProduct inserts a product and returns the saved row.
func (r *Repository) CreateProduct(ctx context.Context, typeID, sellerID int64, name string, price, stockQuantity int64) (Product, error) {
	var product Product

	err := r.pool.QueryRow(ctx, `
		INSERT INTO products (type_id, seller_id, name, price, stock_quantity)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, type_id, seller_id, name, price, stock_quantity
	`, typeID, sellerID, name, price, stockQuantity).Scan(
		&product.ID,
		&product.TypeID,
		&product.SellerID,
		&product.Name,
		&product.Price,
		&product.StockQuantity,
	)
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func (r *Repository) GetProduct(ctx context.Context, id int64) (Product, error) {
	var product Product

	err := r.pool.QueryRow(ctx, `
		SELECT id, type_id, seller_id, name, price, stock_quantity
		FROM products
		WHERE id = $1
	`, id).Scan(
		&product.ID,
		&product.TypeID,
		&product.SellerID,
		&product.Name,
		&product.Price,
		&product.StockQuantity,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrProductNotFound
	}
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

// ListProducts returns all products. If a type ID is supplied, only products
// belonging to that type are returned.
func (r *Repository) ListProducts(ctx context.Context, typeID ...int64) ([]Product, error) {
	var (
		rows pgx.Rows
		err  error
	)

	switch len(typeID) {
	case 0:
		rows, err = r.pool.Query(ctx, `
			SELECT id, type_id, seller_id, name, price, stock_quantity
			FROM products
			ORDER BY id
		`)
	case 1:
		rows, err = r.pool.Query(ctx, `
			SELECT id, type_id, seller_id, name, price, stock_quantity
			FROM products
			WHERE type_id = $1
			ORDER BY id
		`, typeID[0])
	default:
		return nil, errors.New("at most one product type ID may be supplied")
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.TypeID, &product.SellerID, &product.Name, &product.Price, &product.StockQuantity); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

// GetProducts is an alias for ListProducts without a type filter.
func (r *Repository) GetProducts(ctx context.Context) ([]Product, error) {
	return r.ListProducts(ctx)
}

func (r *Repository) GetProductsByType(ctx context.Context, typeID int64) ([]Product, error) {
	return r.ListProducts(ctx, typeID)
}

// GetProductsBySeller returns all products owned by a seller.
func (r *Repository) GetProductsBySeller(ctx context.Context, sellerID int64) ([]Product, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, type_id, seller_id, name, price, stock_quantity
		FROM products
		WHERE seller_id = $1
		ORDER BY id
	`, sellerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.TypeID, &product.SellerID, &product.Name, &product.Price, &product.StockQuantity); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *Repository) UpdateProduct(ctx context.Context, id, typeID, sellerID int64, name string, price, stockQuantity int64) (Product, error) {
	var product Product

	err := r.pool.QueryRow(ctx, `
		UPDATE products
		SET type_id = $2, seller_id = $3, name = $4, price = $5, stock_quantity = $6
		WHERE id = $1
		RETURNING id, type_id, seller_id, name, price, stock_quantity
	`, id, typeID, sellerID, name, price, stockQuantity).Scan(
		&product.ID,
		&product.TypeID,
		&product.SellerID,
		&product.Name,
		&product.Price,
		&product.StockQuantity,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Product{}, ErrProductNotFound
	}
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

// DecreaseStockTx atomically reserves quantity units of a product. The caller
// owns tx and decides when the whole transaction is committed.
func (r *Repository) DecreaseStockTx(ctx context.Context, tx pgx.Tx, productID, quantity int64) (Product, error) {
	if quantity <= 0 {
		return Product{}, ErrInvalidQuantity
	}

	var product Product
	err := tx.QueryRow(ctx, `
		UPDATE products
		SET stock_quantity = stock_quantity - $1
		WHERE id = $2 AND stock_quantity >= $1
		RETURNING id, type_id, seller_id, name, price, stock_quantity
	`, quantity, productID).Scan(
		&product.ID,
		&product.TypeID,
		&product.SellerID,
		&product.Name,
		&product.Price,
		&product.StockQuantity,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		var exists bool
		if existsErr := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM products WHERE id = $1)`, productID).Scan(&exists); existsErr != nil {
			return Product{}, existsErr
		}
		if !exists {
			return Product{}, ErrProductNotFound
		}
		return Product{}, ErrInsufficientStock
	}
	if err != nil {
		return Product{}, err
	}

	return product, nil
}

func (r *Repository) DeleteProduct(ctx context.Context, id int64) error {
	result, err := r.pool.Exec(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrProductNotFound
	}

	return nil
}

func (r *Repository) CreateProductType(ctx context.Context, code, name string) (ProductType, error) {
	var productType ProductType

	err := r.pool.QueryRow(ctx, `
		INSERT INTO product_types (code, name)
		VALUES ($1, $2)
		RETURNING id, code, name
	`, code, name).Scan(
		&productType.ID,
		&productType.Code,
		&productType.Name,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return ProductType{}, ErrProductTypeAlreadyExists
	}
	if err != nil {
		return ProductType{}, err
	}

	return productType, nil
}

func (r *Repository) GetProductType(ctx context.Context, id int64) (ProductType, error) {
	var productType ProductType

	err := r.pool.QueryRow(ctx, `
		SELECT id, code, name
		FROM product_types
		WHERE id = $1
	`, id).Scan(&productType.ID, &productType.Code, &productType.Name)
	if errors.Is(err, pgx.ErrNoRows) {
		return ProductType{}, ErrProductTypeNotFound
	}
	if err != nil {
		return ProductType{}, err
	}

	return productType, nil
}

func (r *Repository) ListProductTypes(ctx context.Context) ([]ProductType, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, code, name
		FROM product_types
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	types := make([]ProductType, 0)
	for rows.Next() {
		var productType ProductType
		if err := rows.Scan(&productType.ID, &productType.Code, &productType.Name); err != nil {
			return nil, err
		}
		types = append(types, productType)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return types, nil
}
