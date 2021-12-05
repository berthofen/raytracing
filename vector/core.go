package vector

import (
	"math"
)

type Vector struct {
	X, Y, Z float64
}

func (a Vector) Add(b Vector) Vector {
	a.X += b.X
	a.Y += b.Y
	a.Z += b.Z
	return a
}

func (a Vector) Sub(b Vector) Vector {
	return Vector{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func (a Vector) MultiplyByScalar(s float64) Vector {
	return Vector{
		X: a.X * s,
		Y: a.Y * s,
		Z: a.Z * s,
	}
}

func (a Vector) Dot(b Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vector) Length() float64 {
	return math.Sqrt(a.Dot(a))
}

func (a Vector) Cross(b Vector) Vector {
	return Vector{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

func (a Vector) Normalize() Vector {
	return a.MultiplyByScalar(1. / a.Length())
}

func (a Vector) Inside() Vector {
	return a.MultiplyByScalar(1. / a.Length())
}

func (normal Vector) Reflect(a Vector) Vector {
	return normal.MultiplyByScalar(2.*a.Dot(normal)).Sub(a)
}