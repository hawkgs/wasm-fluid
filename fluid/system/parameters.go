package system

const (
	// The SystemScale downsizes the particle field.
	// For example, 800px-wide canvas will be downscaled to 20, if the SystemScale is 40.
	// This is needed because the kernel functions are very sensitive to the smoothing radius.
	// They are actually meant to work with specific smoothing radiuses. You can easily check
	// that by plotting the functions (or their derivates, resp.) and examining how they change over time.
	// SPH is very sensitive to parameter tuning. A hundredth of parameter change could be
	// the difference between a stable simulation and absolute chaos.
	SystemScale   = 40.0
	particleMass  = 1.0
	velocityLimit = 10.0
	// gravityForce     = 0
	smoothingRadiusH = 0.8
	// gasConstK        = 800.0
	// restDensity      = 1.0
	// viscosityConst   = 0.0
	collisionDamping = 0.1   // The smaller, the more damping
	timestep         = 0.004 // i.e. DT; Delta Time
)

type Parameters struct {
	GravityForce float64
	// smoothingRadiusH float64
	GasConstK      float64
	RestDensity    float64
	ViscosityConst float64
	// Timestep       float64
}
