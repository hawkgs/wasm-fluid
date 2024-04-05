package vectors

import "math"

// RadToDeg converts radians to degrees
func RadToDeg(rad float64) float64 {
	return rad * 180 / math.Pi
}

// DegToRad converts degrees to radians
func DegToRad(deg float64) float64 {
	return 2 * math.Pi * (deg / 360)
}
