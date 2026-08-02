package point_test

import (
	"context"
	"store-service/internal/point"

	"github.com/stretchr/testify/mock"
)

type mockPointGateway struct {
	mock.Mock
}

func (gateway *mockPointGateway) GetPoints(ctx context.Context, userID int) ([]point.Point, error) {
	argument := gateway.Called(ctx, userID)
	return argument.Get(0).([]point.Point), argument.Error(1)
}

func (gateway *mockPointGateway) CreatePoint(ctx context.Context, userID int, pointItem point.Point) (point.Point, error) {
	argument := gateway.Called(ctx, userID, pointItem)
	return argument.Get(0).(point.Point), argument.Error(1)
}

func (m *mockPointGateway) CalculateEarnedPoints(
	ctx context.Context,
	amountTHB float64,
) (point.CalculatePointResponse, error) {
	args := m.Called(ctx, amountTHB)

	var result point.CalculatePointResponse
	if args.Get(0) != nil {
		result = args.Get(0).(point.CalculatePointResponse)
	}

	return result, args.Error(1)
}
