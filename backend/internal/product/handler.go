package product

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	serv ProductService
}

func NewHandler(service ProductService) *Handler {
	return &Handler{serv: service}
}

func (h *Handler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	var req ListRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	products, err := h.serv.GetList(r.Context(), req.Number, req.ProductType...)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var resp ListResponse
	for _, prod := range products {
		el := ListElementResponse{prod.ID, prod.Type, prod.Name, prod.Price}
		resp.Prosucts = append(resp.Prosucts, el)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	p, err := h.serv.GetProductByID(r.Context(), id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := ProductResponse{p.ID, p.Type, p.SellerID, p.Name, p.Price, p.StockQuantity}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
