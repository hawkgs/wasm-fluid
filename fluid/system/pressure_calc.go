package system

import (
	"math"

	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

// Adapted for 2D SPH
var spikyNormalizationConst = -30 / (math.Pi * math.Pow(smoothingRadiusH, 5))

// Derivative of the pressure kernel; Eqn. (21) spiky
func pressureSmoothingKernelDerivative(distR float64) float64 {
	delta := smoothingRadiusH - distR
	return spikyNormalizationConst * delta * delta
}

// Adapted for 2d sph
var viscosityNormalizationConst = 40 / (math.Pi * math.Pow(smoothingRadiusH, 5))

// Laplacian of the viscosity kernel; Eqn. (22) visc
func viscositySmoothingKernelLaplacian(distR float64) float64 {
	return viscosityNormalizationConst * (smoothingRadiusH - distR)
}

// Eqn. (12)
func calculatePressure(density float64) float64 {
	return gasConstK * (density - restDensity)
}

func calculateNavierStokesForces(system *System, selected *Particle) *vectors.Vector {
	pressure := vectors.NewVector(0, 0)
	viscosity := vectors.NewVector(0, 0)

	// avoid nan
	if selected.density == 0 {
		return pressure
	}

	neighborParticles := system.getParticleNeighbors(selected)
	selectedPressure := calculatePressure(selected.density)

	for _, p := range neighborParticles {
		delta := selected.position.ImmutSubtract(p.position)
		distance := delta.Magnitude()

		// Check if within smoothing radius
		if distance >= smoothingRadiusH {
			continue
		}

		// Eqn. (10)
		// calculate pressure
		wPres := pressureSmoothingKernelDerivative(distance)

		pressureMean := (selectedPressure + calculatePressure(p.density)) / 2
		presScalarStep := -particleMass * pressureMean * wPres / p.density

		dir := delta.ImmutNormalize()
		dir.Multiply(presScalarStep)

		pressure.Add(dir)

		// calculate viscosity; Eqn. (14)
		wVisc := viscositySmoothingKernelLaplacian(distance)
		viscScalarStep := particleMass * viscosityConst * wVisc
		velocityDiff := p.velocity.ImmutSubtract(selected.velocity).ImmutDivide(p.density)
		velocityDiff.Multiply(viscScalarStep)

		viscosity.Add(velocityDiff)
	}

	return pressure.ImmutAdd(viscosity)
}
