package cart

import (
	"context"
	"sort"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCart(ctx context.Context, userID int64) ([]Item, error) {
	items, err := s.repo.GetCart(ctx, userID)
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ProductID < items[j].ProductID })
	return items, nil
}

func (s *Service) AddItem(ctx context.Context, userID, productID, quantity int64) (Item, error) {
	if productID <= 0 {
		return Item{}, ErrInvalidProductID
	}
	if quantity <= 0 {
		return Item{}, ErrInvalidQuantity
	}
	return s.repo.AddItem(ctx, userID, productID, quantity)
}

func (s *Service) SetItem(ctx context.Context, userID, productID, quantity int64) (Item, error) {
	if productID <= 0 {
		return Item{}, ErrInvalidProductID
	}
	if quantity <= 0 {
		return Item{}, ErrInvalidQuantity
	}
	return s.repo.SetItem(ctx, userID, productID, quantity)
}

func (s *Service) DeleteItem(ctx context.Context, userID, productID int64) error {
	if productID <= 0 {
		return ErrInvalidProductID
	}
	return s.repo.DeleteItem(ctx, userID, productID)
}

func (s *Service) ClearCart(ctx context.Context, userID int64) error {
	return s.repo.ClearCart(ctx, userID)
}
