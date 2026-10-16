package product

import "context"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetList(ctx context.Context, limit int64, nameType ...string) ([]Product, error) {
	var products []Product
	if limit < 10 {
		limit = 10
	}
	if limit%10 > 0 {
		limit += 10 - (limit % 10)
	}

	products, err := s.repo.GetListProduct(ctx, limit, nameType...)

	return products, err
}

func (s *Service) GetProductByID(ctx context.Context, id int64) (Product, error) {
	return s.repo.GetProductByID(ctx, id)
}
