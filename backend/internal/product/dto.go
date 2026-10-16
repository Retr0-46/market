package product

type ProductResponse struct {
	ID            int64       `json:"id"`
	Type          ProductType `json:"type"`
	SellerID      int64       `json:"seller_id"`
	Name          string      `json:"name"`
	Price         int64       `json:"price"`
	StockQuantity int64       `json:"stock_quantity"`
}

type ProductRequest struct {
	ID int64 `json:"id"`
}

type ListRequest struct {
	Number      int64    `json:"number"`
	ProductType []string `json:"p_types"`
}

type ListResponse struct {
	Prosucts []ListElementResponse `json:"product"`
}

type ListElementResponse struct {
	ID    int64       `json:"id"`
	Type  ProductType `json:"type"`
	Name  string      `json:"name"`
	Price int64       `json:"price"`
}
