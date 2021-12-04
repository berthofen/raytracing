package common

import "math"

type Material struct {
	Col   Color
	Ka    float64 // ambient reflection constant
	Kd    float64 // diffusive reflection constant
	Ks    float64 // specular reflection constant
	Kr    float64 // reflection constant
	Alpha float64 // shininess
}

type Color struct {
	R uint
	G uint
	B uint
}

func (a Color) Add(b Color) Color {
	return Color{min(255, a.R+b.R),
		min(255, a.G+b.G),
		min(255, a.B+b.B)}
}

func (a Color) Scale(b float64) Color {
	return Color{uint(math.Max(0., math.Min(255., math.Round(float64(a.R)*b)))),
		uint(math.Max(0., math.Min(255., math.Round(float64(a.G)*b)))),
		uint(math.Max(0., math.Min(255., math.Round(float64(a.B)*b))))}
}

func min(a uint, b uint) uint {
	if a > b {
		return b
	} else {
		return a
	}
}

func max(a uint, b uint) uint {
	if a > b {
		return a
	} else {
		return b
	}
}
