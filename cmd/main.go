package main

import (
	"fmt"
	"log"
	//"math"

	"raytracer/pkg/layout"
	"raytracer/pkg/rt"
)

func main() {
/*
	var t *rt.Matrix
	var err error

	floor := rt.NewSphere()
	floor.Transform = rt.Scaling(10, 0.01, 10)
	floor.Material = rt.NewMaterial()
	floor.Material.Color = rt.NewColor(1, 0.9, 0.9)
	floor.Material.Specular = 0

	leftWall := rt.NewSphere()
	t, err = rt.Multiply(rt.Translation(0, 0, 5), rt.RotationY(-math.Pi/4), rt.RotationX(math.Pi/2), rt.Scaling(10, 0.01, 10))
	if err != nil {
		log.Fatalf("fail: %v", err)
	}
	leftWall.Transform = t
	leftWall.Material = floor.Material

	rightWall := rt.NewSphere()
	t, err = rt.Multiply(rt.Translation(0, 0, 5), rt.RotationY(math.Pi/4), rt.RotationX(math.Pi/2), rt.Scaling(10, 0.01, 10))
	if err != nil {
		log.Fatalf("fail: %v", err)
	}
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
	t, err = rt.Multiply(rt.Translation(1.5, 0.5, -0.5), rt.Scaling(0.5, 0.5, 0.5))
	if err != nil {
		log.Fatalf("fail: %v", err)
	}
	right.Transform = t
	right.Material = rt.NewMaterial()
	right.Material.Color = rt.NewColor(0.5, 1, 0.1)
	right.Material.Diffuse = 0.7
	right.Material.Specular = 0.3

	left := rt.NewSphere()
	t, err = rt.Multiply(rt.Translation(-1.5, 0.33, -0.75), rt.Scaling(0.33, 0.33, 0.33))
	if err != nil {
		log.Fatalf("fail: %v", err)
	}
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

	if err != nil {
		log.Fatalf("fail: %v", err)
	}
	camera.Transform = t
	// render the result to a canvas.​ ​
	canvas, err := camera.Render(world)
	if err != nil {
		log.Fatalf("fail: %v", err)
	}

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
	*/
	/*
	entries, err := layout.LoadScene("cover.yaml")
	if err != nil {
		log.Fatalf("failed to load scene: %v", err)
	}

	for _, e := range entries {
		switch {
		case e.Add == "camera":
			fmt.Printf("camera: %dx%d fov=%.3f from=%v to=%v up=%v\n",
				e.Width, e.Height, e.FieldOfView, e.From, e.To, e.Up)

		case e.Add == "light":
			fmt.Printf("light: at=%v intensity=%v\n", e.At, e.Intensity)

		case e.Add != "": // plane, sphere, cube, etc.
			fmt.Printf("shape %q: material=%+v transform=%+v\n",
				e.Add, e.Material, e.Transform)

		case e.Define != "":
			fmt.Printf("define %q (extends %q)\n", e.Define, e.Extend)
			fmt.Printf("\t%#v", e.Value)
		}
	}
	*/
	sceneFile, err := layout.LoadSceneFile("cover.json")
	if err != nil {
		log.Fatalf("failed to load scene: %v", err)
	}

	scene := layout.NewScene()
	scene.AddCamera(sceneFile.Camera)
	scene.AddLight(sceneFile.Lights[0])
	addDefines(scene, sceneFile)
	addShapes(scene, sceneFile)

	err = scene.Save("cover.ppm")
	if err != nil {
		log.Fatalf("failed to save scene: %v", err)
	}
}

func combineMaterials(base, override *rt.Material) *rt.Material {
	result := *base // copy all fields from base
	if override.Color != nil {
		result.Color = override.Color
	}
	if override.Diffuse != 0 {
		result.Diffuse = override.Diffuse
	}
	if override.Ambient != 0 {
		result.Ambient = override.Ambient
	}
	if override.Specular != 0 {
		result.Specular = override.Specular
	}
	if override.Shininess != 0 {
		result.Shininess = override.Shininess
	}
	if override.Reflective != 0 {
		result.Reflective = override.Reflective
	}
	if override.Transparency != 0 {
		result.Transparency = override.Transparency
	}
	if override.RefractiveIndex != 0 {
		result.RefractiveIndex = override.RefractiveIndex
	}
	return &result
}

func addDefines(scene layout.Scene, sceneFile *layout.SceneFile) {
	for _, define := range sceneFile.Defines {
		switch define.Type {
		case "material":
			material, err := define.DefineMaterial()
			if err != nil {
				log.Fatalf("failed to decode material %q: %v", define.Name, err)
			}
			if define.Extends != "" {
				baseMaterial, ok := scene.GetMaterial(define.Extends)
				if !ok {
					log.Fatalf("material %q extends unknown material %q", define.Name, define.Extends)
				}
				material = combineMaterials(baseMaterial, material)
			}
			scene.CacheMaterial(define.Name, material)

		case "transform":
			transform, err := define.DefineTransforms()
			if err != nil {
				log.Fatalf("failed to decode transform %q: %v", define.Name, err)
			}
			if define.Extends != "" {
				baseTransform, ok := scene.GetTransform(define.Extends)
				if !ok {
					log.Fatalf("transform %q extends unknown transform %q", define.Name, define.Extends)
				}
				transform, err = baseTransform.MultiplyMatrix(transform)
				if err != nil {
					log.Fatalf("failed to combine transform %q with base %q: %v", define.Name, define.Extends, err)
				}
			}
			scene.CacheTransform(define.Name, transform)

		default:
			log.Fatalf("unknown define type %q for %q", define.Type, define.Name)
		}
	}
}

func createShape (shapeType string) (*rt.Shape, error) {
	switch shapeType {
	case "plane":
		return rt.NewPlane(), nil
	case "sphere":
		return rt.NewSphere(), nil
	case "cube":
		return rt.NewCube(), nil
	default:
		return nil, fmt.Errorf("unknown shape type %q", shapeType)
	}
}

func addShapes(scene layout.Scene, sceneFile *layout.SceneFile) {
	for _, shapeDef := range sceneFile.Shapes {
		shape, err := createShape(shapeDef.Type)
		if err != nil {
			log.Fatalf("failed to create shape %q: %v", shapeDef.Type, err)
		}
		// Apply material and transform to the shape
		if shapeDef.Material.Ref != "" {
			material, ok := scene.GetMaterial(shapeDef.Material.Ref)
			if !ok {
				log.Fatalf("shape material %q not found", shapeDef.Material.Ref)
			}
			shape.Material = material
		} else if shapeDef.Material.Def != nil {
			shape.Material = shapeDef.Material.Def.ToMaterial()
		}

		transform := rt.IdentityMatrix()
		for _, t := range shapeDef.Transform {
			var tMatrix *rt.Matrix
			if t.Ref != "" {
				var ok bool
				tMatrix, ok = scene.GetTransform(t.Ref)
				if !ok {
					log.Fatalf("shape transform %q not found", t.Ref)
				}
			} else if t.Def != nil {
				var err error
				tMatrix, err = t.Def.ToMatrix()
				if err != nil {
					log.Fatalf("failed to convert inline transform to matrix: %v", err)
				}
			} else {
				log.Fatalf("shape transform must have either ref or def")
			}
			transform, err = tMatrix.MultiplyMatrix(transform)
			if err != nil {
				log.Fatalf("failed to combine shape transform: %v", err)
			}
		}
		shape.Transform = transform

		scene.AddShape(shape)
			
	}
}