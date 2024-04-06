package utils

import "math"

func Round(num float64, precision uint) float64 {
	multiplier := math.Pow(10, float64(precision))

	return math.Round(num*multiplier) / multiplier
}
