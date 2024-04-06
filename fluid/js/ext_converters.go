package js

import (
	"github.com/hawkgs/wasm-fluid/fluid/system"
	"github.com/hawkgs/wasm-fluid/fluid/utils"
	"github.com/hawkgs/wasm-fluid/fluid/vectors"
)

func convertVectorToMap(v *vectors.Vector) map[string]any {
	m := make(map[string]any)
	m["x"] = v.X
	m["y"] = v.Y

	return m
}

func convertParticlesToLocationsArray(particles []*system.Particle) []any {
	mapped := make([]any, len(particles))

	for i := range mapped {
		location := convertVectorToMap(particles[i].GetLocation())
		location["x"] = utils.Round(location["x"].(float64), 6)
		location["y"] = utils.Round(location["y"].(float64), 6)

		mapped[i] = location
	}

	return mapped
}
