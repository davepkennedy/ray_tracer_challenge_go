package main

import (
	"os"
	"fmt"
	"math"

	"log"
	"raytracer/pkg/rt"
)

func main() {
	camera, err := createCamera()
	if err != nil {
		log.Fatalf("failed to create camera: %v", err)
	}
	world := createWorld()
	render(world, camera)
}

func createWorld() *rt.World {
	world := rt.NewWorld()
	mainLight := createMainLight()
	backdrop, err := createBackdrop()
	if err != nil {
		log.Fatalf("failed to create backdrop: %v", err)
	}
	world.Add(backdrop)
	sphere, err := createSphere()
	if err != nil {
		log.Fatalf("failed to create sphere: %v", err)
	}
	world.Add(sphere)
	cubes, err := createCubes()
	if err != nil {
		log.Fatalf("failed to create cubes: %v", err)
	}
	for _, cube := range cubes {
		world.Add(cube)
	}
	world.Light = mainLight

	return world
}

func render(world *rt.World, camera *rt.Camera) {
	canvas, err := camera.Render(world)
	if err != nil {
		log.Fatalf("failed to save PPM: %v", err)
	}
	ppm := canvas.ToPPM()

	file, err := os.Create("cover.ppm")
	if err != nil {
		log.Fatalf("couldn't open output: %v", err)
	}
	defer file.Close()
	_, err = file.WriteString(ppm)
	if err != nil {
		log.Fatalf("failed to write PPM: %v", err)
	}

}

func createCamera() (*rt.Camera, error) {
	camera := rt.NewCamera(320, 240, 0.885)
	transform, err := rt.ViewTransformation(
		rt.NewPoint(-6, 6, -10),
		rt.NewPoint(6, 0, 6),
		rt.NewVector(-0.45, 1, 0))
	if err != nil {
		return nil, fmt.Errorf("failed to create camera transform: %w", err)
	}
	camera.Transform = transform
	return camera, err
}

func createMainLight() *rt.Light {
	return rt.NewPointLight(
		rt.NewPoint(50, 100, -50),
		rt.NewColor(1, 1, 1))
}

func createWhiteMaterial() *rt.Material {
	return &rt.Material{
		Color:    rt.NewColor(1, 1, 1),	
		Diffuse:  0.7,
		Ambient: 0.1,
		Specular: 0.0,
		Reflective: 0.1,
	}
}

func createBlueMaterial() *rt.Material {
	m := createWhiteMaterial()
	m.Color = rt.NewColor(0.537, 0.831, 0.914)
	return m
}

func createRedMaterial() *rt.Material {
	m := createWhiteMaterial()
	m.Color = rt.NewColor(0.941, 0.322, 0.388)
	return m
}

func createPurpleMaterial() *rt.Material {
	m := createWhiteMaterial()
	m.Color = rt.NewColor(0.373, 0.404, 0.550)
	return m
}

func createStandardTransform() (*rt.Matrix, error) {
	translate := rt.Translation(1, -1, 1)
	scale := rt.Scaling(0.5, 0.5, 0.5)
	return rt.Multiply(translate, scale)
}

func createLargeObjectTransform() (*rt.Matrix, error) {
	standard, err := createStandardTransform()
	if err != nil {
		return nil, fmt.Errorf("failed to create standard transform: %w", err)
	}
	scale := rt.Scaling(3.5, 3.5, 3.5)
	return rt.Multiply(standard, scale)
}

func createMediumObjectTransform() (*rt.Matrix, error) {
	standard, err := createStandardTransform()
	if err != nil {
		return nil, fmt.Errorf("failed to create standard transform: %w", err)
	}
	scale := rt.Scaling(3.0, 3.0, 3.0)
	return rt.Multiply(standard, scale)
}

func createSmallObjectTransform() (*rt.Matrix, error) {
	standard, err := createStandardTransform()
	if err != nil {
		return nil, fmt.Errorf("failed to create standard transform: %w", err)
	}
	scale := rt.Scaling(2, 2, 2)
	return rt.Multiply(standard, scale)
}

func createBackdrop() (*rt.Shape, error) {
	shape := rt.NewPlane()
	transform, err := rt.Multiply(
		rt.Translation(0, 0, 500),
		rt.RotationX(math.Pi/2),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create backdrop transform: %w", err)
	}

	shape.Transform = transform
	shape.Material.Color = rt.NewColor(1, 1, 1)
	shape.Material.Ambient = 1.0
	return shape, nil
}

func createSphere() (*rt.Shape, error) {
	shape := rt.NewSphere()
	transform, err := createLargeObjectTransform()
	if err != nil {
		return nil, fmt.Errorf("failed to create sphere transform: %w", err)
	}
	shape.Transform = transform

	shape.Material.Color = rt.NewColor(0.373, 0.404, 0.550)
	shape.Material.Diffuse = 0.2
	shape.Material.Ambient = 0.0
	shape.Material.Specular = 1.0
	shape.Material.Reflective = 0.7
	shape.Material.Shininess = 200.0
	shape.Material.Transparency = 0.7
	shape.Material.RefractiveIndex = 1.5

	return shape, nil
}

func createCubes() ([]*rt.Shape, error) {
	whiteMaterial := createWhiteMaterial()
	blueMaterial := createBlueMaterial()
	redMaterial := createRedMaterial()
	purpleMaterial := createPurpleMaterial()

	largeObject, err := createLargeObjectTransform()
	if err != nil {
		return nil, fmt.Errorf("failed to create large object transform: %w", err)
	}
	mediumObject, err := createMediumObjectTransform()
	if err != nil {
		return nil, fmt.Errorf("failed to create medium object transform: %w", err)
	}
	smallObject, err := createSmallObjectTransform()
	if err != nil {
		return nil, fmt.Errorf("failed to create small object transform: %w", err)
	}
	defs := []struct {
		material *rt.Material
		baseTransform *rt.Matrix
		objectTransform *rt.Matrix
	}{
		{whiteMaterial, mediumObject, rt.Translation(4.0, 0.0, 0.0)},
		{blueMaterial, largeObject, rt.Translation(8.5, 1.5, -0.5)},
		{redMaterial, largeObject, rt.Translation(0, 0, 4)},
		{whiteMaterial, smallObject, rt.Translation(4, 0, 4)},
		{purpleMaterial, mediumObject, rt.Translation(7.5, 0.5, 4)},
		{whiteMaterial, mediumObject, rt.Translation(-0.25, 0.25, 8)},
		{blueMaterial, largeObject, rt.Translation(4.0, 1.0, 7.5)},
		{redMaterial, mediumObject, rt.Translation(10, 2, 7.5)},
		{whiteMaterial, smallObject, rt.Translation(8, 2, 12)},
		{whiteMaterial, smallObject, rt.Translation(20, 1, 9)},
		{blueMaterial, largeObject, rt.Translation(-0.5, -5, 0.25)},
		{redMaterial, largeObject, rt.Translation(4, -4, 0)},
		{whiteMaterial, largeObject, rt.Translation(8.5, -4, 0)},
		{whiteMaterial, largeObject, rt.Translation(0, -4, 4)},
		{purpleMaterial, largeObject, rt.Translation(-0.5, -4.5, 8)},
		{whiteMaterial, largeObject, rt.Translation(0, -8, 4)},
		{whiteMaterial, largeObject, rt.Translation(-0.5, -8.5, 8)}}

	cubes := make([]*rt.Shape, len(defs))
	for i, def := range defs {
		cube := rt.NewCube()
		transform, err := rt.Multiply(def.objectTransform, def.baseTransform)
		if err != nil {
			return nil, fmt.Errorf("failed to create cube transform: %w", err)
		}
		cube.Transform = transform
		cube.Material = def.material
		cubes[i] = cube
	}
	return cubes, nil
}