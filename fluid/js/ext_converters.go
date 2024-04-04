package js

import "github.com/hawkgs/wasm-fluid/fluid/vectors"

func convertVectorToMap(v *vectors.Vector) map[string]any {
	m := make(map[string]any)
	m["x"] = v.X
	m["y"] = v.Y

	return m
}

func convertVectorsToArray(v []*vectors.Vector) []any {
	mapped := make([]any, len(v))

	for i := range mapped {
		mapped[i] = convertVectorToMap(v[i])
	}

	return mapped
}
