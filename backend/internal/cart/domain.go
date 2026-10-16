package cart

import "errors"

var (
	ErrInvalidProductID = errors.New("product id must be positive")
	ErrInvalidQuantity  = errors.New("quantity must be positive")
)

type Item struct {
	ProductID int64
	Quantity  int64
}
