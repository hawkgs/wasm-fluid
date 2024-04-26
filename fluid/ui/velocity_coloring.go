package ui

import (
	"math"
	"strconv"
)

const (
	maxVelocity = 20        // Based on a downscaled field
	stepSize    = 1.0 / 3.0 // 0.33, equal to the gradient steps
)

type color = [3]float64

type colorData struct {
	step  float64
	color *color
}

// Velocity color gradient (red-orange-green-blue)
var gradient = []*colorData{
	{
		step:  0,
		color: &color{0, 127, 255},
	},
	{
		step:  0.33,
		color: &color{30, 201, 103},
	},
	{
		step:  0.66,
		color: &color{255, 143, 0},
	},
	{
		step:  1,
		color: &color{224, 86, 86},
	},
}

// getColorMix mixes two colors by a provided influence (for the 2nd color)
func getColorMix(color1 *color, color2 *color, clr2Influence float64) *color {
	ci2 := clr2Influence
	ci1 := 1 - ci2

	return &color{
		math.Round(color1[0]*ci1 + color2[0]*ci2),
		math.Round(color1[1]*ci1 + color2[1]*ci2),
		math.Round(color1[2]*ci1 + color2[2]*ci2),
	}
}

// toColorString returns a JS-compatible RGB color string
func toColorString(color *color) string {
	return "rgb(" +
		strconv.Itoa(int(color[0])) + ", " +
		strconv.Itoa(int(color[1])) + ", " +
		strconv.Itoa(int(color[2])) + ")"
}

// GetParticleVelocityColor determines the color of the particle based on its velocity
func GetParticleVelocityColor(v float64) string {
	if v >= maxVelocity {
		return toColorString(gradient[len(gradient)-1].color)
	}

	vNorm := v / maxVelocity

	for i := len(gradient) - 1; i >= 0; i-- {
		curr := gradient[i]

		if vNorm >= curr.step {
			influence := (vNorm - curr.step) / stepSize
			newColor := getColorMix(curr.color, gradient[i+1].color, influence)

			return toColorString(newColor)
		}
	}

	// Should never be called, if the V >= 0
	return toColorString(gradient[0].color)
}
