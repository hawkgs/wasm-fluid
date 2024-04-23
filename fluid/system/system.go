package system

import (
	"fmt"
	"math"
	"strconv"

	"github.com/hawkgs/wasm-fluid/fluid/forces"
	"github.com/hawkgs/wasm-fluid/fluid/utils"
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
		density := calculateDensity(s, p)
		p.SetDensity(density)
	}

	nsForces := []*vectors.Vector{}

	for _, p := range s.particles {
		nsForces = append(nsForces, calculateNavierStokesForces(s, p))
	}

	for i, particle := range s.particles {
		particle.ApplyForce(nsForces[i])
		s.applyForces(particle)

		particle.Update()
		particle.Contain()

		if math.IsNaN(particle.position.X) || math.IsNaN(particle.position.Y) {
			fmt.Println("nan position")
		}
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
	percX := p.position.X / s.config.Width
	percY := p.position.Y / s.config.Height

	x := uint(math.Floor(percX * float64(s.gridWidth)))
	y := uint(math.Floor(percY * float64(s.gridHeight)))

	return [2]int{int(y), int(x)}
}

func (s *System) getParticleCellKey(cell [2]int) string {
	return strconv.Itoa(cell[0]) + "," + strconv.Itoa(cell[1])
}

func (s *System) getParticleNeighbors(p *Particle) []*Particle {
	cell := s.getParticleCell(p)

	cells := [8][2]int{
		{cell[0] - 1, cell[1] - 1}, // top left
		{cell[0] - 1, cell[1]},     // top
		{cell[0] - 1, cell[1] + 1}, // top right
		{cell[0], cell[1] - 1},     // left
		{cell[0], cell[1] + 1},     // right
		{cell[0] + 1, cell[1]},     // bottom
		{cell[0] + 1, cell[1] - 1}, // bottom left
		{cell[0] + 1, cell[1] + 1}, // bottom right
	}

	// Add the target cell particles without the target particle P
	particles := s.grid[s.getParticleCellKey(cell)]
	particles = utils.FilterSlice(particles, func(cp *Particle) bool {
		return p != cp
	})

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
	container := vectors.NewVector(
		cfg.Width-cfg.ParticleUiRadius,
		cfg.Height-cfg.ParticleUiRadius,
	)

	margin := spawnedParticleMargin
	height := cfg.Height - margin*4
	cursor := vectors.NewVector(margin*16, margin)

	for i := range particles {
		position := cursor.Copy()
		particles[i] = NewParticle(position, container, cfg)

		if cursor.Y > float64(height) {
			cursor.Y = margin
			cursor.X += margin
		} else {
			cursor.Y += margin
		}
	}

	return particles
}
