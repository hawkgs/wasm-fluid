package system

import (
	"math"
)

var densityKernelScaler = 315 / (64 * math.Pi * math.Pow(smoothingRadiusH, 9))

// Eqn. (20) poly6
func densitySmoothingKernel(distR float64) float64 {
	deltaRoots := smoothingRadiusH*smoothingRadiusH - distR*distR
	rootsPow3 := deltaRoots * deltaRoots * deltaRoots

	return rootsPow3 * densityKernelScaler
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
		if distance >= smoothingRadiusH {
			continue
		}

		w := densitySmoothingKernel(distance)

		density += particleMass * w
	}

	// fmt.Println("density", density)

	// Check why density is rapidly approaching zero which
	// causes a rapid inflection in the calc. +0.05 is
	// a band aid
	return density + 0.1
}
