package product

import "context"

type Service struct {
	repo ProductRepository
}

// ProductService is the application boundary used by HTTP handlers.
type ProductService interface {
	GetList(ctx context.Context, limit int64, nameType ...string) ([]Product, error)
	GetProductByID(ctx context.Context, id int64) (Product, error)
}

var _ ProductService = (*Service)(nil)

func NewService(repo ProductRepository) *Service {
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
