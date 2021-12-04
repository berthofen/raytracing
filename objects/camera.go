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

func (s Sphere) Intersect(a Ray) (*vec.Vector, float64, *Material, *vec.Vector, *vec.Vector) {
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
	
	ref1, ref2 := n1.MultiplyByScalar(2. * a.From.Sub(int1).Dot(n1)).Sub(a.From.Sub(int1)),
					n2.MultiplyByScalar(2. * a.From.Sub(int2).Dot(n2)).Sub(a.From.Sub(int2))

	if s.contains(a.From) || (a1 < 0. && a2 < 0.) || math.IsNaN(a1) {
		return nil, -1., nil, nil, nil
	} else if a1 == a2 {
		return &int1, a1, &s.mat, &n1, &ref1
	} else {
		if a1 > 0. && int1.Sub(a.From).Length() < int2.Sub(a.From).Length() {
			return &int1, a1, &s.mat, &n1, &ref1
		} else if a2 > 0. {
			return &int2, a2, &s.mat, &n2, &ref2
		} else {
			return nil, -1., nil, nil, nil
		}
	}
}

func calcSpecular(m Material, pixelPos *vec.Vector, n *vec.Vector, lights []LightSource, int *vec.Vector) Color {
	view := pixelPos.Sub(*int).Normalize()
	var lint, ref vec.Vector

	total := Color{0, 0, 0}
	for _, light := range lights {
		lint = light.Pos.Sub(*int).Normalize()
		ref = n.MultiplyByScalar(2. * lint.Dot(*n)).Sub(lint).Normalize()

		if overlap := ref.Dot(view); overlap > 0. {
			total = total.Add(light.Col.Scale(m.Ks * light.Intens * math.Pow(overlap, m.Alpha)))
		}
	}

	return total
}

func calcDiffuse(m Material, n *vec.Vector, lights []LightSource, int *vec.Vector) Color {

	total := Color{0, 0, 0}
	for _, light := range lights {
		total = total.Add(m.Col.Scale(m.Kd * light.Intens * n.Dot((light.Pos).Sub(*int).Normalize())))
	}

	return total
}

func calcAmbient(m Material) Color {
	return m.Col.Scale(m.Ka)
}

func calcLight(m Material, pixelPos *vec.Vector, n *vec.Vector, l []LightSource, int *vec.Vector) Color {
	var col []Color

	col = append(col, calcAmbient(m))
	col = append(col, calcDiffuse(m, n, l, int))
	col = append(col, calcSpecular(m, pixelPos, n, l, int))

	total := Color{0, 0, 0}
	for _, c := range col {
		total = total.Add(c)
	}

	return total
}

func CreateExampleSphereImage(a *Camera, sc Scene, d []byte) {
	fmt.Println("creating example Sphere image")

	specPos := a.middle.Add(a.dir.MultiplyByScalar(a.spectatorDistance))

	i := 0
	for i < len(a.pixel) {

		color := castRay(Ray{a.pixel[i], a.pixel[i].Sub(specPos)}, sc, MinDepth)

		writeColor(d[a.colorChannel*i:a.colorChannel*i+3], color)
		i++
	}
}

func writeColor(d []byte, a Color) {
	d[0] = byte(a.R)
	d[1] = byte(a.G)
	d[2] = byte(a.B)
}

func castRay(ray Ray, scene Scene, depth uint) Color {
		if depth > scene.MaxDepth {
			return scene.AmbCol
		}

		var (
			closest_pos  *vec.Vector
			closest_mat  *Material
			closest_norm *vec.Vector
			closest_refl *vec.Vector
		)
		closest_len := math.Inf(1)
		for _, obj := range scene.Objects {
			x, xl, xm, xn, xr := obj.Intersect(ray)

			if xl > 0 && xl < closest_len {
				closest_pos = x
				closest_len = xl
				closest_mat = xm
				closest_norm = xn
				closest_refl = xr
			}
		}

		if closest_pos == nil {
			if depth == MinDepth {
				return scene.AmbCol
			} else {
				return scene.AmbCol
			}
		} else {
			c := castRay(Ray{*closest_pos, *closest_refl}, scene, depth + 1).Scale(closest_mat.Kr)
			return calcLight(*closest_mat, &ray.From, closest_norm, scene.Lights, closest_pos).Add(c)
		}
}

func doNothing(a interface{}) {
	return
}