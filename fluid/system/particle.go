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

// ApplyForce adds the force vector the object's acceleration vector
func (p *Particle) ApplyForce(force *vectors.Vector) {
	// Newton's 2nd law: Acceleration = Sum of all forces / Mass
	fCopy := force.Copy()
	fCopy.Divide(particleMass) // Change mass to density based on Eqn. (8)

	p.acceleration.Add(fCopy)
}

// Update modifies the object's position depending on the applied forces;
// Should be called on every rendering iteration
func (p *Particle) Update() {
	// We keep the velocity only for correctness based on physics laws
	p.velocity.Add(p.acceleration)
	p.position.Add(p.velocity)

	// Limit the velocity
	p.velocity.Limit(velocityLimit)

	// Clear the acceleration
	p.acceleration.Multiply(0)
}

// Contain keeps the mover within its container (bounces off) when it reaches an edge
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

func (p *Particle) SetDensity(density float64) {
	p.density = density
}
