package product

import (
	"context"
	"reflect"
	"testing"
)

type productRepositoryFake struct {
	getProductByID func(context.Context, int64) (Product, error)
	getListProduct func(context.Context, int64, ...string) ([]Product, error)
}

func (f *productRepositoryFake) CreateProduct(context.Context, string, int64, string, int64, int64) (Product, error) {
	return Product{}, nil
}

func (f *productRepositoryFake) GetProductByID(ctx context.Context, id int64) (Product, error) {
	return f.getProductByID(ctx, id)
}

func (f *productRepositoryFake) GetListProduct(ctx context.Context, number int64, pTypes ...string) ([]Product, error) {
	return f.getListProduct(ctx, number, pTypes...)
}

func TestServiceGetListNormalizesLimit(t *testing.T) {
	var gotLimit int64
	var gotTypes []string
	service := NewService(&productRepositoryFake{
		getListProduct: func(_ context.Context, number int64, pTypes ...string) ([]Product, error) {
			gotLimit = number
			gotTypes = append([]string(nil), pTypes...)
			return []Product{{ID: 1}}, nil
		},
	})

	got, err := service.GetList(context.Background(), 11, "books", "sport")

	if err != nil {
		t.Fatalf("GetList() error = %v", err)
	}
	if gotLimit != 20 || !reflect.DeepEqual(gotTypes, []string{"books", "sport"}) || len(got) != 1 {
		t.Fatalf("unexpected result: products=%+v limit=%d types=%v", got, gotLimit, gotTypes)
	}
}

func TestServiceGetListUsesMinimumLimit(t *testing.T) {
	var gotLimit int64
	service := NewService(&productRepositoryFake{
		getListProduct: func(_ context.Context, number int64, _ ...string) ([]Product, error) {
			gotLimit = number
			return nil, nil
		},
	})

	if _, err := service.GetList(context.Background(), 0); err != nil {
		t.Fatalf("GetList() error = %v", err)
	}
	if gotLimit != 10 {
		t.Fatalf("GetList() limit = %d, want 10", gotLimit)
	}
}

func TestServiceGetProductByID(t *testing.T) {
	var gotID int64
	want := Product{ID: 9, Type: ProductType("books"), Name: "Book", Price: 100}
	service := NewService(&productRepositoryFake{
		getProductByID: func(_ context.Context, id int64) (Product, error) {
			gotID = id
			return want, nil
		},
	})

	got, err := service.GetProductByID(context.Background(), 9)

	if err != nil {
		t.Fatalf("GetProductByID() error = %v", err)
	}
	if gotID != 9 || !reflect.DeepEqual(got, want) {
		t.Fatalf("GetProductByID() = %+v, id=%d, want %+v", got, gotID, want)
	}
}
