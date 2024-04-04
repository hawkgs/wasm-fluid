package vectors

type Vector struct {
	X float64
	Y float64
}

func NewVector(x float64, y float64) *Vector {
	return &Vector{x, y}
}
