package js

import (
	"github.com/hawkgs/wasm-fluid/fluid/system"
	"github.com/hawkgs/wasm-fluid/fluid/ui"
	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

// convertVectorToMap converts a Vector to a JS-compatible map
func convertVectorToMap(v *vectors.Vector, scalingFactor float64) map[string]any {
	m := make(map[string]any)
	m["x"] = v.X * scalingFactor
	m["y"] = v.Y * scalingFactor

	return m
}

// convertParticlesToJsArray converts an array of particles to an array of
// JS-compatible objects that contain the position and velocity color
func convertParticlesToJsArray(particles []*system.Particle, scalingFactor float64) []any {
	mapped := make([]any, len(particles))

	for i := range mapped {
		p := particles[i]
		particleData := convertVectorToMap(p.GetPosition(), scalingFactor)

		v := p.GetVelocity().Magnitude()
		particleData["vColor"] = ui.GetParticleVelocityColor(v)

		mapped[i] = particleData
	}

	return mapped
}
