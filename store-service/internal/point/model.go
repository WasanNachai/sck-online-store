package point

type SubmitedPoint struct {
	Amount int `json:"amount"`
}

type Point struct {
	OrgID  int `json:"orgId"`
	UserID int `json:"userId"`
	Amount int `json:"amount"`
}

type TotalPoint struct {
	Point int `json:"point"`
}

type CalculatePointRequest struct {
	AmountTHB float64 `json:"amountTHB"`
}

type CalculatePointResponse struct {
	EarnedPoints    int     `json:"earnedPoints"`
	AmountTHB       float64 `json:"amountTHB"`
	RateTHBPerPoint int     `json:"rateTHBPerPoint"`
}
