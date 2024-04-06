package system

import (
	"math"

	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

var hPow6 float64 = math.Pow(smoothingRadiusH, 6)

// Derivative of the pressure kernel; Eqn. (21)
func pressureSmoothingKernelDerivative(r float64) float64 {
	if r > smoothingRadiusH {
		return 0
	}

	pow := smoothingRadiusH - r
	numerator := -45 * pow * pow
	denominator := math.Pi * hPow6

	return numerator / denominator
}

// Eqn. (12)
func CalculatePressure(density float64) float64 {
	return gasConstK * (density - restDensity)
}

// Eqn. (10)
func CalculatePressureGradient(particles []*Particle, particle *vectors.Vector) float64 {
	var pressure float64 = 0

	for _, p := range particles {
		mag := particle.Subtract(p.position).Magnitude()
		w := pressureSmoothingKernelDerivative(mag)

		// Todo: multiply by the arithmetic mean of the pressures of the interacting
		// particles => (p(particle) + p(ParticleI)) / 2 * density(ParticleI)
		pressure += -particleMass * w
	}

	return pressure
}
