package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"store-service/internal/point"

	"github.com/gin-gonic/gin"
)

type PointAPI struct {
	PointService point.PointInterface
}

// @Summary Deduct points from user
// @Description Deduct points from user's point balance
// @Tags point
// @Accept json
// @Produce json
// @Param request body point.SubmitedPoint true "Point deduction request"
// @Success 200 {object} point.Point
// @Failure 400 {string} string "Bad request error"
// @Failure 500
// @Router /api/v1/point [post]
func (api PointAPI) DeductPointHandler(context *gin.Context) {
	ctx := context.Request.Context()

	var request point.SubmitedPoint
	if err := context.BindJSON(&request); err != nil {
		slog.ErrorContext(ctx, "Point deduct bad request",
			"log_type", "error",
			"error_code", "INVALID_REQUEST",
			"error_message", err.Error(),
			"user_id", 0,
		)
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	uid, uidErr := strconv.Atoi(context.GetHeader("uid"))
	if uidErr != nil {
		uid = context.GetInt("userID")
	}

	res, err := api.PointService.DeductPoint(ctx, uid, request)
	if err != nil {
		slog.ErrorContext(ctx, "PointService.DeductPoint failed",
			"log_type", "error",
			"error_code", "POINT_DEDUCTION_FAILED",
			"error_message", err.Error(),
			"user_id", uid,
			slog.Any("request", map[string]any{"amount": request.Amount}),
		)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.InfoContext(ctx, "Points deducted",
		"log_type", "business",
		"event", "points_deducted",
		"entity_type", "point",
		"entity_id", uid,
		"actor_id", uid,
		slog.Any("metadata", map[string]any{
			"amount":          request.Amount,
			"remaining_point": res.Point,
		}),
	)

	context.JSON(http.StatusOK, res)
}

// @Summary Get total points
// @Description Get user's total point balance
// @Tags point
// @Accept json
// @Produce json
// @Success 200 {object} point.Point
// @Failure 500
// @Router /api/v1/point [get]
func (api PointAPI) TotalPointHandler(context *gin.Context) {
	uid, uidErr := strconv.Atoi(context.GetHeader("uid"))
	if uidErr != nil {
		uid = context.GetInt("userID")
	}

	ctx := context.Request.Context()
	res, err := api.PointService.TotalPoint(ctx, uid)

	if err != nil {
		slog.ErrorContext(ctx, "PointService.TotalPoint failed",
			"log_type", "error",
			"error_code", "POINT_QUERY_FAILED",
			"error_message", err.Error(),
			"user_id", uid,
		)
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, res)
}

// @Summary Get point summary
// @Description Get user's points grouped by lifecycle statuses
// @Tags point
// @Accept json
// @Produce json
// @Success 200 {object} point.PointSummary
// @Failure 500
// @Router /api/v1/point/summary [get]
func (api PointAPI) PointSummaryHandler(context *gin.Context) {
	uid := context.GetInt("userID")
	ctx := context.Request.Context()

	res, err := api.PointService.GetPointSummary(ctx, uid)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, res)
}

// @Summary Create pending earn points
// @Description Create pending approval earned points for the authenticated user
// @Tags point
// @Accept json
// @Produce json
// @Param request body point.SubmitedPendingEarnPoint true "Pending earn point request"
// @Success 200 {object} point.Point
// @Failure 400 {string} string "Bad request error"
// @Failure 500
// @Router /api/v1/point/earn/pending [post]
func (api PointAPI) CreatePendingEarnPointHandler(context *gin.Context) {
	ctx := context.Request.Context()
	uid := context.GetInt("userID")

	var request point.SubmitedPendingEarnPoint
	if err := context.BindJSON(&request); err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	if request.Amount <= 0 {
		context.String(
			http.StatusBadRequest,
			"amount must be greater than zero",
		)
		return
	}

	res, err := api.PointService.CreatePendingEarnPoint(
		ctx,
		uid,
		request,
	)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, res)
}

// @Summary Approve pending point
// @Description Approve pending point item for authenticated user
// @Tags point
// @Accept json
// @Produce json
// @Param id path int true "Point ID"
// @Success 200 {object} point.Point
// @Failure 400 {string} string "Bad request error"
// @Failure 500
// @Router /api/v1/point/{id}/approve [patch]
func (api PointAPI) ApprovePointHandler(context *gin.Context) {
	ctx := context.Request.Context()
	uid := context.GetInt("userID")
	pointID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.String(http.StatusBadRequest, "invalid point id")
		return
	}

	res, err := api.PointService.ApprovePoint(ctx, uid, pointID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, res)
}

// @Summary Redeem approved point
// @Description Redeem approved point item for authenticated user
// @Tags point
// @Accept json
// @Produce json
// @Param id path int true "Point ID"
// @Success 200 {object} point.Point
// @Failure 400 {string} string "Bad request error"
// @Failure 500
// @Router /api/v1/point/{id}/redeem [patch]
func (api PointAPI) RedeemPointHandler(context *gin.Context) {
	ctx := context.Request.Context()
	uid := context.GetInt("userID")
	pointID, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		context.String(http.StatusBadRequest, "invalid point id")
		return
	}

	res, err := api.PointService.RedeemPoint(ctx, uid, pointID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, res)
}

// @Summary Expire approved points
// @Description Mark approved points that passed expire date as expired for authenticated user
// @Tags point
// @Accept json
// @Produce json
// @Param request body point.SubmitedExpirePoint true "Expire point request"
// @Success 200 {object} point.ExpirePointResponse
// @Failure 400 {string} string "Bad request error"
// @Failure 500
// @Router /api/v1/point/expire [patch]
func (api PointAPI) ExpirePointHandler(context *gin.Context) {
	ctx := context.Request.Context()
	uid := context.GetInt("userID")

	var request point.SubmitedExpirePoint
	if err := context.BindJSON(&request); err != nil {
		context.String(http.StatusBadRequest, err.Error())
		return
	}

	if request.BeforeDate == "" {
		context.String(http.StatusBadRequest, "before_date is required")
		return
	}

	res, err := api.PointService.ExpirePoints(
		ctx,
		uid,
		request,
	)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"expired_count": res.ExpiredCount,
	})
}
