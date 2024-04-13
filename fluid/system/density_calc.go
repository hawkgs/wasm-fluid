package system

import (
	"fmt"
	"math"
)

var densityKernelScaler = 315 / (64 * math.Pi * math.Pow(smoothingRadiusH, 9))

// Eqn. (20) poly6
func densitySmoothingKernel(r float64) float64 {
	deltaRoots := smoothingRadiusH*smoothingRadiusH - r*r
	rootsPow3 := deltaRoots * deltaRoots * deltaRoots

	return rootsPow3 * densityKernelScaler
}

// Eqn. (3)
func calculateDensity(system *System, selected *Particle) float64 {
	var density float64 = 0

	// neighborParticles := system.getParticleNeighbors(selected)

	for _, p := range system.particles {
		if selected == p {
			continue
		}

		pPos := p.position.Copy()
		selectedPos := selected.position.Copy()
		deltaMag := selectedPos.Subtract(pPos).Magnitude()

		// Check if withing smoothing radius ....
		if deltaMag > smoothingRadiusH {
			continue
		}

		w := densitySmoothingKernel(deltaMag)

		density += particleMass * w
	}

	fmt.Println("density", density)

	// Check why density is rapidly approaching zero which
	// causes a rapid inflection in the calc. +0.05 is
	// a band aid
	return density + 0.05
}
