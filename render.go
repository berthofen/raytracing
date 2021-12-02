package main

import (
	"os"
	"strconv"
	"math/rand"
	"time"
	vec "raytracing/vector"
	. "raytracing/objects"
	. "raytracing/common"
	. "raytracing/scene"
)

const (
	// resolution needs to be an even number
	resolution_x = 2000
	resolution_z = 2000
	maxColor     = 255
	colorChannel = 3

	cameraHeightX = cameraHeightZ * (float64(resolution_x) / float64(resolution_z))
	cameraHeightZ = 50.
	cameraDirX = 0.
	cameraDirY = 1.
	cameraDirZ = 0.

	spectatorDistance = -10000.

	backgroundR = 0
	backgroundG = 255
	backgroundB = 255
	backgroundInt = 1.
)

type header struct {
	t byte
	width uint
	height uint
	maxColor byte
}

func(h header) print() []byte {
	header := "P" + strconv.Itoa(int(h.t)) + "\n" +
			strconv.Itoa(int(h.width)) + " " + strconv.Itoa(int(h.height)) + "\n" +
			strconv.Itoa(int(h.maxColor)) + "\n"
	return []byte(header)
}

func(h header) String() string {
	return string(h.print())
}


func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	h := header{6, resolution_x, resolution_z, maxColor}.print()
	data := make([]byte, resolution_x * resolution_z * colorChannel)

	c := CameraCreate(
		vec.Vector{0., 1., 0.}, 
		vec.Vector{0., 0., 1.}, 
		vec.Vector{0., 0., 0.}, 
		float64(cameraHeightX), 
		float64(cameraHeightZ), 
		resolution_x, 
		resolution_z, 
		spectatorDistance,
		colorChannel)

	var objects []RayIntersector
	objects = append(objects, NewSphere(vec.Vector{0., 20., 0.}, 15., Material{byte(150), byte(150), byte(0)}))
	objects = append(objects, NewSphere(vec.Vector{-5., 50., 10.}, 10., Material{byte(60), byte(25), byte(25)}))
	objects = append(objects, NewSphere(vec.Vector{-5., 20., 20.}, 4., Material{byte(30), byte(30), byte(30)}))

	sc := Scene{Color{backgroundR, backgroundG, backgroundB},
	 			backgroundInt, 
	 			LightSource{vec.Vector{-40., -15., 20.}},
				objects}

	CreateExampleSphereImage(c, sc, data)

	os.WriteFile("/tmp/dat1.ppm", append(h, data...), 0644)
}