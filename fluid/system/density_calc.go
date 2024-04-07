package system

import (
	"math"
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
func calculateDensityGradient(system *System, particle *Particle) float64 {
	var density float64 = 0

	neighborParticles := system.getParticleNeighbors(particle)

	for _, p := range neighborParticles {
		if particle == p {
			continue
		}

		pPos := p.position.Copy()
		particlePos := particle.position.Copy()

		mag := pPos.Subtract(particlePos).Magnitude()

		w := densitySmoothingKernelDerivative(mag)
		density += particleMass * w
	}

	return density
}
