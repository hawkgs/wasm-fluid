package system

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/hawkgs/wasm-fluid/fluid/forces"
	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

type System struct {
	config         *SystemConfig
	particles      []*Particle
	grid           map[string][]*Particle
	externalForces []forces.Force
	gridWidth      uint
	gridHeight     uint
}

func NewSystem(cfg *SystemConfig, extForces []forces.Force) *System {
	particles := createParticles(cfg)

	gridWidth := uint(math.Ceil(float64(cfg.Width) / float64(smoothingRadiusH)))
	gridHeight := uint(math.Ceil(float64(cfg.Height) / float64(smoothingRadiusH)))

	return &System{
		cfg,
		particles,
		make(map[string][]*Particle, gridWidth*gridHeight),
		extForces,
		gridWidth,
		gridHeight,
	}
}

func (s *System) Update() []*Particle {
	s.updateGrid()

	for _, p := range s.particles {
		density := calculateDensityGradient(s, p)
		p.SetDensity(density)
	}

	pressures := []float64{}
	for _, p := range s.particles {
		pressures = append(pressures, calculatePressureGradient(s, p))
	}

	for _, particle := range s.particles {
		s.applyForces(particle)
		// pres := pressures[i]
		// particle.ApplyForce(vectors.NewVector(pres, pres))

		particle.Contain()
		particle.Update()
	}

	return s.particles
}

func (s *System) applyForces(p *Particle) {
	for _, f := range s.externalForces {
		p.ApplyForce(f.GetVector())
	}
}

func (s *System) updateGrid() {
	cellsCount := s.gridWidth * s.gridHeight
	grid := make(map[string][]*Particle, cellsCount)

	for _, p := range s.particles {
		key := s.getParticleCellKey(s.getParticleCell(p))

		if cell := grid[key]; cell != nil {
			cell = append(cell, p)
			grid[key] = cell
		} else {
			grid[key] = []*Particle{p}
		}
	}

	s.grid = grid
}

func (s *System) getParticleCell(p *Particle) [2]int {
	x := uint(p.position.X) / s.gridWidth
	y := uint(p.position.Y) / s.gridHeight

	return [2]int{int(x), int(y)}
}

func (s *System) getParticleCellKey(cell [2]int) string {
	return strconv.Itoa(cell[0]) + "," + strconv.Itoa(cell[1])
}

// Includes the target particle as well
func (s *System) getParticleNeighbors(p *Particle) []*Particle {
	cell := s.getParticleCell(p)

	cells := [9][2]int{
		{cell[0] - 1, cell[1] - 1}, // top left
		{cell[0] - 1, cell[1]},     // top
		{cell[0] - 1, cell[1] + 1}, // top right
		{cell[0], cell[1] - 1},     // left
		{cell[0], cell[1] + 1},     // right
		{cell[0] + 1, cell[1]},     // bottom
		{cell[0] + 1, cell[1] - 1}, // bottom left
		{cell[0] + 1, cell[1] + 1}, // bottom right
	}

	particles := []*Particle{}

	for _, c := range cells {
		key := s.getParticleCellKey(c)

		if c := s.grid[key]; c != nil {
			particles = append(particles, c...)
		}
	}

	return particles
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
