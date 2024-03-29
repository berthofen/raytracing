package scene

import (
	. "raytracing/common"
	vec "raytracing/vector"
)

const (
	MinDepth = 0
)

type Scene struct {
	AmbCol   Color
	AmbInt   float64
	Lights   []LightSource
	Objects  []RayIntersector
	MaxDepth uint
}

type LightSource struct {
	Pos    vec.Vector
	Col    Color
	Intens float64
}

type RayIntersector interface {
	Intersect(a Ray) (*vec.Vector, float64, *Material, *vec.Vector, *vec.Vector)
}

type RayIntersection struct {
	Intersection   vec.Vector
	DistanceFactor float64
	Material       Material
	Normal         vec.Vector
	Reflection     vec.Vector
}

type Ray struct {
	From vec.Vector
	Dir  vec.Vector
}

func (a LightSource) Visible(sc Scene, fromPos vec.Vector, fromObj int) bool {
	vecToLight := a.Pos.Sub(fromPos)

	for ind, obj := range sc.Objects {
		if ind != fromObj {
			_, xl, _, _, _ := obj.Intersect(Ray{fromPos, vecToLight})
			if xl > 0 {
				return false
			}
		}
	}
	return true
}
