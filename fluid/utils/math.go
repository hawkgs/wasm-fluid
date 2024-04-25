package utils

import "math"

// RoundNum rounds a number with a given decimal point precision
func RoundNum(num float64, precision uint) float64 {
	multiplier := math.Pow(10, float64(precision))

	return math.Round(num*multiplier) / multiplier
}
