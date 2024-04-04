package system

import (
	"math/rand"

	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

type System struct {
	config    *SystemConfig
	particles []*vectors.Vector
}

func NewSystem(cfg *SystemConfig) *System {
	particles := createParticles(cfg)

	return &System{cfg, particles}
}

func (s *System) Update() []*vectors.Vector {
	for _, p := range s.particles {
		p.X = rand.Float64() * float64(s.config.Width)
		p.Y = rand.Float64() * float64(s.config.Height)
	}

	return s.particles
}

func createParticles(cfg *SystemConfig) []*vectors.Vector {
	particles := make([]*vectors.Vector, cfg.Particles)

	for i := range particles {
		x := rand.Float64() * float64(cfg.Width)
		y := rand.Float64() * float64(cfg.Height)

		particles[i] = vectors.NewVector(x, y)
	}

	return particles
}
