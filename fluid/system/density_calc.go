package system

import (
	"math"
)

// Normalization constant for the Poly6 kernel adapted for 2D SPH
var poly6NormalizationConst = 4 / (math.Pi * math.Pow(smoothingRadiusH, 8))

// Müller et al – Eqn. (20) Poly6 kernel
func densitySmoothingKernel(distR float64) float64 {
	deltaRoots := smoothingRadiusH*smoothingRadiusH - distR*distR
	rootsPow3 := deltaRoots * deltaRoots * deltaRoots

	return rootsPow3 * poly6NormalizationConst
}

// Müller et al – Eqn. (3)
func calculateDensity(system *System, selected *Particle) float64 {
	var density float64 = 0

	neighborParticles := system.getParticleNeighbors(selected)

	for _, p := range neighborParticles {
		distance := selected.position.ImmutSubtract(p.position).Magnitude()

		if distance < smoothingRadiusH {
			w := densitySmoothingKernel(distance)
			density += particleMass * w
		}
	}

	// This is wrong (a temp band-aid solution). Single particles that are reintroduced to
	// to the particle stack and are at the very edge of a smoothing radius of another particle
	// produce a very-very low density. Dividing the forces by that density leads to a very
	// high acceleration (large magnitude) which ejects the particle from the stack again.
	// This fix ensures that the density will never be lower than 1 unless 0.
	// if 0 < density && density < 1 {
	// 	return 1
	// }
	return density
}
