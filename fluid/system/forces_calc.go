package system

import (
	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

// Gravity force
// var gravity = vectors.NewVector(0, gravityForce)

// Derivative of the pressure kernel; Müller et al – Eqn. (21) Spiky kernel
func pressureSmoothingKernelDerivative(distR float64, cfg *SystemConfig) float64 {
	delta := cfg.SmoothingRadiusH - distR
	return cfg.spikyNormalizationConst * delta * delta
}

// Laplacian of the viscosity kernel; Müller et al – Eqn. (22) Viscosity kernel
func viscositySmoothingKernelLaplacian(distR float64, cfg *SystemConfig) float64 {
	return cfg.viscosityNormalizationConst * (cfg.SmoothingRadiusH - distR)
}

// Müller et al – Eqn. (12)
func calculatePressure(density float64, cfg *SystemConfig) float64 {
	return cfg.GasConstK * (density - cfg.RestDensity)
}

func calculateNavierStokesForces(system *System, selected *Particle) *vectors.Vector {
	c := system.config
	pressure := vectors.NewVector(0, 0)
	viscosity := vectors.NewVector(0, 0)
	gravity := vectors.NewVector(0, c.GravityForce)

	// If a sole particle (density = 0), return only ext. forces
	if selected.density == 0 {
		return gravity
	}

	neighborParticles := system.getParticleNeighbors(selected)
	selectedPressure := calculatePressure(selected.density, c)

	for _, p := range neighborParticles {
		delta := selected.position.ImmutSubtract(p.position)
		distance := delta.Magnitude()

		// Same – if a density = 0, skip
		if distance > system.config.SmoothingRadiusH || p.density == 0 {
			continue
		}

		// Calculate pressure gradient; Müller et al – Eqn. (10)
		wPres := pressureSmoothingKernelDerivative(distance, c)

		pressureMean := (selectedPressure + calculatePressure(p.density, c)) / 2
		presScalarStep := -c.ParticleMass * pressureMean * wPres / p.density

		dir := delta.Normalized()
		dir.Multiply(presScalarStep)

		pressure.Add(dir)

		// Calculate viscosity Laplacian; Müller et al – Eqn. (14)
		wVisc := viscositySmoothingKernelLaplacian(distance, c)

		viscScalarStep := c.ParticleMass * c.ViscosityConst * wVisc
		velocityDiff := p.velocity.ImmutSubtract(selected.velocity).ImmutDivide(p.density)
		velocityDiff.Multiply(viscScalarStep)

		viscosity.Add(velocityDiff)
	}

	// Add all forces: Navier-Stokes equation; Müller et al – Eqn. (7)
	// => Pressure + Viscosity + Ext. forces (i.e. Gravity)
	return pressure.ImmutAdd(viscosity).ImmutAdd(gravity)
}
