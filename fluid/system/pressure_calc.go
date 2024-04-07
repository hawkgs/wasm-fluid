package system

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

var scaler = -45 / (math.Pi * math.Pow(smoothingRadiusH, 6))

// Derivative of the pressure kernel; Eqn. (21)
func pressureSmoothingKernelDerivative(r float64) float64 {
	if smoothingRadiusH < r {
		return 0
	}

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

	neighborParticles := system.getParticleNeighbors(selected)
	selectedPressure := calculatePressure(selected.density)

	for _, p := range neighborParticles {
		if selected == p {
			continue
		}

		pPos := p.position.Copy()
		selectedPos := selected.position.Copy()

		delta := selectedPos.Subtract(pPos)
		deltaMag := delta.Magnitude()
		var dir *vectors.Vector

		if deltaMag == 0 {
			dir = vectors.NewVector(rand.Float64(), rand.Float64())
		} else {
			dir = delta.Copy().Divide(deltaMag)
		}

		w := pressureSmoothingKernelDerivative(deltaMag)

		// Todo: multiply by the arithmetic mean of the pressures of the interacting
		// particles => (p(selParticle) + p(ParticleI)) / 2
		pressureMean := (selectedPressure + calculatePressure(p.density)) / 2
		scalarStep := -particleMass * pressureMean * w / p.density

		dir.Multiply(scalarStep)
		pressure.Add(dir)
	}

	fmt.Println(pressure)

	return pressure
}
