package test

import (
	"context"
	"fmt"
	"math"
	"regexp"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func selectTransform(s string) (func(x, y, z float64) *rt.Matrix, error) {
	switch s {
	case "translation":
		return rt.Translation, nil
	case "scaling":
		return rt.Scaling, nil
	}
	return nil, fmt.Errorf("no transform matching %s", s)
}

func transformFromString(s string) (*rt.Matrix, error) {
	re, err := regexp.Compile(`^(\w+)\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`)
	if err != nil {
		return nil, err
	}

	matches := re.FindStringSubmatch(s)
	if len(matches) < 5 {
		return nil, fmt.Errorf("failed to match %s to transform pattern", s)
	}
	transformation, err := selectTransform(matches[1])
	if err != nil {
		return nil, err
	}

	x, y, z, err := parseXYZ(matches[2], matches[3], matches[4])
	if err != nil {
		return nil, err
	}

	return transformation(x, y, z), nil
}

func InitializeTransformationScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w+) ← translation\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) context.Context {
			return setMatrix(ctx, name, rt.Translation(x, y, z))
		})
	sc.Given(
		`^(\w+) ← scaling\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) context.Context {
			return setMatrix(ctx, name, rt.Scaling(x, y, z))
		})
	sc.Given(
		`^(\w+) ← shearing\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, xy, xz, yx, yz, zx, zy float64) context.Context {
			return setMatrix(ctx, name, rt.Shearing(xy, xz, yx, yz, zx, zy))
		})

	sc.Given(
		`^(\w+) ← rotation_x\(π / (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, f float64) context.Context {
			return setMatrix(ctx, name, rt.RotationX(math.Pi/f))
		})
	sc.Given(
		`^(\w+) ← rotation_y\(π / (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, f float64) context.Context {
			return setMatrix(ctx, name, rt.RotationY(math.Pi/f))
		})
	sc.Given(
		`^(\w+) ← rotation_z\(π / (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, f float64) context.Context {
			return setMatrix(ctx, name, rt.RotationZ(math.Pi/f))
		})
	sc.Given(
		`^m ← scaling\(1, 0.5, 1\) \* rotation_z\(π/5\)`,
		func(ctx context.Context) (context.Context, error) {
			scaling := rt.Scaling(1, 0.5, 1)
			rotation := rt.RotationZ(math.Pi / 5)
			m, err := scaling.MultiplyMatrix(rotation)
			if err != nil {
				return ctx, err
			}
			return setMatrix(ctx, "m", m), nil
		})
	sc.Given(
		`^(\w+) is added to (\w+)$`,
		func(ctx context.Context, source, dest string) (context.Context, error) {
			shape, err := getShape(ctx, source)
			if err != nil {
				return ctx, err
			}
			world, err := getWorld(ctx, dest)
			if err != nil {
				return ctx, nil
			}

			world.Add(shape)
			return ctx, nil
		})

	sc.When(
		`^(\w+) ← transform\((\w+), (\w+)\)$`,
		func(ctx context.Context, destName, rayName, matrixName string) (context.Context, error) {
			ray, err := getRay(ctx, rayName)
			if err != nil {
				return ctx, err
			}

			matrix, err := getMatrix(ctx, matrixName)
			if err != nil {
				return ctx, err
			}

			ray, err = ray.Transform(matrix)
			if err != nil {
				return ctx, fmt.Errorf("Could not transform ray: %w", err)
			}
			return setRay(ctx, destName, ray), nil
		})

	sc.When(
		`^(\w) ← view_transform\((\w+), (\w+), (\w+)\)$`,
		func(ctx context.Context, dest, fromName, toName, upName string) (context.Context, error) {
			from, err := getTuple(ctx, fromName)
			if err != nil {
				return ctx, nil
			}
			to, err := getTuple(ctx, toName)
			if err != nil {
				return ctx, nil
			}
			up, err := getTuple(ctx, upName)
			if err != nil {
				return ctx, nil
			}

			m, err := rt.ViewTransformation(from, to, up)
			if err != nil {
				return ctx, nil
			}

			return setMatrix(ctx, dest, m), nil
		})
	sc.When(
		`(\w) ← (\w) \* (\w) \* (\w)$`,
		func(ctx context.Context, dest, name1, name2, name3 string) (context.Context, error) {
			matrix1, err := getMatrix(ctx, name1)
			if err != nil {
				return ctx, err
			}
			matrix2, err := getMatrix(ctx, name2)
			if err != nil {
				return ctx, err
			}
			matrix3, err := getMatrix(ctx, name3)
			if err != nil {
				return ctx, err
			}

			matrix, err := matrix1.MultiplyMatrix(matrix2)
			if err != nil {
				return ctx, fmt.Errorf("cannot multiply %s by %s", matrix1, matrix2)
			}
			matrix, err = matrix.MultiplyMatrix(matrix3)
			if err != nil {
				return ctx, fmt.Errorf("cannot multiply %s by %s", matrix, matrix3)
			}
			return setMatrix(ctx, dest, matrix), nil
		})

	sc.Then(
		`^(\w) = translation\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			t := rt.Translation(x, y, z)
			if !matrix.Equal(t) {
				return fmt.Errorf("expected %s to equal %s, but it did not", matrix, t)
			}
			return nil
		})
	sc.Then(
		`^(\w) = scaling\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			t := rt.Scaling(x, y, z)
			if !matrix.Equal(t) {
				return fmt.Errorf("expected %s to equal %s, but it did not", matrix, t)
			}
			return nil
		})

	sc.Then(
		`^(\w+) \* (\w) = (\w)$`,
		func(ctx context.Context, matrixName, tupleName string) error {
			matrix, err := getMatrix(ctx, matrixName)
			if err != nil {return err}

			tuple, err := getTuple(ctx, tupleName)
			if err != nil {return err}
			
			result, err := matrix.MultiplyTuple(tuple)
			if err != nil {return nil}

			if !result.Equal(tuple) {
				return fmt.Errorf("expected %s, got %s", tuple, result)
			}
			return nil
		})
}
