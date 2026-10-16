package product

import "errors"

var (
	ErrProductNotFound          = errors.New("product not found")
	ErrProductTypeNotFound      = errors.New("product type not found")
	ErrProductTypeAlreadyExists = errors.New("product type with this code already exists")
	ErrInvalidQuantity          = errors.New("quantity must be positive")
	ErrInsufficientStock        = errors.New("insufficient product stock")
	ErrWrongIdentificator       = errors.New("wrong identificator given")
	ErrNegativeNumber           = errors.New("number must be positive")
)

// Product is a row from the products table.
type Product struct {
	ID            int64       `redis:"id"`
	Type          ProductType `redis:"type"`
	SellerID      int64       `redis:"seller_id"`
	Name          string      `redis:"name"`
	Price         int64       `redis:"price"`
	StockQuantity int64       `redis:"stock_quantity"`
}

type ProductType string

const (
	ProductTypeFood   ProductType = "food"
	ProductTypeSport  ProductType = "sport"
	ProducTypeClothes ProductType = "clothes"
)
