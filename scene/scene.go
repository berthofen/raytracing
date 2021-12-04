package scene

import (
	. "raytracing/common"
	vec "raytracing/vector"
)

const (
	MinDepth = 0
)

type Scene struct {
	AmbCol Color
	AmbInt float64
	Lights []LightSource
	Objects []RayIntersector
	MaxDepth uint
}

type LightSource struct {
	Pos vec.Vector
	Col Color
	Intens float64
}

type RayIntersector interface {
	Intersect(a Ray) (*vec.Vector, float64, *Material, *vec.Vector, *vec.Vector)
}

type Ray struct {
	From vec.Vector
	Dir vec.Vector
}