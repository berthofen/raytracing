package objects

import (
	"fmt"
	"math"
	vec "raytracing/vector"
	. "raytracing/common"
	. "raytracing/scene"
)


type container interface {
	contains(a vec.Vector) bool
}


type Camera struct {
	dir vec.Vector
	dirZ vec.Vector // needs to be orthogonal to dir
	middle vec.Vector
	width float64
	height float64
	resWidth int
	resHeight int
	spectatorDistance float64
	colorChannel int
	pixel []vec.Vector
}

func CameraCreate(dir vec.Vector, dirZ vec.Vector, middle vec.Vector, width float64, height float64, resWidth int, resHeight int, specDist float64, colorChannel int) *Camera {
	c := Camera{
		dir,
		dirZ,
		middle,
		width,
		height,
		resWidth,
		resHeight,
		specDist,
		colorChannel,
		nil}
	c.SetPixels()
	return &c
}

func(c *Camera) SetPixels() {
	c.pixel = make([]vec.Vector, c.resWidth * c.resHeight)
	dirHeight := c.dirZ.Normalize().MultiplyByScalar(c.height / (float64(c.resHeight) - 1.))
	dirWidth  := c.dir.Cross(c.dirZ).Normalize().MultiplyByScalar(c.width / (float64(c.resWidth) - 1.))

	// for even resolution there is no pixel in the middle thats why we add/sub 1/2 of the normalized dir vectors
	start := c.middle.Add(dirHeight.MultiplyByScalar(float64(c.resHeight) / 2.)).Sub(dirWidth.MultiplyByScalar(float64(c.resWidth) / 2.))
	start = start.Sub(dirHeight.MultiplyByScalar(.5)).Add(dirWidth.MultiplyByScalar(.5))

	scalarWidth, scalarHeight := 0., 0.
	for i := 0; i < len(c.pixel); i++ {
		scalarWidth = float64(i%c.resWidth)
		scalarHeight = float64(i/c.resWidth)

		c.pixel[i] = start.Add(dirWidth.MultiplyByScalar(scalarWidth)).Sub(dirHeight.MultiplyByScalar(scalarHeight))
	}

	// fmt.Println(c.pixel)
}

func(c *Camera) Inside(obj container) bool {
	return obj.contains(c.middle)
}





type Sphere struct {
	pos vec.Vector
	rad float64
	mat Material
}

func (s Sphere) contains(a vec.Vector) bool {
	return math.Pow(a.X - s.pos.X, 2) + math.Pow(a.Y - s.pos.Y, 2) + math.Pow(a.Z - s.pos.Z, 2) <= math.Pow(s.rad, 2)
}







func createExampleImage(d []byte, a Color) {
	fmt.Println("creating example image")
	i := 0
	for i < len(d) {
			d[i] = a.R
			i++
			d[i] = a.G
			i++
			d[i] = a.B
			i++
	}
}

func sphereIntersect(a *Camera, pixel vec.Vector, cameraDir vec.Vector, s Sphere) *vec.Vector {
	a1, a2 := 0., 0.

	p := 2. * (cameraDir.X * (pixel.X - s.pos.X) + cameraDir.Y * (pixel.Y - s.pos.Y) + cameraDir.Z * (pixel.Z - s.pos.Z)) / (math.Pow(cameraDir.X, 2) + math.Pow(cameraDir.Y, 2) + math.Pow(cameraDir.Z, 2))
	q := (math.Pow(pixel.X - s.pos.X, 2) + math.Pow(pixel.Y - s.pos.Y, 2) + math.Pow(pixel.Z - s.pos.Z, 2) - math.Pow(s.rad, 2)) / (math.Pow(cameraDir.X, 2) + math.Pow(cameraDir.Y, 2) + math.Pow(cameraDir.Z, 2))

	a1 = -1. * p / 2. + math.Sqrt(math.Pow(p / 2., 2) - q)
	a2 = -1. * p / 2. - math.Sqrt(math.Pow(p / 2., 2) - q)

	int1 := pixel.Add(cameraDir.MultiplyByScalar(a1))
	int2 := pixel.Add(cameraDir.MultiplyByScalar(a2))

	// implementation flawed, should consider only printing object when camera is outside
	// becomes important for reflections
	if a.Inside(s) {
		return nil
	} else if math.IsNaN(a1) {
		return nil
	} else if a1 == a2 && a1 > 0. {
		return &int1
	} else {
		if a1 > 0. && int1.Sub(pixel).Length() < int2.Sub(pixel).Length() {
			return &int1
		} else if a2 > 0. {
			return &int2
		} else {
			return nil
		}
	}
}

// func intersects(cameraPos vec.Vector, s Sphere) bool {
// 	x1, x2, x3 := 0., 0., 0.
// 	if cameraDirX == 0 {
// 		x1 = cameraPos.X
// 	} else {
// 		x1 = (s.pos.X - cameraPos.X) / cameraDirX
// 	}
// 	if cameraDirY == 0 {
// 		x2 = cameraPos.Y
// 	} else {
// 		x2 = (s.pos.Y - cameraPos.Y) / cameraDirY
// 	}
// 	if cameraDirZ == 0 {
// 		x3 = cameraPos.Z
// 	} else {
// 		x3 = (s.pos.Z - cameraPos.Z) / cameraDirZ
// 	}
// 	return s.rad >= vec.Vector{x1,x2,x3}.Sub(s.pos).Length()
// }

func CreateExampleSphereImage(a *Camera, sc Scene, d []byte) {
	fmt.Println("creating example Sphere image")

	s1 := Sphere{vec.Vector{0., 20., 0.}, 15., Material{byte(122), byte(122), byte(0)}}
	s2 := Sphere{vec.Vector{-5., 50., 10.}, 10., Material{byte(122), byte(50), byte(50)}}

	spec := a.middle.Add(a.dir.MultiplyByScalar(a.spectatorDistance))

	i := 0
	for i < len(a.pixel) {

		x := sphereIntersect(a, a.pixel[i], a.pixel[i].Sub(spec), s1)
		y := sphereIntersect(a, a.pixel[i], a.pixel[i].Sub(spec), s2)
		xl, yl := -1., -1.

		if x != nil {
			xl = x.Sub(a.pixel[i]).Length()
		}
		if y != nil {
			yl = y.Sub(a.pixel[i]).Length()
		}

		if x != nil && y != nil {
			if xl < yl {
				d[a.colorChannel * i]     = s1.mat.ColorR
				d[a.colorChannel * i + 1] = s1.mat.ColorG
				d[a.colorChannel * i + 2] = s1.mat.ColorB
			} else {
				d[a.colorChannel * i]     = s2.mat.ColorR
				d[a.colorChannel * i + 1] = s2.mat.ColorG
				d[a.colorChannel * i + 2] = s2.mat.ColorB
			}
		} else if x != nil {
			d[a.colorChannel * i]     = s1.mat.ColorR
			d[a.colorChannel * i + 1] = s1.mat.ColorG
			d[a.colorChannel * i + 2] = s1.mat.ColorB
		} else if y != nil {
			d[a.colorChannel * i]     = s2.mat.ColorR
			d[a.colorChannel * i + 1] = s2.mat.ColorG
			d[a.colorChannel * i + 2] = s2.mat.ColorB
		} else {
			d[a.colorChannel * i]     = sc.AmbCol.R
			d[a.colorChannel * i + 1] = sc.AmbCol.G
			d[a.colorChannel * i + 2] = sc.AmbCol.B
		}

		i++
	}
}