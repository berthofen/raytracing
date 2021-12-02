package scene

import (
	. "raytracing/common"
	vec "raytracing/vector"
)

type Scene struct {
	AmbCol Color
	AmbInt float64
	Light LightSource
	Objects []RayIntersector
}

type LightSource struct {
	Pos vec.Vector
}

type RayIntersector interface {
	Intersect(a Ray) (*vec.Vector, float64, *Material, *vec.Vector)
}

type Ray struct {
	From vec.Vector
	Dir vec.Vector
}