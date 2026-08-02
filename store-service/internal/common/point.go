package common

import "math"

const THBPerPoint = 50

func CalculatePoint(amount float64) int {
	if amount < 0 {
		return 0
	}

	points := int(math.Floor(amount / THBPerPoint))
	return points
}
