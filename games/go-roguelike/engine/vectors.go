package engine

import (
	"fmt"
	"math"
)

type Vector struct {
	X, Y float64
}

func NewVector(x, y float64) *Vector {
	return &Vector{x, y}
}

func (v *Vector) Str() string {
	return fmt.Sprintf("<%.1f, %.1f>", v.X, v.Y)
}

func (v *Vector) Clone() *Vector {
	return NewVector(v.X, v.Y)
}

func (v *Vector) Add(d *Vector) *Vector {
	return NewVector(v.X+d.X, v.Y+d.Y)
}

func (v *Vector) Sub(d *Vector) *Vector {
	return NewVector(v.X-d.X, v.Y-d.Y)
}

func (v *Vector) Mul(d *Vector) *Vector {
	return NewVector(v.X*d.X, v.Y*d.Y)
}

func (v *Vector) Div(d *Vector) *Vector {
	return NewVector(v.X/d.X, v.Y/d.Y)
}

func (v *Vector) Eq(d *Vector) bool {
	//Equal
	return (v.X == d.X && v.Y == d.Y)
}

func (v *Vector) Lt(d *Vector) bool {
	// Less than
	return (v.X < d.X) || (v.X == d.X && v.Y < d.Y)
}

func (v *Vector) Le(d *Vector) bool {
	// Less than or equal
	return (v.X <= d.X && v.Y <= d.Y)
}

func (v *Vector) Len() float64 {
	// Length
	return math.Sqrt(v.X + v.X + v.Y*v.Y)
}

func (v *Vector) Dist(d *Vector) float64 {
	// Distance
	dx := v.X - d.X
	dy := v.Y - d.Y
	return math.Sqrt(dx*dx + dy*dy)
}
