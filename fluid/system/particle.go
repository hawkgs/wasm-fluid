package system

import (
	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

const (
	particleMass  = 1
	velocityLimit = 25
)

type Particle struct {
	acceleration *vectors.Vector
	velocity     *vectors.Vector
	location     *vectors.Vector
	container    *vectors.Vector
	mass         float64
}

func NewParticle(location *vectors.Vector, container *vectors.Vector) *Particle {
	return &Particle{
		acceleration: vectors.NewVector(0, 0),
		velocity:     vectors.NewVector(0, 0),
		location:     location,
		container:    container,
		mass:         particleMass,
	}
}

func (p *Particle) GetLocation() *vectors.Vector {
	return p.location
}

// ApplyForce adds the force vector the object's acceleration vector
func (p *Particle) ApplyForce(force *vectors.Vector) {
	// Newton's 2nd law: Acceleration = Sum of all forces / Mass
	fCopy := force.Copy()
	fCopy.Divide(p.mass)
	p.acceleration.Add(fCopy)
}

// Update modifies the object's location depending on the applied forces;
// Should be called on every rendering iteration
func (p *Particle) Update() {
	// We keep the velocity only for correctness based on physics laws
	p.velocity.Add(p.acceleration)
	p.location.Add(p.velocity)

	// Limit the velocity
	p.velocity.Limit(velocityLimit)

	// Clear the acceleration
	p.acceleration.Multiply(0)
}

// Contain keeps the mover within its container (bounces off) when it reaches an edge
func (p *Particle) Contain() {
	// Right/left
	if p.location.X > p.container.X {
		p.velocity.X *= -1
		p.location.X = p.container.X
	} else if p.location.X < 0 {
		p.velocity.X *= -1
		p.location.X = 0
	}

	// Bottom/top
	if p.location.Y > p.container.Y {
		p.velocity.Y *= -1
		p.location.Y = p.container.Y
	} else if p.location.Y < 0 {
		p.velocity.Y *= -1
		p.location.Y = 0
	}
}
