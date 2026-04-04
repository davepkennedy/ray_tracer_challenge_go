package test

import (
	"context"
	"fmt"
	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func InitializeSphereScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w+) ← sphere\(\)$`,
		func(ctx context.Context, name string) context.Context {
			return context.WithValue(ctx, shapeKey{name}, rt.NewSphere())
		})

	sc.Given(
		`^(\w+) ← glass_sphere\(\)$`,
		func(ctx context.Context, name string) context.Context {
			return context.WithValue(ctx, shapeKey{name}, rt.GlassSphere())
		})

	sc.Given(
		`^set_transform\((\w+), (\w+)\)$`,
		func(ctx context.Context, shapeName, matrixName string) (context.Context, error) {
			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return ctx, err
			}

			matrix, err := getMatrix(ctx, matrixName)
			if err != nil {
				return ctx, err
			}

			shape.Transform = matrix
			return ctx, nil
		})
	sc.Given(
		`^set_transform\((\w+), translation\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		func(ctx context.Context, shapeName string, x, y, z float64) (context.Context, error) {
			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return ctx, err
			}

			matrix := rt.Translation(x, y, z)
			shape.Transform = matrix
			return ctx, nil
		})
	sc.Given(
		`^(\w+).material.ambient ← (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, val float64) error {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return nil
			}

			shape.Material.Ambient = val
			return nil
		})
	sc.Given(
		`^(\w+).light ← point_light\(point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\), color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		func(ctx context.Context, dest string, x, y, z, r, g, b float64) error {
			world, err := getWorld(ctx, dest)
			if err != nil {
				return err
			}

			world.Light = rt.NewPointLight(
				rt.NewPoint(x, y, z),
				rt.NewColor(r, g, b),
			)
			return nil
		})

	sc.When(
		`^(\w+) ← intersect\((\w), (\w)\)$`,
		func(ctx context.Context, dest, shapeName, rayName string) (context.Context, error) {
			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return ctx, err
			}

			ray, err := getRay(ctx, rayName)
			if err != nil {
				return ctx, err
			}

			intersections := shape.Intersect(ray)
			return context.WithValue(ctx, intersectionsKey{dest}, intersections), nil
		})
	sc.When(
		`^set_transform\((\w+), (\w+)\)$`,
		func(ctx context.Context, shapeName, matrixName string) (context.Context, error) {
			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return ctx, err
			}

			matrix, err := getMatrix(ctx, matrixName)
			if err != nil {
				return ctx, err
			}

			shape.Transform = matrix
			return ctx, nil
		})
	sc.When(
		`^set_transform\((\w+), scaling\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		func(ctx context.Context, shapeName string, x, y, z float64) (context.Context, error) {
			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return ctx, err
			}

			matrix := rt.Scaling(x, y, z)
			shape.Transform = matrix
			return ctx, nil
		})
	sc.When(
		`^set_transform\((\w+), translation\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		func(ctx context.Context, shapeName string, x, y, z float64) (context.Context, error) {
			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return ctx, err
			}

			matrix := rt.Translation(x, y, z)
			shape.Transform = matrix
			return ctx, nil
		})
	sc.When(
		`^(\w) ← normal_at\((\w+), point\((\S+), (\S+), (\S+)\)\)$`,
		func(ctx context.Context, dest, source, x, y, z string) (context.Context, error) {
			shape, err := getShape(ctx, source)
			if err != nil {
				return ctx, err
			}

			point, err := pointFromStrings(x, y, z)
			if err != nil {
				return ctx, err
			}
			normal, err := shape.NormalAt(point)
			if err != nil {
				return ctx, err
			}

			return context.WithValue(ctx, tupleKey{dest}, normal), nil
		})
	sc.When(
		`^(\w+) ← (\w+).material$`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			shape, err := getShape(ctx, source)
			if err != nil {
				return ctx, err
			}

			return context.WithValue(ctx, materialKey{dest}, shape.Material), nil
		})
	sc.When(
		`^(\w+).material ← (\w+)$`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return ctx, err
			}

			material, err := getMaterial(ctx, source)
			if err != nil {
				return ctx, err
			}

			shape.Material = material
			return ctx, nil
		})

	sc.Then(
		`^s.transform = (\w+)$`,
		func(ctx context.Context, matrixName string) error {
			shapeName := "s"
			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return err
			}

			matrix, err := getMatrix(ctx, matrixName)
			if err != nil {
				return err
			}

			if !shape.Transform.Equal(matrix) {
				return fmt.Errorf("expected %s, got %s", matrix, shape.Transform)
			}
			return nil
		})
	sc.Then(
		`^(\w+).material = (\w+)$`,
		func(ctx context.Context, dest, source string) error {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return err
			}

			material, err := getMaterial(ctx, source)
			if err != nil {
				return err
			}

			if !shape.Material.Equal(material) {
				return fmt.Errorf("expected %s, got %s", material, shape.Material)
			}
			return nil
		})

	sc.Then(
		`^(\w+).material.transparency = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, val float64) error {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return err
			}

			if shape.Material.Transparency != val {
				return fmt.Errorf("expected %f, got %f", val, shape.Material.Transparency)
			}
			return nil
		})

	sc.Then(
		`^(\w+).material.refractive_index = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, val float64) error {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return err
			}

			if shape.Material.RefractiveIndex != val {
				return fmt.Errorf("expected %f, got %f", val, shape.Material.RefractiveIndex)
			}
			return nil
		})
}
