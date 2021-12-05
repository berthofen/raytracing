package main

import (
	"math/rand"
	"os"
	. "raytracing/common"
	. "raytracing/objects"
	. "raytracing/scene"
	vec "raytracing/vector"
	"strconv"
	"time"
)

const (
	// resolution needs to be an even number
	resolution_x = 3000
	resolution_z = 3000
	maxColor     = 255
	colorChannel = 3

	cameraHeightX = cameraHeightZ * (float64(resolution_x) / float64(resolution_z))
	cameraHeightZ = 20.

	spectatorDistance = -100.

	backgroundR   = 51
	backgroundG   = 179
	backgroundB   = 204
	backgroundInt = 1.

	sceneMaxDepth = 4
)

type header struct {
	t        byte
	width    uint
	height   uint
	maxColor byte
}

func (h header) print() []byte {
	header := "P" + strconv.Itoa(int(h.t)) + "\n" +
		strconv.Itoa(int(h.width)) + " " + strconv.Itoa(int(h.height)) + "\n" +
		strconv.Itoa(int(h.maxColor)) + "\n"
	return []byte(header)
}

func (h header) String() string {
	return string(h.print())
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	h := header{6, resolution_x, resolution_z, maxColor}.print()
	data := make([]byte, resolution_x*resolution_z*colorChannel)

	c := CameraCreate(
		vec.Vector{0., -0.1, -1.},
		vec.Vector{0.2, 1., 0.},
		vec.Vector{0., 3., 10.},
		float64(cameraHeightX),
		float64(cameraHeightZ),
		resolution_x,
		resolution_z,
		spectatorDistance,
		colorChannel)

	ivory := Material{Color{122, 122, 77}, .0, 0.6, 0.3, 0.1, 50.}
	glass := Material{Color{153, 179, 204}, .0, 0.0, 0.5, 0.1, 125.}
	red_rubber := Material{Color{77, 26, 26}, .0, 0.9, 0.1, 0.0, 10.}
	mirror := Material{Color{0, 0, 0}, .0, 0.0, 10.0, 0.8, 1425.}

	var objects []RayIntersector
	objects = append(objects, NewSphere(vec.Vector{-3., 0., -16.}, 2., ivory))
	objects = append(objects, NewSphere(vec.Vector{-1., -1.5, -12.}, 2., glass))
	objects = append(objects, NewSphere(vec.Vector{1.5, -0.5, -18.}, 3., red_rubber))
	objects = append(objects, NewSphere(vec.Vector{7., 5., -18.}, 4., mirror))

	p := NewPlain(vec.Vector{2, 2, 0},
		vec.Vector{1, 1, 1},
		vec.Vector{1, 0, 0},
		red_rubber,
		5.,
		5.)

	sc := Scene{Color{backgroundR, backgroundG, backgroundB},
		backgroundInt,
		[]LightSource{LightSource{vec.Vector{-20., 20., 20.}, Color{255, 255, 255}, 1.5},
			LightSource{vec.Vector{30., 50., -25.}, Color{255, 255, 255}, 1.8},
			LightSource{vec.Vector{30., 20., 30.}, Color{255, 255, 255}, 1.7}},
		objects,
		sceneMaxDepth}

	CreateExampleSphereImage(c, sc, data)

	os.WriteFile("/Users/jhesselmann/Documents/Projects/go/raytracing/dat1.ppm", append(h, data...), 0644)
}

/*
To Do:
add shadow, reflection
need to first put raycast into seperate function to use recursion

cleanup common

*/
