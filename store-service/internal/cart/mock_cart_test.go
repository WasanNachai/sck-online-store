package cart_test

import (
	"context"
	"store-service/internal/cart"
	"store-service/internal/point"

	"github.com/stretchr/testify/mock"
)

type mockPointInterface struct {
	mock.Mock
}

func (mock *mockPointInterface) CalculateEarnedPoints(
	ctx context.Context,
	amountTHB float64,
) (int, error) {
	return int(amountTHB / 50), nil
}

func (mock *mockPointInterface) TotalPoint(ctx context.Context, uid int) (point.TotalPoint, error) {
	argument := mock.Called(ctx, uid)
	return argument.Get(0).(point.TotalPoint), argument.Error(1)
}

func (mock *mockPointInterface) DeductPoint(ctx context.Context, uid int, submitedPoint point.SubmitedPoint) (point.TotalPoint, error) {
	argument := mock.Called(ctx, uid, submitedPoint)
	return argument.Get(0).(point.TotalPoint), argument.Error(1)
}

func (mock *mockPointInterface) CheckBurnPoint(ctx context.Context, uid int, amount int) (bool, error) {
	argument := mock.Called(ctx, uid, amount)
	return argument.Bool(0), argument.Error(1)
}

type mockCartRepository struct {
	mock.Mock
}

func (repo *mockCartRepository) GetCartDetail(ctx context.Context, userID int) ([]cart.CartDetail, error) {
	argument := repo.Called(ctx, userID)
	return argument.Get(0).([]cart.CartDetail), argument.Error(1)
}

func (repo *mockCartRepository) GetCartByProductID(ctx context.Context, userID int, productID int) (cart.Cart, error) {
	argument := repo.Called(ctx, userID, productID)
	return argument.Get(0).(cart.Cart), argument.Error(1)
}

func (repo *mockCartRepository) CreateCart(ctx context.Context, userID int, productID int, quantity int) (int, error) {
	argument := repo.Called(ctx, userID, productID, quantity)
	return argument.Int(0), argument.Error(1)
}

func (repo *mockCartRepository) UpdateCart(ctx context.Context, userID int, productID int, quantity int) error {
	argument := repo.Called(ctx, userID, productID, quantity)
	return argument.Error(0)
}

func (repo *mockCartRepository) DeleteCart(ctx context.Context, userID int, productID int) error {
	argument := repo.Called(ctx, userID, productID)
	return argument.Error(0)
}
