package forces

import "github.com/hawkgs/wasm-fluid/fluid/vectors"

const windForce = 0.0001

type Wind struct {
	BaseForce
}

func NewWind() *Wind {
	f := BaseForce{vectors.NewVector(windForce, 0)}

	return &Wind{f}
}

func (g *Wind) GetVector() *vectors.Vector {
	return g.vector
}
