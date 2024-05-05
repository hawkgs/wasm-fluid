package system

import "math"

// StaticParams keeps all simulation parameters that are NOT eligible for real-time updates
type StaticParams struct {
	SystemScale      float64
	SmoothingRadiusH float64
	Timestep         float64 // i.e. DT; Delta Time
}

// DynamicParams keeps all simulation parameters that are eligible for real-time updates
type DynamicParams struct {
	ParticleMass     float64
	GravityForce     float64
	GasConstK        float64
	RestDensity      float64
	ViscosityConst   float64
	VelocityLimit    float64
	CollisionDamping float64
}

// SystemConfig keeps all the config properties and simulation parameters, including pre-computed values
type SystemConfig struct {
	// General
	Width            float64
	Height           float64
	Particles        uint
	ParticleUiRadius float64

	// The SystemScale downsizes the particle field.
	// For example, 800px-wide canvas will be downscaled to 20, if the SystemScale is 40.
	// This is needed because the kernel functions are very sensitive to the smoothing radius.
	// They are actually meant to work with specific smoothing radiuses. You can easily check
	// that by plotting the functions (or their derivates, resp.) and examining how they change over time.
	// SPH is very sensitive to parameter tuning. A hundredth of parameter change could be
	// the difference between a stable simulation and absolute chaos.
	SystemScale float64

	// Parameters
	*StaticParams
	*DynamicParams

	// Kernel normalization consts
	poly6NormalizationConst     float64
	spikyNormalizationConst     float64
	viscosityNormalizationConst float64
}

func NewSystemConfig(
	width int,
	height int,
	particles int,
	particleUiRadius float64,
	systemScale float64,
	smoothingRadiusH float64,
	timestep float64,
	particleMass float64,
	gravityForce float64,
	gasConstK float64,
	restDensity float64,
	viscosityConst float64,
	velocityLimit float64,
	collisionDamping float64,
) *SystemConfig {
	staticParams := &StaticParams{
		SmoothingRadiusH: smoothingRadiusH,
		Timestep:         timestep,
	}

	dynamicParams := &DynamicParams{
		GravityForce:     gravityForce,
		GasConstK:        gasConstK,
		RestDensity:      restDensity,
		ViscosityConst:   viscosityConst,
		VelocityLimit:    velocityLimit,
		CollisionDamping: collisionDamping,
		ParticleMass:     particleMass,
	}

	return &SystemConfig{
		Width:            float64(width) / systemScale,
		Height:           float64(height) / systemScale,
		Particles:        uint(particles),
		ParticleUiRadius: float64(particleUiRadius) / systemScale,
		SystemScale:      systemScale,
		StaticParams:     staticParams,
		DynamicParams:    dynamicParams,

		// Normalization constant for the Poly6 kernel adapted for 2D SPH
		poly6NormalizationConst: 4 / (math.Pi * math.Pow(smoothingRadiusH, 8)),

		// Normalization constant for the Spiky kernel adapted for 2D SPH
		spikyNormalizationConst: -30 / (math.Pi * math.Pow(smoothingRadiusH, 5)),

		// Normalization constant for the Viscosity kernel adapted for 2D SPH
		viscosityNormalizationConst: 40 / (math.Pi * math.Pow(smoothingRadiusH, 5)),
	}
}
