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
	gravityForce     = 750
	smoothingRadiusH = 0.5
	gasConstK        = 100.0
	restDensity      = 10.0
	viscosityConst   = 50.0
	collisionDamping = 0.5  // The higher, the more damping
	timestep         = 0.02 // i.e. DT; Delta Time
)
