package main

import (
	"log"
	"math"
	"os"
	"raytracer/pkg/rt"
)

func main() {

	var t *rt.Matrix
	var err error
	
	floor := rt.NewSphere()
	floor.Transform = rt.Scaling(10, 0.01, 10) 
	floor.Material = rt.NewMaterial()
	floor.Material.Color = rt.NewColor(1, 0.9, 0.9)
	floor.Material.Specular = 0
	
	leftWall := rt.NewSphere()
	t, err = rt.Multiply(rt.Translation(0, 0, 5), rt.RotationY(-math.Pi/4), rt.RotationX(math.Pi/2), rt.Scaling(10, 0.01, 10))
	if err != nil {log.Fatalf("fail: %v", err)}
	leftWall.Transform = t
	leftWall.Material = floor.Material
	
	rightWall := rt.NewSphere()
	t, err = rt.Multiply(rt.Translation(0, 0, 5), rt.RotationY(math.Pi/4), rt.RotationX(math.Pi/2), rt.Scaling(10, 0.01, 10))
	if err != nil {log.Fatalf("fail: %v", err)}
	rightWall.Transform = t
	rightWall.Material = floor.Material

	middle := rt.NewSphere()
	middle.Transform = rt.Translation(-0.5, 1, 0.5)
	middle.Material = rt.NewMaterial()
	middle.Material.Color = rt.NewColor(0.1, 1, 0.5)
	middle.Material.Diffuse = 0.7
	middle.Material.Specular = 0.3
	middle.Material.Pattern = rt.NewCheckersPattern(rt.NewColor(.7, .7, .7), rt.NewColor(0, 0.8, 0))
	middle.Material.Pattern.Transform = rt.Scaling(0.25, 0.25, 0.25)

	right := rt.NewSphere()
	t, err = rt.Multiply( rt.Translation(1.5, 0.5, -0.5), rt.Scaling(0.5, 0.5, 0.5))
	if err != nil {log.Fatalf("fail: %v", err)}
	right.Transform = t
	right.Material = rt.NewMaterial()
	right.Material.Color = rt.NewColor(0.5, 1, 0.1)
	right.Material.Diffuse = 0.7
	right.Material.Specular = 0.3

	left := rt.NewSphere()
	t, err = rt.Multiply(rt.Translation(-1.5, 0.33, -0.75), rt.Scaling(0.33, 0.33, 0.33))
	if err != nil {log.Fatalf("fail: %v", err)}
	left.Transform = t
	left.Material = rt.NewMaterial()
	left.Material.Color = rt.NewColor(1, 0.8, 0.1)
	left.Material.Diffuse = 0.7
	left.Material.Specular = 0.3

	world := rt.NewWorld()
	world.Light = rt.NewPointLight(rt.NewPoint(-10, 10, -10), rt.NewColor(1, 1, 1))
	world.Add(leftWall)
	world.Add(rightWall)
	world.Add(floor)
	world.Add(left)
	world.Add(middle)
	world.Add(right)

	camera := rt.NewCamera(640, 480, math.Pi/3)
	t, err = rt.ViewTransformation(
		rt.NewPoint(0, 1.5, -5),
		rt.NewPoint(0, 1, 0),
		rt.NewVector(0, 1, 0))

	if err != nil {log.Fatalf("fail: %v", err)}
	camera.Transform = t
	// render the result to a canvas.​ ​  
	canvas, err := camera.Render(world)
	if err != nil {log.Fatalf("fail: %v", err)}

	ppm := canvas.ToPPM()
	file, err := os.Create("render.ppm")
	if err != nil {
		log.Fatalf("Couldn't open output: %v", err)
	}
	defer file.Close()
	_, err = file.WriteString(ppm)
	if err != nil {
		log.Fatalf("Couldn't write file: %v", err)
	}
}