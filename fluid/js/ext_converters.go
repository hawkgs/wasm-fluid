package js

import (
	"github.com/hawkgs/wasm-fluid/fluid/system"
	"github.com/hawkgs/wasm-fluid/fluid/ui"
	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

func convertVectorToMap(v *vectors.Vector) map[string]any {
	m := make(map[string]any)
	m["x"] = v.X * system.SystemScale
	m["y"] = v.Y * system.SystemScale

	return m
}

func convertParticlesToJsArray(particles []*system.Particle) []any {
	mapped := make([]any, len(particles))

	for i := range mapped {
		p := particles[i]
		particleData := convertVectorToMap(p.GetPosition())

		v := p.GetVelocity().Magnitude()
		particleData["vColor"] = ui.GetParticleVelocityColor(v)

		mapped[i] = particleData
	}

	return mapped
}
