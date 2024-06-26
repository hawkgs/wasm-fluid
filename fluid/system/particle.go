package system

import (
	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

type Particle struct {
	velocity     *vectors.Vector
	velocityHalf *vectors.Vector
	position     *vectors.Vector
	container    *vectors.Vector
	density      float64
	cfg          *SystemConfig
}

func NewParticle(position *vectors.Vector, container *vectors.Vector, cfg *SystemConfig) *Particle {
	return &Particle{
		velocity:     vectors.NewVector(0, 0),
		velocityHalf: vectors.NewVector(0, 0),
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

func (p *Particle) SetDensity(density float64) {
	p.density = density
}

func (p *Particle) getDensity() float64 {
	// Since our simulation is still unstable, we have cases
	// where a single particle could be outside of the smoothing radius of
	// any other particle, so the density is 0. However, ext. forces are
	// still applicable, so in order to avoid division by zero,
	// we use the mass instead.
	// Why density is used?: Müller et al – Eqn. (8)
	density := p.density
	if density <= 0 {
		density = p.cfg.ParticleMass
	}
	return density
}

// For integration, we use Leapfrog method since Navier-Stokes PDE is of 2nd order

// ApplyForce updates vector's velocity and position based on the provided force
func (p *Particle) ApplyForce(force *vectors.Vector) {
	// Newton's 2nd law: Acceleration = Sum of all forces / Mass (or density in our case)
	acceleration := force.ImmutDivide(p.getDensity())

	p.velocityHalf.Add(acceleration.ImmutMultiply(p.cfg.Timestep))
	p.velocityHalf.Limit(p.cfg.VelocityLimit)

	p.velocity.Add(p.velocityHalf.ImmutAdd(acceleration.ImmutMultiply(p.cfg.Timestep / 2)))
	p.velocity.Limit(p.cfg.VelocityLimit)

	p.position.Add(p.velocityHalf.ImmutMultiply(p.cfg.Timestep))
	p.contain()
}

// ApplyInitialForce should be called the first time the particle is updated in order to
// update the velocity at half step (Leapfrog)
func (p *Particle) ApplyInitialForce(force *vectors.Vector) {
	acceleration := force.ImmutDivide(p.getDensity())

	p.velocityHalf = acceleration.ImmutMultiply(p.cfg.Timestep / 2)
	p.velocity = acceleration.ImmutMultiply(p.cfg.Timestep)
	p.position.Add(p.velocityHalf.ImmutMultiply(p.cfg.Timestep))

	p.contain()
}

// contain keeps the particle within its container (bounces off) when it reaches an edge
func (p *Particle) contain() {
	reflect := -1 * (1 - p.cfg.CollisionDamping)

	// Right/left
	if p.position.X > p.container.X {
		p.velocity.X *= reflect
		p.velocityHalf.X *= reflect
		p.position.X = p.container.X
	} else if p.position.X < p.cfg.ParticleUiRadius {
		p.velocity.X *= reflect
		p.velocityHalf.X *= reflect
		p.position.X = p.cfg.ParticleUiRadius
	}

	// Bottom/top
	if p.position.Y > p.container.Y {
		p.velocity.Y *= reflect
		p.velocityHalf.Y *= reflect
		p.position.Y = p.container.Y
	} else if p.position.Y < p.cfg.ParticleUiRadius {
		p.velocity.Y *= reflect
		p.velocityHalf.Y *= reflect
		p.position.Y = p.cfg.ParticleUiRadius
	}
}
