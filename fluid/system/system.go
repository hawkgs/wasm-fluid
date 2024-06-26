package system

import (
	"math"
	"strconv"

	"github.com/hawkgs/wasm-fluid/fluid/utils"
	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

type System struct {
	config      *SystemConfig
	particles   []*Particle
	grid        map[string][]*Particle
	gridWidth   uint
	gridHeight  uint
	firstUpdate bool

	// For debugging
	devFramesCt    uint
	devNanDetected bool
}

func NewSystem(cfg *SystemConfig) *System {
	particles := createParticles(cfg)

	gridWidth := uint(math.Ceil(cfg.Width / cfg.SmoothingRadiusH))
	gridHeight := uint(math.Ceil(cfg.Height / cfg.SmoothingRadiusH))

	return &System{
		cfg,
		particles,
		make(map[string][]*Particle, gridWidth*gridHeight),
		gridWidth,
		gridHeight,
		true,
		0,
		false,
	}
}

func (s *System) GetConfig() *SystemConfig {
	return s.config
}

// UpdateDynamicParams updates all dynamic simulation params
func (s *System) UpdateDynamicParams(params *DynamicParams) {
	s.config.DynamicParams = params
}

// Update updates the system
func (s *System) Update() []*Particle {
	if s.firstUpdate {
		s.calculateForces(func(p *Particle, f *vectors.Vector) {
			p.ApplyInitialForce(f)
		})
		s.firstUpdate = false
	}

	s.calculateForces(func(p *Particle, f *vectors.Vector) {
		p.ApplyForce(f)
		s.devAlarmForNanPos(p)
	})

	s.devFramesCt++

	return s.particles
}

// calculateForces calculates all densities, the Navier-Stokes forces based on SPH and then applies them
func (s *System) calculateForces(op func(p *Particle, f *vectors.Vector)) {
	// We need to know the cell of every particle prior the calculations
	s.updateGrid()

	// Calculate densities
	for _, p := range s.particles {
		density := calculateDensity(s, p)
		p.SetDensity(density)
	}

	// Calculates pressure, viscosity and ext. forces for each particle and then apply them
	for _, p := range s.particles {
		nsForces := calculateNavierStokesForces(s, p)
		op(p, nsForces)
	}
}

// updateGrid creates/updates a grid with particles with the size of the smoothing radius.
// Since SPH checks the particles only within the smoothing radius, this optimization helps
// us to reduce the number of calculation (i.e. not iterate over particles outside the radius).
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

// getParticleCell returns the cell indices of a particle
func (s *System) getParticleCell(p *Particle) [2]int {
	percX := p.position.X / s.config.Width
	percY := p.position.Y / s.config.Height

	x := uint(math.Floor(percX * float64(s.gridWidth)))
	y := uint(math.Floor(percY * float64(s.gridHeight)))

	return [2]int{int(y), int(x)}
}

// getParticleCellKey returns the grid key of a particle
func (s *System) getParticleCellKey(cell [2]int) string {
	return strconv.Itoa(cell[0]) + "," + strconv.Itoa(cell[1])
}

// getParticleNeighbors returns all neighbors of a given particle that
// might be within the smoothing radius (i.e. the particles inside the
// neighbor grid cells; excl. the target particle).
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

	// Add the target cell particles without the target particle
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

// createParticles creates the initial stack of particles that is dropped onto the field
func createParticles(cfg *SystemConfig) []*Particle {
	particles := make([]*Particle, cfg.Particles)
	container := vectors.NewVector(
		cfg.Width-cfg.ParticleUiRadius,
		cfg.Height-cfg.ParticleUiRadius,
	)

	margin := cfg.ParticleUiRadius * 2

	// The hardcoded values are arbitrary and affect
	// the initial position of the particle stack
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
