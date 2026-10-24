package cart

import (
	"encoding/json"
	"errors"
	"market/internal/platform"
	"net/http"
	"strconv"
)

type Handler struct {
	serv CartService
}

func NewHandler(service CartService) *Handler {
	return &Handler{serv: service}
}

func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := platform.UserIDFromContext(r.Context())
	items, err := h.serv.GetCart(r.Context(), userID)
	if err != nil {
		h.writeError(w, err)
		return
	}

	resp := CartResponse{Items: make([]ItemResponse, 0, len(items))}
	for _, item := range items {
		resp.Items = append(resp.Items, ItemResponse{ProductID: item.ProductID, Quantity: item.Quantity})
	}
	h.writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := platform.UserIDFromContext(r.Context())
	var req ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	item, err := h.serv.AddItem(r.Context(), userID, req.ProductID, req.Quantity)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, ItemResponse{ProductID: item.ProductID, Quantity: item.Quantity})
}

func (h *Handler) SetItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := platform.UserIDFromContext(r.Context())
	productID, err := h.productID(r)
	if err != nil {
		h.writeError(w, err)
		return
	}

	var req QuantityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	item, err := h.serv.SetItem(r.Context(), userID, productID, req.Quantity)
	if err != nil {
		h.writeError(w, err)
		return
	}
	h.writeJSON(w, http.StatusOK, ItemResponse{ProductID: item.ProductID, Quantity: item.Quantity})
}

func (h *Handler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, _ := platform.UserIDFromContext(r.Context())
	productID, err := h.productID(r)
	if err != nil {
		h.writeError(w, err)
		return
	}
	if err := h.serv.DeleteItem(r.Context(), userID, productID); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, _ := platform.UserIDFromContext(r.Context())
	if err := h.serv.ClearCart(r.Context(), userID); err != nil {
		h.writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) productID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("product_id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, ErrInvalidProductID
	}
	return id, nil
}

func (h *Handler) writeError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrInvalidProductID) || errors.Is(err, ErrInvalidQuantity) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Error(w, "internal server error", http.StatusInternalServerError)
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
