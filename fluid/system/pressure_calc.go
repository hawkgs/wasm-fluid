package system

import (
	"math"
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
func calculatePressureGradient(system *System, particle *Particle) float64 {
	var pressure float64 = 0

	neighborParticles := system.getParticleNeighbors(particle)
	particlePressure := CalculatePressure(particle.density)

	for _, p := range neighborParticles {
		if particle == p {
			continue
		}

		pPos := p.position.Copy()
		particlePos := particle.position.Copy()

		mag := pPos.Subtract(particlePos).Magnitude()
		w := pressureSmoothingKernelDerivative(mag)

		// Todo: multiply by the arithmetic mean of the pressures of the interacting
		// particles => (p(particle) + p(ParticleI)) / 2 * density(ParticleI)
		pressureMean := (particlePressure + CalculatePressure(p.density)) / (2 * p.density)
		pressure += -particleMass * w * pressureMean
	}

	return pressure
}
