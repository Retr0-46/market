package product

import "errors"

var (
	ErrProductNotFound          = errors.New("product not found")
	ErrProductTypeNotFound      = errors.New("product type not found")
	ErrProductTypeAlreadyExists = errors.New("product type with this code already exists")
	ErrInvalidQuantity          = errors.New("quantity must be positive")
	ErrInsufficientStock        = errors.New("insufficient product stock")
)

// Product is a row from the products table.
type Product struct {
	ID            int64
	TypeID        int64
	SellerID      int64
	Name          string
	Price         int64
	StockQuantity int64
}

// ProductType is a row from the product_types table.
type ProductType struct {
	ID   int64
	Code string
	Name string
}
