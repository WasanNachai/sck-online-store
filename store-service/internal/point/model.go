package point

type SubmitedPoint struct {
	Amount int `json:"amount"`
}

type SubmitedPendingEarnPoint struct {
	OrderNumber int64  `json:"order_number"`
	Amount      int    `json:"amount"`
	ExpireDate  string `json:"expire_date"`
}

type SubmitedExpirePoint struct {
	BeforeDate string `json:"before_date"`
}

type Point struct {
	ID         int    `json:"id,omitempty"`
	OrgID      int    `json:"orgId"`
	UserID     int    `json:"userId"`
	Amount     int    `json:"amount"`
	Status     string `json:"status,omitempty"`
	ExpireDate string `json:"expireDate,omitempty"`
}

type TotalPoint struct {
	Point int `json:"point"`
}

type PointSummary struct {
	AvailablePoints int `json:"available_points"`
	PendingPoints   int `json:"pending_points"`
	RedeemedPoints  int `json:"redeemed_points"`
	ExpiredPoints   int `json:"expired_points"`
}

type PointServiceSummary struct {
	AvailablePoints int `json:"availablePoints"`
	PendingPoints   int `json:"pendingPoints"`
	RedeemedPoints  int `json:"redeemedPoints"`
	ExpiredPoints   int `json:"expiredPoints"`
}

type PointTransitionRequest struct {
	UserID int `json:"userId"`
}

type CreatePendingPointRequest struct {
	OrgID      int    `json:"orgId"`
	UserID     int    `json:"userId"`
	Amount     int    `json:"amount"`
	ExpireDate string `json:"expireDate"`
}

type ExpirePointRequest struct {
	BeforeDate string `json:"beforeDate"`
	UserID     int    `json:"userId"`
}

type ExpirePointResponse struct {
	ExpiredCount int `json:"expiredCount"`
}

type CalculatePointRequest struct {
	AmountTHB float64 `json:"amountTHB"`
}

type CalculatePointResponse struct {
	EarnedPoints    int     `json:"earnedPoints"`
	AmountTHB       float64 `json:"amountTHB"`
	RateTHBPerPoint int     `json:"rateTHBPerPoint"`
}
