package system

import (
	"fmt"
	"math"
	"strconv"

	"github.com/hawkgs/wasm-fluid/fluid/utils"
	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

type System struct {
	config     *SystemConfig
	particles  []*Particle
	grid       map[string][]*Particle
	gridWidth  uint
	gridHeight uint

	// For debugging
	devFramesCt    uint
	devNanDetected bool
}

func NewSystem(cfg *SystemConfig) *System {
	particles := createParticles(cfg)

	gridWidth := uint(math.Ceil(cfg.Width / smoothingRadiusH))
	gridHeight := uint(math.Ceil(cfg.Height / smoothingRadiusH))

	return &System{
		cfg,
		particles,
		make(map[string][]*Particle, gridWidth*gridHeight),
		gridWidth,
		gridHeight,
		0,
		false,
	}
}

// Update calculates all densities, the Navier-Stokes forces based on SPH and then applies them
func (s *System) Update() []*Particle {
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

		p.ApplyForce(nsForces)
		p.Update()
		p.Contain()

		s.devAlarmForNanPos(p)
	}

	s.devFramesCt++

	return s.particles
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

// devAlarmForNanPos alerts for particles with NaN position
func (s *System) devAlarmForNanPos(p *Particle) {
	if !s.devNanDetected && (math.IsNaN(p.position.X) || math.IsNaN(p.position.Y)) {
		fmt.Println("NaN position detected!")
		fmt.Println("Position =", p.position)
		fmt.Println("Acceleration =", p.acceleration)
		fmt.Println("Velocity =", p.velocity)
		fmt.Println("Density =", p.density)
		fmt.Println("")

		s.devNanDetected = true
	}
}

// DevPrintStats prints the current status of the system. Used for debugging.
func (s *System) DevPrintStats() {
	var nanParticles uint = 0
	var infParticles uint = 0
	var cornerParticles uint = 0
	var totalDensity float64 = 0
	var minDensity float64 = math.Inf(1)
	var maxDensity float64 = math.Inf(-1)

	uiRad := s.config.ParticleUiRadius
	cont := vectors.NewVector(
		s.config.Width-uiRad,
		s.config.Height-uiRad,
	)

	for _, p := range s.particles {
		pos := p.position
		totalDensity += p.density

		if math.IsNaN(pos.X) || math.IsNaN(pos.Y) {
			nanParticles++
		} else if math.IsInf(pos.X, 1) || math.IsInf(pos.Y, 1) || math.IsInf(pos.X, -1) || math.IsInf(p.position.Y, -1) {
			infParticles++
		}

		if pos.X == uiRad && pos.Y == uiRad ||
			pos.X == cont.X && pos.Y == uiRad ||
			pos.X == uiRad && pos.Y == cont.Y ||
			pos.X == cont.X && pos.Y == cont.Y {
			cornerParticles++
		}

		if p.density > maxDensity {
			maxDensity = p.density
		}
		if p.density < minDensity {
			minDensity = p.density
		}
	}

	fmt.Println("")
	fmt.Println("CURRENT SYSTEM STATS")
	fmt.Println("********************")
	fmt.Println("Current frame:", s.devFramesCt)
	fmt.Println("")
	fmt.Println("Okay particles:", s.config.Particles-nanParticles-infParticles)
	fmt.Println("NaN particles:", nanParticles)
	fmt.Println("Inf particles:", infParticles)
	fmt.Println("")
	fmt.Println("Particles at the corners:", cornerParticles)
	fmt.Println("")
	fmt.Println("Average density:", totalDensity/float64(len(s.particles)))
	fmt.Println("Max density:", maxDensity)
	fmt.Println("Min density:", minDensity)
	fmt.Println("....................")
}
