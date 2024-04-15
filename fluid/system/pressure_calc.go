package system

import (
	"math"
	"math/rand"

	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

var scaler = -45 / (math.Pi * math.Pow(smoothingRadiusH, 6))

// Derivative of the pressure kernel; Eqn. (21)
func pressureSmoothingKernelDerivative(distR float64) float64 {
	delta := smoothingRadiusH - distR
	return scaler * delta * delta
}

// Eqn. (12)
func calculatePressure(density float64) float64 {
	return gasConstK * (density - restDensity)
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

		delta := pPos.Subtract(selectedPos)
		distance := delta.Magnitude()

		// Check if within smoothing radius
		if distance > smoothingRadiusH || p.density == 0 {
			continue
		}

		var dir *vectors.Vector

		if distance == 0 {
			dir = vectors.NewVector(rand.Float64()/SystemScale, rand.Float64()/SystemScale)
		} else {
			dir = delta.Copy().Divide(distance)
		}

		w := pressureSmoothingKernelDerivative(distance)

		pressureMean := (selectedPressure + calculatePressure(p.density)) / 2
		scalarStep := -particleMass * pressureMean * w / p.density

		// fmt.Println("w", w)

		dir.Multiply(scalarStep)
		pressure.Add(dir)
	}

	return pressure
}
