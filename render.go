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
	resolution_x = 1000
	resolution_z = 1000
	maxColor     = 255
	colorChannel = 3

	cameraHeightX = cameraHeightZ * (float64(resolution_x) / float64(resolution_z))
	cameraHeightZ = 50.

	spectatorDistance = -100.

	backgroundR   = 0
	backgroundG   = 255
	backgroundB   = 255
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
		vec.Vector{0., 1., 0.},
		vec.Vector{0., 0., 1.},
		vec.Vector{-13., 0., 13.},
		float64(cameraHeightX),
		float64(cameraHeightZ),
		resolution_x,
		resolution_z,
		spectatorDistance,
		colorChannel)

	var objects []RayIntersector
	objects = append(objects, NewSphere(vec.Vector{0., 20., 0.}, 15., Material{Color{150, 150, 0}, .5, .5, .9, 50.}))
	objects = append(objects, NewSphere(vec.Vector{-15., 70., 10.}, 10., Material{Color{60, 0, 25}, 1., .9, .01, 1.}))
	objects = append(objects, NewSphere(vec.Vector{20., 200., 40.}, 80., Material{Color{30, 30, 30}, 1., .5, .8, 50.}))

	sc := Scene{Color{backgroundR, backgroundG, backgroundB},
		backgroundInt,
		[]LightSource{LightSource{vec.Vector{-20., 15., 30.}, Color{255, 255, 255}, 1.},
		LightSource{vec.Vector{-200., -200., 100.}, Color{255, 255, 255}, 1.}},
		objects,
		sceneMaxDepth}

	CreateExampleSphereImage(c, sc, data)

	os.WriteFile("/tmp/dat1.ppm", append(h, data...), 0644)
}

/*
To Do:
add shadow, reflection
need to first put raycast into seperate function to use recursion

cleanup common

*/
