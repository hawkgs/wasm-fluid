package system

import (
	"math"

	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

// Gravity vector
var gravityVector = vectors.NewVector(0, gravityForce)

// Normalization constant for the Spiky kernel adapted for 2D SPH
var spikyNormalizationConst = -30 / (math.Pi * math.Pow(smoothingRadiusH, 5))

// Derivative of the pressure kernel; Müller et al – Eqn. (21) Spiky kernel
func pressureSmoothingKernelDerivative(distR float64) float64 {
	delta := smoothingRadiusH - distR
	return spikyNormalizationConst * delta * delta
}

// Normalization constant for the Viscosity kernel adapted for 2D SPH
var viscosityNormalizationConst = 40 / (math.Pi * math.Pow(smoothingRadiusH, 5))

// Laplacian of the viscosity kernel; Müller et al – Eqn. (22) Viscosity kernel
func viscositySmoothingKernelLaplacian(distR float64) float64 {
	return viscosityNormalizationConst * (smoothingRadiusH - distR)
}

// Müller et al – Eqn. (12)
func calculatePressure(density float64) float64 {
	return gasConstK * (density - restDensity)
}

func calculateNavierStokesForces(system *System, selected *Particle) *vectors.Vector {
	pressure := vectors.NewVector(0, 0)
	viscosity := vectors.NewVector(0, 0)

	// If a sole particle (density = 0), return only ext. forces
	if selected.density == 0 {
		return gravityVector.Copy()
	}

	neighborParticles := system.getParticleNeighbors(selected)
	selectedPressure := calculatePressure(selected.density)

	for _, p := range neighborParticles {
		delta := selected.position.ImmutSubtract(p.position)
		distance := delta.Magnitude()

		if distance >= smoothingRadiusH {
			continue
		}

		// Calculate pressure gradient; Müller et al – Eqn. (10)
		wPres := pressureSmoothingKernelDerivative(distance)

		pressureMean := (selectedPressure + calculatePressure(p.density)) / 2
		presScalarStep := -particleMass * pressureMean * wPres / p.density

		dir := delta.Normalized()
		dir.Multiply(presScalarStep)

		pressure.Add(dir)

		// Calculate viscosity Laplacian; Müller et al – Eqn. (14)
		wVisc := viscositySmoothingKernelLaplacian(distance)

		viscScalarStep := particleMass * viscosityConst * wVisc
		velocityDiff := p.velocity.ImmutSubtract(selected.velocity).ImmutDivide(p.density)
		velocityDiff.Multiply(viscScalarStep)

		viscosity.Add(velocityDiff)
	}

	// Add all forces: Navier-Stokes equation; Müller et al – Eqn. (7)
	// => Pressure + Viscosity + Ext. forces (i.e. Gravity)
	return pressure.ImmutAdd(viscosity).ImmutAdd(gravityVector)
}
