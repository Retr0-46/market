package cart

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type cartRepositoryFake struct {
	getCart    func(context.Context, int64) ([]Item, error)
	addItem    func(context.Context, int64, int64, int64) (Item, error)
	setItem    func(context.Context, int64, int64, int64) (Item, error)
	deleteItem func(context.Context, int64, int64) error
	clearCart  func(context.Context, int64) error
}

func (f *cartRepositoryFake) GetCart(ctx context.Context, userID int64) ([]Item, error) {
	return f.getCart(ctx, userID)
}

func (f *cartRepositoryFake) AddItem(ctx context.Context, userID, productID, quantity int64) (Item, error) {
	return f.addItem(ctx, userID, productID, quantity)
}

func (f *cartRepositoryFake) SetItem(ctx context.Context, userID, productID, quantity int64) (Item, error) {
	return f.setItem(ctx, userID, productID, quantity)
}

func (f *cartRepositoryFake) DeleteItem(ctx context.Context, userID, productID int64) error {
	return f.deleteItem(ctx, userID, productID)
}

func (f *cartRepositoryFake) ClearCart(ctx context.Context, userID int64) error {
	return f.clearCart(ctx, userID)
}

func TestServiceGetCartSortsItems(t *testing.T) {
	service := NewService(&cartRepositoryFake{
		getCart: func(context.Context, int64) ([]Item, error) {
			return []Item{{ProductID: 3, Quantity: 1}, {ProductID: 1, Quantity: 2}, {ProductID: 2, Quantity: 4}}, nil
		},
	})

	got, err := service.GetCart(context.Background(), 10)

	if err != nil {
		t.Fatalf("GetCart() error = %v", err)
	}
	want := []Item{{ProductID: 1, Quantity: 2}, {ProductID: 2, Quantity: 4}, {ProductID: 3, Quantity: 1}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("GetCart() = %+v, want %+v", got, want)
	}
}

func TestServiceAddItem(t *testing.T) {
	var gotUserID, gotProductID, gotQuantity int64
	service := NewService(&cartRepositoryFake{
		addItem: func(_ context.Context, userID, productID, quantity int64) (Item, error) {
			gotUserID, gotProductID, gotQuantity = userID, productID, quantity
			return Item{ProductID: productID, Quantity: quantity}, nil
		},
	})

	got, err := service.AddItem(context.Background(), 4, 8, 3)

	if err != nil {
		t.Fatalf("AddItem() error = %v", err)
	}
	if got != (Item{ProductID: 8, Quantity: 3}) || gotUserID != 4 || gotProductID != 8 || gotQuantity != 3 {
		t.Fatalf("unexpected result: %+v, args=(%d,%d,%d)", got, gotUserID, gotProductID, gotQuantity)
	}
}

func TestServiceAddItemValidation(t *testing.T) {
	tests := []struct {
		name      string
		productID int64
		quantity  int64
		wantErr   error
	}{
		{name: "invalid product", productID: 0, quantity: 1, wantErr: ErrInvalidProductID},
		{name: "invalid quantity", productID: 1, quantity: 0, wantErr: ErrInvalidQuantity},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			service := NewService(&cartRepositoryFake{
				addItem: func(context.Context, int64, int64, int64) (Item, error) {
					called = true
					return Item{}, nil
				},
			})

			_, err := service.AddItem(context.Background(), 1, tt.productID, tt.quantity)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("AddItem() error = %v, want %v", err, tt.wantErr)
			}
			if called {
				t.Fatal("repository was called for invalid input")
			}
		})
	}
}

func TestServiceSetDeleteAndClear(t *testing.T) {
	var setArgs, deleteArgs []int64
	var clearedUserID int64
	service := NewService(&cartRepositoryFake{
		setItem: func(_ context.Context, userID, productID, quantity int64) (Item, error) {
			setArgs = []int64{userID, productID, quantity}
			return Item{ProductID: productID, Quantity: quantity}, nil
		},
		deleteItem: func(_ context.Context, userID, productID int64) error {
			deleteArgs = []int64{userID, productID}
			return nil
		},
		clearCart: func(_ context.Context, userID int64) error {
			clearedUserID = userID
			return nil
		},
	})

	if _, err := service.SetItem(context.Background(), 2, 5, 6); err != nil {
		t.Fatalf("SetItem() error = %v", err)
	}
	if err := service.DeleteItem(context.Background(), 2, 5); err != nil {
		t.Fatalf("DeleteItem() error = %v", err)
	}
	if err := service.ClearCart(context.Background(), 2); err != nil {
		t.Fatalf("ClearCart() error = %v", err)
	}

	if !reflect.DeepEqual(setArgs, []int64{2, 5, 6}) || !reflect.DeepEqual(deleteArgs, []int64{2, 5}) || clearedUserID != 2 {
		t.Fatalf("unexpected repository calls: set=%v delete=%v clear=%d", setArgs, deleteArgs, clearedUserID)
	}
}
