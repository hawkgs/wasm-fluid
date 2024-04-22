package system

import (
	"math"

	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

// Adapted for 2D SPH
var spikyNormalizationConst = -30 / (math.Pi * math.Pow(smoothingRadiusH, 5))

// Derivative of the pressure kernel; Eqn. (21)
func pressureSmoothingKernelDerivative(distR float64) float64 {
	delta := smoothingRadiusH - distR
	return spikyNormalizationConst * delta * delta
}

// Eqn. (12)
func calculatePressure(density float64) float64 {
	p := gasConstK * (density - restDensity)

	// fmt.Println(p)

	return p
}

// Eqn. (10)
func calculatePressureGradient(system *System, selected *Particle) *vectors.Vector {
	pressure := vectors.NewVector(0, 0)

	if selected.density == 0 {
		return pressure
	}

	neighborParticles := system.getParticleNeighbors(selected)
	selectedPressure := calculatePressure(selected.density)

	for _, p := range neighborParticles {
		pPos := p.position.Copy()
		selectedPos := selected.position.Copy()

		delta := selectedPos.Subtract(pPos)
		distance := delta.Magnitude()

		// Check if within smoothing radius
		if distance > smoothingRadiusH {
			continue
		}

		w := pressureSmoothingKernelDerivative(distance)

		pressureMean := (selectedPressure + calculatePressure(p.density)) / 2
		scalarStep := -particleMass * pressureMean * w / p.density

		dir := delta.Copy().Normalize()
		dir.Multiply(scalarStep)
		pressure.Add(dir)
	}

	return pressure
}
