package system

const (
	// The SystemScale downsizes the particle field.
	// For example, 800px-wide canvas will be downscaled to 20, if the SystemScale is 40.
	// This is needed because the kernel functions are very sensitive to the smoothing radius.
	// They are actually meant to work with specific smoothing radiuses. You can easily check
	// that by plotting the functions (or their derivates, resp.) and examining how they change over time.
	// SPH is very sensitive to parameter tuning. A hundredth of parameter change could be
	// the difference between a stable simulation and absolute chaos.
	SystemScale      = 40.0
	particleMass     = 1.0
	velocityLimit    = 10.0
	gravityForce     = 0
	smoothingRadiusH = 0.45
	gasConstK        = 500.0
	restDensity      = 2.0
	viscosityConst   = 0.0
	collisionDamping = 0.3   // The smaller, the more damping
	timestep         = 0.003 // i.e. DT; Delta Time
)
