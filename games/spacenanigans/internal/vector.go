package internal

import "math"

type Vector struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (v Vector) Add(w Vector) Vector {
	return Vector{
		X: v.X + w.X,
		Y: v.Y + w.Y,
	}
}

func (v Vector) Addf(f float64) Vector {
	return Vector{
		X: v.X + f,
		Y: v.Y + f,
	}
}

func (v Vector) Sub(w Vector) Vector {
	return Vector{
		X: v.X - w.X,
		Y: v.Y - w.Y,
	}
}

func (v Vector) Subf(f float64) Vector {
	return Vector{
		X: v.X - f,
		Y: v.Y - f,
	}
}

func (v Vector) Mul(w Vector) Vector {
	return Vector{
		X: v.X * w.X,
		Y: v.Y * w.Y,
	}
}

func (v Vector) Mulf(f float64) Vector {
	return Vector{
		X: v.X * f,
		Y: v.Y * f,
	}
}

func (v Vector) Div(w Vector) Vector {
	return Vector{
		X: v.X / w.X,
		Y: v.Y / w.Y,
	}
}

func (v Vector) Divf(f float64) Vector {
	return Vector{
		X: v.X / f,
		Y: v.Y / f,
	}
}

func (v Vector) Round() Vector {
	return Vector{
		X: math.Round(v.X),
		Y: math.Round(v.Y),
	}
}

////////////////////////////////////////////////////////////////////////////////

func (v Vector) Zero() bool {
	return v.X == 0 && v.Y == 0
}

func (v Vector) Eq(w Vector) bool {
	return v.X == w.X && v.Y == w.Y
}

func (v Vector) Distance(w Vector) float64 {
	// Manhattan distance
	//return math.Abs(w.X-v.X) + math.Abs(w.Y-v.Y)
	// Euclidean distance
	dx, dy := w.X-v.X, w.Y-v.Y
	return math.Sqrt(dx*dx + dy*dy)
}
