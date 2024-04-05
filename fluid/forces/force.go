package forces

import "github.com/hawkgs/wasm-fluid/fluid/vectors"

type Force interface {
	GetVector() *vectors.Vector
}

type BaseForce struct {
	vector *vectors.Vector
}
