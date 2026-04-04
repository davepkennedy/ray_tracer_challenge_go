package test

import (
	"context"
	"fmt"
	"strconv"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func newShapeOfType(s string) *rt.Shape {
	switch s {
	case "sphere":
		return rt.NewSphere()
	case "plane":
		return rt.NewPlane()
	}
	return NewTestShape()
}

func InitializeWorldScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`(\w) ← world\(\)`,
		func(ctx context.Context, dest string) context.Context {
			return setWorld(ctx, dest, rt.NewWorld())
		})
	sc.Given(
		`(\w) ← default_world\(\)`,
		func(ctx context.Context, dest string) context.Context {
			return setWorld(ctx, dest, rt.DefaultWorld())
		})

	sc.Given(
		`(\w+) ← (\w+)\(\) with:`,
		func(ctx context.Context, dest, shapeType string, table *godog.Table) (context.Context, error) {
			sphere := newShapeOfType(shapeType)
			sphere.Material = rt.NewMaterial()
			for _, row := range table.Rows {
				switch row.Cells[0].Value {
				case "material.color":
					v, err := colorFromString(row.Cells[1].Value)
					if err != nil {
						return ctx, err
					}
					sphere.Material.Color = v
				case "material.diffuse":
					v, err := strconv.ParseFloat(row.Cells[1].Value, 64)
					if err != nil {
						return ctx, err
					}
					sphere.Material.Diffuse = v
				case "material.specular":
					v, err := strconv.ParseFloat(row.Cells[1].Value, 64)
					if err != nil {
						return ctx, err
					}
					sphere.Material.Specular = v
				case "material.reflective":
					v, err := strconv.ParseFloat(row.Cells[1].Value, 64)
					if err != nil {
						return ctx, err
					}
					sphere.Material.Reflective = v
				case "material.refractive_index":
					v, err := strconv.ParseFloat(row.Cells[1].Value, 64)
					if err != nil {
						return ctx, err
					}
					sphere.Material.RefractiveIndex = v
				case "material.transparency":
					v, err := strconv.ParseFloat(row.Cells[1].Value, 64)
					if err != nil {
						return ctx, err
					}
					sphere.Material.Transparency = v
				case "material.ambient":
					v, err := strconv.ParseFloat(row.Cells[1].Value, 64)
					if err != nil {
						return ctx, err
					}
					sphere.Material.Ambient = v
				case "transform":
					t, err := transformFromString(row.Cells[1].Value)
					if err != nil {
						return ctx, err
					}
					sphere.Transform = t
				default:
					return ctx, fmt.Errorf("unrecognized property: %s", row.Cells[0].Value)
				}
			}
			return setShape(ctx, dest, sphere), nil
		})

	sc.Given(
		`^(\w+) ← the first object in (\w+)$`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			world, err := getWorld(ctx, source)
			if err != nil {
				return ctx, err
			}

			if len(world.Shapes) < 1 {
				return ctx, fmt.Errorf("insufficient elements in world")
			}
			return setShape(ctx, dest, world.Shapes[0]), nil
		})
	sc.Given(
		`^(\w+) ← the second object in (\w+)$`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			world, err := getWorld(ctx, source)
			if err != nil {
				return ctx, err
			}

			if len(world.Shapes) < 2 {
				return ctx, fmt.Errorf("insufficient elements in world")
			}
			return setShape(ctx, dest, world.Shapes[1]), nil
		})

	sc.When(
		`(\w) ← default_world\(\)`,
		func(ctx context.Context, dest string) context.Context {
			return setWorld(ctx, dest, rt.DefaultWorld())
		})
	sc.When(
		`(\w+) ← intersect_world\((\w+), (\w+)\)`,
		func(ctx context.Context, dest, worldName, rayName string) (context.Context, error) {
			world, err := getWorld(ctx, worldName)
			if err != nil {
				return ctx, err
			}
			ray, err := getRay(ctx, rayName)
			if err != nil {
				return ctx, err
			}

			xs := world.Intersect(ray)
			return setIntersections(ctx, dest, xs), nil
		})
	sc.When(
		`^(\w+) ← shade_hit\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, worldName, compsName string) (context.Context, error) {
			world, err := getWorld(ctx, worldName)
			if err != nil {
				return ctx, err
			}
			comps, err := getComps(ctx, compsName)
			if err != nil {
				return ctx, err
			}

			color, err := world.ShadeHit(comps, 5)
			if err != nil {
				return ctx, err
			}

			return setColor(ctx, dest, color), err
		})
	sc.When(
		`^(\w+) ← shade_hit\((\w+), (\w+), (\d+)\)$`,
		func(ctx context.Context, dest, worldName, compsName string, remaining int) (context.Context, error) {
			world, err := getWorld(ctx, worldName)
			if err != nil {
				return ctx, err
			}
			comps, err := getComps(ctx, compsName)
			if err != nil {
				return ctx, err
			}

			color, err := world.ShadeHit(comps, remaining)
			if err != nil {
				return ctx, err
			}

			return setColor(ctx, dest, color), err
		})
	sc.When(
		`^(\w+) ← color_at\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, worldName, rayName string) (context.Context, error) {
			world, err := getWorld(ctx, worldName)
			if err != nil {
				return ctx, err
			}
			ray, err := getRay(ctx, rayName)
			if err != nil {
				return ctx, err
			}

			color, err := world.ColorAt(ray, 5)
			if err != nil {
				return ctx, err
			}

			return setColor(ctx, dest, color), nil
		})
	sc.When(
		`^(\w+) ← reflected_color\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, source, compsName string) (context.Context, error) {
			world, err := getWorld(ctx, source)
			if err != nil {
				return ctx, err
			}

			comps, err := getComps(ctx, compsName)
			if err != nil {
				return ctx, err
			}

			color, err := world.ReflectedColor(comps, 5)
			if err != nil {
				return ctx, err
			}

			return setColor(ctx, dest, color), nil
		})

	sc.When(
		`^(\w+) ← reflected_color\((\w+), (\w+), 0\)$`,
		func(ctx context.Context, dest, source, compsName string) (context.Context, error) {
			world, err := getWorld(ctx, source)
			if err != nil {
				return ctx, err
			}

			comps, err := getComps(ctx, compsName)
			if err != nil {
				return ctx, err
			}

			color, err := world.ReflectedColor(comps, 0)
			if err != nil {
				return ctx, err
			}

			return setColor(ctx, dest, color), nil
		})

	sc.When(
		`^(\w+) ← refracted_color\((\w+), (\w+), (\d+)\)$`,
		func(ctx context.Context, dest, source, compsName string, remaining int) (context.Context, error) {
			world, err := getWorld(ctx, source)
			if err != nil {
				return ctx, err
			}

			comps, err := getComps(ctx, compsName)
			if err != nil {
				return ctx, err
			}

			color, err := world.RefractedColor(comps, remaining)
			if err != nil {
				return ctx, err
			}

			return setColor(ctx, dest, color), nil
		})

	sc.Then(
		`^(\w) contains no objects$`,
		func(ctx context.Context, dest string) error {
			world, err := getWorld(ctx, dest)
			if err != nil {
				return err
			}
			if len(world.Shapes) > 0 {
				return fmt.Errorf("expected world to contain no objects")
			}
			return nil
		})
	sc.Then(
		`^(\w) has no light source$`,
		func(ctx context.Context, dest string) error {
			world, err := getWorld(ctx, dest)
			if err != nil {
				return err
			}
			if world.Light != nil {
				return fmt.Errorf("expected world to contain no lights")
			}
			return nil
		})

	sc.Then(
		`^(\w+).light = (\w+)$`,
		func(ctx context.Context, dest, name string) error {
			world, err := getWorld(ctx, dest)
			if err != nil {
				return err
			}

			light, err := getLight(ctx, name)
			if err != nil {
				return err
			}

			if !world.Light.Equal(light) {
				return fmt.Errorf("expected %s, got %s", light, world.Light)
			}
			return nil
		})

	sc.Then(
		`^(\w+) contains (\w+)$`,
		func(ctx context.Context, dest, shapeName string) error {
			world, err := getWorld(ctx, dest)
			if err != nil {
				return err
			}
			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return err
			}

			if !world.Contains(shape) {
				return fmt.Errorf("expected %s to contain %s but it did not", dest, shapeName)
			}
			return nil
		})
	sc.Then(
		`^(\w+) = (\w+).material.color$`,
		func(ctx context.Context, dest, source string) error {
			color, err := getColor(ctx, dest)
			if err != nil {
				return err
			}
			shape, err := getShape(ctx, source)
			if err != nil {
				return err
			}

			if !color.Equal(shape.Material.Color) {
				return fmt.Errorf("expected %s, got %s", shape.Material.Color, color)
			}
			return nil
		})
	sc.Then(
		`^is_shadowed\((\w+), (\w+)\) is (true|false)$`,
		func(ctx context.Context, worldName, pointName, shadowedFlag string) error {
			world, err := getWorld(ctx, worldName)
			if err != nil {
				return err
			}
			point, err := getTuple(ctx, pointName)
			if err != nil {
				return err
			}

			expect := shadowedFlag == "true"
			actual := world.IsShadowed(point)
			if actual != expect {
				return fmt.Errorf("expected %v, got %v", expect, actual)
			}

			return nil
		})
	sc.Then(
		`^color_at\((\w+), (\w+)\) should terminate successfully$`,
		func(ctx context.Context, worldName, rayName string) error {
			world, err := getWorld(ctx, worldName)
			if err != nil {
				return err
			}

			ray, err := getRay(ctx, rayName)
			if err != nil {
				return err
			}

			_, err = world.ColorAt(ray, 5)
			if err != nil {
				return fmt.Errorf("expected color_at to terminate successfully, but got error: %v", err)
			}
			return nil
		})
}
