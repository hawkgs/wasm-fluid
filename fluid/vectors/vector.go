package vectors

type Vector struct {
	X float64
	Y float64
}

func NewVector(x float64, y float64) *Vector {
	return &Vector{x, y}
}

func (v *Vector) ToMap() map[string]any {
	m := make(map[string]any)
	m["x"] = v.X
	m["y"] = v.Y

	return m
}
