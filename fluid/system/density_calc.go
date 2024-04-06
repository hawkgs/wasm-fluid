package system

import (
	"math"

	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

var hPow9 float64 = math.Pow(smoothingRadiusH, 9)

// Derivative of Eqn. (20)
func densitySmoothingKernelDerivative(r float64) float64 {
	if r > smoothingRadiusH {
		return 0
	}

	pow := smoothingRadiusH*smoothingRadiusH - r*r
	numerator := -945 * r * pow * pow
	denominator := 32 * math.Pi * hPow9

	return numerator / denominator
}

// Eqn. (3)
func CalculateDensityGradient(particles []*Particle, particle *vectors.Vector) float64 {
	var density float64 = 0

	for _, p := range particles {
		mag := particle.Subtract(p.position).Magnitude()
		w := densitySmoothingKernelDerivative(mag)
		density += particleMass * w
	}

	return density
}
