package system

import (
	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

type Particle struct {
	acceleration *vectors.Vector
	velocity     *vectors.Vector
	position     *vectors.Vector
	container    *vectors.Vector
	density      float64
	cfg          *SystemConfig
}

func NewParticle(position *vectors.Vector, container *vectors.Vector, cfg *SystemConfig) *Particle {
	return &Particle{
		acceleration: vectors.NewVector(0, 0),
		velocity:     vectors.NewVector(0, 0),
		position:     position,
		container:    container,
		cfg:          cfg,
	}
}

func (p *Particle) GetPosition() *vectors.Vector {
	return p.position
}

func (p *Particle) GetVelocity() *vectors.Vector {
	return p.velocity
}

func (p *Particle) GetAcceleration() *vectors.Vector {
	return p.acceleration
}

func (p *Particle) SetDensity(density float64) {
	p.density = density
}

// ApplyForce adds the force vector the object's acceleration vector
func (p *Particle) ApplyForce(force *vectors.Vector) {
	// Since our simulation is still unstable, we have cases
	// where a single particle could be outside of the smoothing radius of
	// any other particle, so the density is 0. In order to avoid NaN use the
	// mass in that case instead
	// Why density is used?: Müller et al – Eqn. (8)
	density := p.density
	if density <= 0 {
		density = particleMass
	}

	// Newton's 2nd law: Acceleration = Sum of all forces / Mass (or density in our case)
	force.Divide(density)
	force.Multiply(timestep) // Euler method (integration; Muller's SPH assumes leap frog)

	p.acceleration.Add(force)
}

// Update modifies the object's position depending on the applied forces on each rendering iteration
func (p *Particle) Update() {
	// We keep the velocity only for correctness based on physics laws
	p.velocity.Add(p.acceleration)
	p.velocity.Multiply(timestep) // Euler method (integration; Muller's SPH assumes leap frog)
	p.position.Add(p.velocity)

	// Limit the velocity
	p.velocity.Limit(velocityLimit)

	// Clear the acceleration
	p.acceleration.Multiply(0)
}

// Contain keeps the particle within its container (bounces off) when it reaches an edge
func (p *Particle) Contain() {
	// Right/left
	if p.position.X > p.container.X {
		p.velocity.X *= -1 * collisionDamping
		p.position.X = p.container.X
	} else if p.position.X < p.cfg.ParticleUiRadius {
		p.velocity.X *= -1 * collisionDamping
		p.position.X = p.cfg.ParticleUiRadius
	}

	// Bottom/top
	if p.position.Y > p.container.Y {
		p.velocity.Y *= -1 * collisionDamping
		p.position.Y = p.container.Y
	} else if p.position.Y < p.cfg.ParticleUiRadius {
		p.velocity.Y *= -1 * collisionDamping
		p.position.Y = p.cfg.ParticleUiRadius
	}
}
