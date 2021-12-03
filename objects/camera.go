package objects

import (
	"fmt"
	"math"
	. "raytracing/common"
	. "raytracing/scene"
	vec "raytracing/vector"
)

type container interface {
	contains(a vec.Vector) bool
}

type Camera struct {
	dir               vec.Vector
	dirZ              vec.Vector // needs to be orthogonal to dir
	middle            vec.Vector
	width             float64
	height            float64
	resWidth          int
	resHeight         int
	spectatorDistance float64
	colorChannel      int
	pixel             []vec.Vector
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

func (c *Camera) SetPixels() {
	c.pixel = make([]vec.Vector, c.resWidth*c.resHeight)
	dirHeight := c.dirZ.Normalize().MultiplyByScalar(c.height / (float64(c.resHeight) - 1.))
	dirWidth := c.dir.Cross(c.dirZ).Normalize().MultiplyByScalar(c.width / (float64(c.resWidth) - 1.))

	// for even resolution there is no pixel in the middle thats why we add/sub 1/2 of the normalized dir vectors
	start := c.middle.Add(dirHeight.MultiplyByScalar(float64(c.resHeight) / 2.)).Sub(dirWidth.MultiplyByScalar(float64(c.resWidth) / 2.))
	start = start.Sub(dirHeight.MultiplyByScalar(.5)).Add(dirWidth.MultiplyByScalar(.5))

	scalarWidth, scalarHeight := 0., 0.
	for i := 0; i < len(c.pixel); i++ {
		scalarWidth = float64(i % c.resWidth)
		scalarHeight = float64(i / c.resWidth)

		c.pixel[i] = start.Add(dirWidth.MultiplyByScalar(scalarWidth)).Sub(dirHeight.MultiplyByScalar(scalarHeight))
	}

	// fmt.Println(c.pixel)
}

func (c *Camera) Inside(obj container) bool {
	return obj.contains(c.middle)
}

type Sphere struct {
	pos vec.Vector
	rad float64
	mat Material
}

func NewSphere(a vec.Vector, r float64, m Material) Sphere {
	return Sphere{a, r, m}
}

func (s Sphere) contains(a vec.Vector) bool {
	return math.Pow(a.X-s.pos.X, 2)+math.Pow(a.Y-s.pos.Y, 2)+math.Pow(a.Z-s.pos.Z, 2) <= math.Pow(s.rad, 2)
}

func (s Sphere) Intersect(a Ray) (*vec.Vector, float64, *Material, *vec.Vector) {
	a1, a2 := 0., 0.

	p := 2. * (a.Dir.X*(a.From.X-s.pos.X) + a.Dir.Y*(a.From.Y-s.pos.Y) + a.Dir.Z*(a.From.Z-s.pos.Z)) / (math.Pow(a.Dir.X, 2) + math.Pow(a.Dir.Y, 2) + math.Pow(a.Dir.Z, 2))
	q := (math.Pow(a.From.X-s.pos.X, 2) + math.Pow(a.From.Y-s.pos.Y, 2) + math.Pow(a.From.Z-s.pos.Z, 2) - math.Pow(s.rad, 2)) / (math.Pow(a.Dir.X, 2) + math.Pow(a.Dir.Y, 2) + math.Pow(a.Dir.Z, 2))

	a1 = -1.*p/2. + math.Sqrt(math.Pow(p/2., 2)-q)
	a2 = -1.*p/2. - math.Sqrt(math.Pow(p/2., 2)-q)

	int1 := a.From.Add(a.Dir.MultiplyByScalar(a1))
	int2 := a.From.Add(a.Dir.MultiplyByScalar(a2))

	var n1, n2 vec.Vector

	if &int1 != nil {
		n1 = int1.Sub(s.pos).Normalize()
	}
	if &int2 != nil {
		n2 = int2.Sub(s.pos).Normalize()
	}

	if s.contains(a.From) || (a1 < 0. && a2 < 0.) || math.IsNaN(a1) {
		return nil, -1., nil, nil
	} else if a1 == a2 {
		return &int1, a1, &s.mat, &n1
	} else {
		if a1 > 0. && int1.Sub(a.From).Length() < int2.Sub(a.From).Length() {
			return &int1, a1, &s.mat, &n1
		} else if a2 > 0. {
			return &int2, a2, &s.mat, &n2
		} else {
			return nil, -1., nil, nil
		}
	}
}

func calcSpecular(m Material, viewer *vec.Vector, n *vec.Vector, l *vec.Vector, int *vec.Vector) Color {
	lint := l.Sub(*int).Normalize()
	ref := n.MultiplyByScalar(2. * lint.Dot(*n)).Sub(lint).Normalize()
	view := viewer.Sub(*int).Normalize()

	if l.Dot(*n) < 0. {
		return Color{0, 0, 0}
	} else {
		return m.Col.Scale(m.Ks * math.Pow(ref.Dot(view), m.Alpha))
	}
}

func calcDiffuseL(m Material, n *vec.Vector, l *vec.Vector, int *vec.Vector) Color {
	//return byte(math.Min(255, float64(channel)+math.Max(0, math.RoundToEven(n.Dot((*l).Sub(*int).Normalize())*float64(channel)))))
	return m.Col.Scale(m.Kd * n.Dot((*l).Sub(*int).Normalize()))
}

func calcAmbient(m Material) Color {
	return m.Col.Scale(m.Ka)
}

func calcLight(m Material, viewer *vec.Vector, n *vec.Vector, l *vec.Vector, int *vec.Vector) Color {
	var col []Color

	col = append(col, calcAmbient(m))
	col = append(col, calcDiffuseL(m, n, l, int))
	col = append(col, calcSpecular(m, viewer, n, l, int))

	total := Color{0, 0, 0}
	for _, c := range col {
		total = total.Add(c)
	}

	return total
}

func CreateExampleSphereImage(a *Camera, sc Scene, d []byte) {
	fmt.Println("creating example Sphere image")

	spec := a.middle.Add(a.dir.MultiplyByScalar(a.spectatorDistance))

	i := 0
	for i < len(a.pixel) {

		var (
			closest_pos  *vec.Vector
			closest_mat  *Material
			closest_norm *vec.Vector
		)
		closest_len := math.Inf(1)
		for _, obj := range sc.Objects {
			x, xl, xm, xn := obj.Intersect(Ray{a.pixel[i], a.pixel[i].Sub(spec)})

			if xl > 0 && xl < closest_len {
				closest_pos = x
				closest_len = xl
				closest_mat = xm
				closest_norm = xn
			}
		}

		if closest_pos == nil {
			writeColor(d[a.colorChannel*i:a.colorChannel*i+3], sc.AmbCol)
		} else {
			writeColor(d[a.colorChannel*i:a.colorChannel*i+3], calcLight(*closest_mat, &a.pixel[i], closest_norm, &sc.Light.Pos, closest_pos))
		}
		i++
	}
}

func writeColor(d []byte, a Color) {
	d[0] = byte(a.R)
	d[1] = byte(a.G)
	d[2] = byte(a.B)
}
