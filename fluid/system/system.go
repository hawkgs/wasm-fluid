package system

import (
	"math"
	"math/rand"

	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

type System struct {
	config    *SystemConfig
	particles []*vectors.Vector
}

func NewSystem(cfg *SystemConfig) *System {
	return &System{cfg, []*vectors.Vector{}}
}

func (s *System) Update() []*vectors.Vector {
	// Note(Georgi): Temp. Testing
	x := math.Min(rand.Float64()*100, float64(s.config.Width))
	x = math.Max(0, x)

	y := math.Min(rand.Float64()*100, float64(s.config.Height))
	y = math.Max(0, y)

	v := vectors.NewVector(x, y)
	s.particles = nil
	s.particles = append(s.particles, v)

	return s.particles
}
