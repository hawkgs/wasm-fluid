package system

import (
	"fmt"
	"math"

	"github.com/hawkgs/wasm-fluid/fluid/utils"
	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

// System dev/debugging functions

// devAlarmForNanPos alerts for particles with NaN position
func (s *System) devAlarmForNanPos(p *Particle) {
	if !s.devNanDetected && (math.IsNaN(p.position.X) || math.IsNaN(p.position.Y)) {
		fmt.Println("NaN position detected!")
		fmt.Println("Position =", p.position)
		fmt.Println("Velocity =", p.velocity)
		fmt.Println("Velocity 1/2 =", p.velocityHalf)
		fmt.Println("Density =", p.density)
		fmt.Println("")

		s.devNanDetected = true
	}
}

// DevPrintStats prints the current status of the system. Used for debugging.
func (s *System) DevPrintStats() {
	var (
		nanParticles      uint    = 0
		infParticles      uint    = 0
		topLeftCorner     uint    = 0
		topRightCorner    uint    = 0
		bottomLeftCorner  uint    = 0
		bottomRightCorner uint    = 0
		totalDensity      float64 = 0
		minDensity        float64 = math.Inf(1)
		maxDensity        float64 = math.Inf(-1)
	)

	cfg := s.config
	uiRad := s.config.ParticleUiRadius
	cont := vectors.NewVector(
		cfg.Width-uiRad,
		cfg.Height-uiRad,
	)

	for _, p := range s.particles {
		pos := p.position
		totalDensity += p.density

		if math.IsNaN(pos.X) || math.IsNaN(pos.Y) {
			nanParticles++
		} else if math.IsInf(pos.X, 1) || math.IsInf(pos.Y, 1) || math.IsInf(pos.X, -1) || math.IsInf(p.position.Y, -1) {
			infParticles++
		}

		if pos.X == uiRad && pos.Y == uiRad {
			topLeftCorner++
		}
		if pos.X == cont.X && pos.Y == uiRad {
			topRightCorner++
		}
		if pos.X == uiRad && pos.Y == cont.Y {
			bottomLeftCorner++
		}
		if pos.X == cont.X && pos.Y == cont.Y {
			bottomRightCorner++
		}

		if p.density > maxDensity {
			maxDensity = p.density
		}
		if p.density < minDensity {
			minDensity = p.density
		}
	}

	okayParticles := cfg.Particles - nanParticles - infParticles
	atCorner := topLeftCorner + topRightCorner + bottomLeftCorner + bottomRightCorner
	avgDensity := totalDensity / float64(len(s.particles))

	fmt.Println("")
	fmt.Println("SYSTEM STATS SNAPSHOT")
	fmt.Println("*********************")
	fmt.Println("Current frame:", s.devFramesCt)
	fmt.Println(
		"Params: Field =", cfg.Width,
		"x", cfg.Height,
		"| h =", cfg.SmoothingRadiusH,
		"| p⌀ =", cfg.ParticleUiRadius*2,
		"| 𝚫t =", cfg.Timestep,
		"| m =", cfg.ParticleMass,
	)
	fmt.Println(
		"G =", cfg.GravityForce,
		"| k =", cfg.GasConstK,
		"| ⍴₀ =", cfg.RestDensity,
		"| μ =", cfg.ViscosityConst,
		"| V limit =", cfg.VelocityLimit,
		"| Col. damp =", cfg.CollisionDamping,
	)
	fmt.Println("Particles => OK:", okayParticles, "/ NaN:", nanParticles, "/ Inf:", infParticles)
	fmt.Println("Particles at the corners:", atCorner)
	fmt.Println("⌜", topLeftCorner, "  ⌝", topRightCorner)
	fmt.Println("⌞", bottomLeftCorner, "  ⌟", bottomRightCorner)
	fmt.Println(
		"Density => Avg:", utils.RoundNum(avgDensity, 7),
		"/ Max:", utils.RoundNum(maxDensity, 7),
		"/ Min:", utils.RoundNum(minDensity, 7))
	fmt.Println("")
}
