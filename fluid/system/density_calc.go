package system

import (
	"math"
)

// Adapted for 2D SPH
var poly6NormalizationConst = 4 / (math.Pi * math.Pow(smoothingRadiusH, 8))

// var densityKernelScaler = 1 / (3 * math.Pi * math.Pow(smoothingRadiusH, 6))

// Eqn. (20) poly6
func densitySmoothingKernel(distR float64) float64 {
	deltaRoots := smoothingRadiusH*smoothingRadiusH - distR*distR
	rootsPow3 := deltaRoots * deltaRoots * deltaRoots

	return rootsPow3 * poly6NormalizationConst
}

// Eqn. (3)
func calculateDensity(system *System, selected *Particle) float64 {
	var density float64 = 0

	neighborParticles := system.getParticleNeighbors(selected)

	for _, p := range neighborParticles {
		pPos := p.position.Copy()
		selectedPos := selected.position.Copy()
		distance := selectedPos.Subtract(pPos).Magnitude()

		// Check if within smoothing radius
		if distance < smoothingRadiusH {
			w := densitySmoothingKernel(distance)
			density += particleMass * w
		}
	}

	return density
}
