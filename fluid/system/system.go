package system

import (
	"math/rand"

	"github.com/hawkgs/wasm-fluid/fluid/forces"
	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

type System struct {
	config    *SystemConfig
	particles []*Particle
	forces    []forces.Force
}

func NewSystem(cfg *SystemConfig, forces []forces.Force) *System {
	particles := createParticles(cfg)

	return &System{cfg, particles, forces}
}

func (s *System) Update() []*Particle {
	for _, particle := range s.particles {
		s.applyForces(particle)
		particle.Contain()
		particle.Update()
	}

	return s.particles
}

func createParticles(cfg *SystemConfig) []*Particle {
	particles := make([]*Particle, cfg.Particles)
	container := vectors.NewVector(float64(cfg.Width), float64(cfg.Height))

	for i := range particles {
		x := rand.Float64() * float64(cfg.Width)
		y := rand.Float64() * float64(cfg.Height)

		position := vectors.NewVector(x, y)

		particles[i] = NewParticle(position, container)
	}

	return particles
}

func (s *System) applyForces(p *Particle) {
	for _, f := range s.forces {
		p.ApplyForce(f.GetVector())
	}
}
