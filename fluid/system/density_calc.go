package system

import (
	"math"
)

var densityKernelScaler = 315 / (64 * math.Pi * math.Pow(smoothingRadiusH, 9))

// Eqn. (20) poly6
func densitySmoothingKernel(r float64) float64 {
	if r < 0 || smoothingRadiusH < r {
		return 0
	}

	deltaRoots := smoothingRadiusH*smoothingRadiusH - r*r
	rootsPow3 := deltaRoots * deltaRoots * deltaRoots

	return rootsPow3 * densityKernelScaler
}

// Eqn. (3)
func calculateDensity(system *System, selected *Particle) float64 {
	var density float64 = 0

	neighborParticles := system.getParticleNeighbors(selected)

	for _, p := range neighborParticles {
		if selected == p {
			continue
		}

		pPos := p.position.Copy()
		selectedPos := selected.position.Copy()
		deltaMag := selectedPos.Subtract(pPos).Magnitude()

		w := densitySmoothingKernel(deltaMag)

		density += particleMass * w
	}

	return density
}
