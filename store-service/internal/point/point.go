package point

import (
	"context"
	"fmt"
	"log/slog"
)

type PointInterface interface {
	TotalPoint(
		ctx context.Context,
		uid int,
	) (TotalPoint, error)

	DeductPoint(
		ctx context.Context,
		uid int,
		submitedPoint SubmitedPoint,
	) (TotalPoint, error)

	CheckBurnPoint(
		ctx context.Context,
		uid int,
		amount int,
	) (bool, error)

	CalculateEarnedPoints(
		ctx context.Context,
		amountTHB float64,
	) (int, error)
}

type PointService struct {
	PointGateway PointGatewayInterface
}

type PointGatewayInterface interface {
	GetPoints(
		ctx context.Context,
		uid int,
	) ([]Point, error)

	CreatePoint(
		ctx context.Context,
		uid int,
		body Point,
	) (Point, error)

	CalculateEarnedPoints(
		ctx context.Context,
		amountTHB float64,
	) (CalculatePointResponse, error)
}

func (pointService PointService) TotalPoint(
	ctx context.Context,
	uid int,
) (TotalPoint, error) {
	points, err := pointService.PointGateway.GetPoints(ctx, uid)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"PointGateway.GetPoints failed",
			"log_type", "error",
			"error_code", "POINT_GATEWAY_FAILED",
			"error_message", err.Error(),
			"user_id", uid,
		)

		return TotalPoint{}, err
	}

	total := 0

	for _, point := range points {
		total += point.Amount
	}

	return TotalPoint{
		Point: total,
	}, nil
}

func (pointService PointService) DeductPoint(
	ctx context.Context,
	uid int,
	submitedPoint SubmitedPoint,
) (TotalPoint, error) {
	_, err := pointService.CheckBurnPoint(
		ctx,
		uid,
		submitedPoint.Amount,
	)
	if err != nil {
		return TotalPoint{}, err
	}

	point := Point{
		OrgID:  1,
		UserID: uid,
		Amount: submitedPoint.Amount,
	}

	_, err = pointService.PointGateway.CreatePoint(
		ctx,
		uid,
		point,
	)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"PointGateway.CreatePoint failed",
			"log_type", "error",
			"error_code", "POINT_CREATE_FAILED",
			"error_message", err.Error(),
			"user_id", uid,
			"amount", submitedPoint.Amount,
		)

		return TotalPoint{}, err
	}

	return pointService.TotalPoint(ctx, uid)
}

func (pointService PointService) CheckBurnPoint(
	ctx context.Context,
	uid int,
	amount int,
) (bool, error) {
	total, err := pointService.TotalPoint(ctx, uid)
	if err != nil {
		slog.ErrorContext(
			ctx,
			"PointService.TotalPoint failed",
			"log_type", "error",
			"error_code", "POINT_CHECK_FAILED",
			"error_message", err.Error(),
			"user_id", uid,
		)

		return false, err
	}

	if amount+total.Point < 0 {
		return false, fmt.Errorf(
			"points are not enough, please try again",
		)
	}

	return true, nil
}

func (pointService PointService) CalculateEarnedPoints(
	ctx context.Context,
	amountTHB float64,
) (int, error) {
	response, err :=
		pointService.PointGateway.CalculateEarnedPoints(
			ctx,
			amountTHB,
		)

	if err != nil {
		slog.ErrorContext(
			ctx,
			"PointGateway.CalculateEarnedPoints failed",
			"log_type", "error",
			"error_code", "POINT_CALCULATION_FAILED",
			"error_message", err.Error(),
			"amount_thb", amountTHB,
		)

		return 0, err
	}

	return response.EarnedPoints, nil
}
