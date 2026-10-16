package cart

type ItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

type QuantityRequest struct {
	Quantity int64 `json:"quantity"`
}

type ItemResponse struct {
	ProductID int64 `json:"product_id"`
	Quantity  int64 `json:"quantity"`
}

type CartResponse struct {
	Items []ItemResponse `json:"items"`
}
