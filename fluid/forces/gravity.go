package forces

import "github.com/hawkgs/wasm-fluid/fluid/vectors"

const gravitationalPull = 0.01

type Gravity struct {
	BaseForce
}

func NewGravity() *Gravity {
	f := BaseForce{vectors.NewVector(0, gravitationalPull)}

	return &Gravity{f}
}

func (g *Gravity) GetVector() *vectors.Vector {
	return g.vector
}
