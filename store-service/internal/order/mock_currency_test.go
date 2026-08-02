package order_test

import (
	"context"
	"store-service/internal/common"
)

type mockCurrencyService struct{}

func (m mockCurrencyService) ConvertToThb(_ context.Context, amount float64) common.Decimal {
	return common.ConvertToThb(amount, common.DefaultUSDToTHBRate)
}
