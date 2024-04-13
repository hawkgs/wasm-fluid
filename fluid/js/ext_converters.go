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

func convertParticlesToPositionsArray(particles []*system.Particle) []any {
	mapped := make([]any, len(particles))

	for i := range mapped {
		position := convertVectorToMap(particles[i].GetPosition())
		position["x"] = utils.RoundNum(position["x"].(float64), 6)
		position["y"] = utils.RoundNum(position["y"].(float64), 6)

		mapped[i] = position
	}

	return mapped
}
