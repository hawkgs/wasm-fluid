package system

import (
	"math"
	"math/rand"

	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

var scaler = -45 / (math.Pi * math.Pow(smoothingRadiusH, 6))

// Derivative of the pressure kernel; Eqn. (21)
func pressureSmoothingKernelDerivative(r float64) float64 {
	delta := smoothingRadiusH - r
	return scaler * delta * delta
}

// Eqn. (12)
func calculatePressure(density float64) float64 {
	return gasConstK * (density - restDensity)
}

// Eqn. (10)
func calculatePressureGradient(system *System, selected *Particle) *vectors.Vector {
	pressure := vectors.NewVector(0, 0)

	// neighborParticles := system.getParticleNeighbors(selected)
	selectedPressure := calculatePressure(selected.density)

	for _, p := range system.particles {
		if selected == p {
			continue
		}

		pPos := p.position.Copy()
		selectedPos := selected.position.Copy()

		delta := pPos.Subtract(selectedPos)
		deltaMag := delta.Magnitude()

		// Check if withing smoothing radius ....
		if deltaMag > smoothingRadiusH {
			continue
		}

		var dir *vectors.Vector

		if deltaMag == 0 {
			dir = vectors.NewVector(rand.Float64(), rand.Float64())
		} else {
			dir = delta.Copy().Divide(deltaMag)
		}

		w := pressureSmoothingKernelDerivative(deltaMag)

		pressureMean := (selectedPressure + calculatePressure(p.density)) / 2
		scalarStep := -particleMass * pressureMean * w / p.density

		dir.Multiply(scalarStep)
		pressure.Add(dir)
	}

	return pressure
}
