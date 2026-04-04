package test

import (
	"fmt"
	"math"
	"context"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func InitializeCylinderScenario(sc *godog.ScenarioContext) {
	sc.Given (
		`^(\w+) ← cylinder\(\)$`,
		func(ctx context.Context, name string) context.Context {
			cylinder := rt.NewCylinder()
			return setShape(ctx, name, cylinder)
		})
	sc.Given(
		`^(\w+) ← normalize\(vector\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		func(ctx context.Context, name string, x, y, z float64) context.Context {
			vector := rt.NewVector(x, y, z)
			normalized := vector.Normalize()
			return setTuple(ctx, name, normalized)
		})
	sc.Given(
		`^(\w+) ← ray\(point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\), (\w+)\)$`,
		func(ctx context.Context, dest string, x, y, z float64, vectorName string) (context.Context, error) {
			vector, err := getTuple(ctx, vectorName)
			if err != nil {return ctx, err}
			point := rt.NewPoint(x, y, z)
			return setRay(ctx, dest, rt.NewRay(point, vector)), nil
		})
	sc.Given (
		`(\w+).minimum ← (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, val float64) (context.Context, error) {
			shape, err := getShape(ctx, name)
			if err != nil {return ctx, err}
			shape.Trait.AsCapped().Minimum = val
			return ctx, nil
		})
	sc.Given (
		`(\w+).maximum ← (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, val float64) (context.Context, error) {
			shape, err := getShape(ctx, name)
			if err != nil {return ctx, err}
			shape.Trait.AsCapped().Maximum = val
			return ctx, nil
		})
	sc.Given (
		`(\w+).closed ← (true|false)$`,
		func(ctx context.Context, name string, val string) (context.Context, error) {
			shape, err := getShape(ctx, name)
			if err != nil {return ctx, err}
			shape.Trait.AsCapped().Closed = val == "true"
			return ctx, nil
		})

	sc.When(
		`(\w+) ← local_normal_at\((\w+), (\w+)\)`,
		func(ctx context.Context, dest, source, pointName string) (context.Context, error) {
			shape, err := getShape(ctx, source)
			if err != nil {return ctx, err}

			point, err := getTuple(ctx, pointName)
			if err != nil {return ctx, err}

			normal, err := shape.Trait.LocalNormalAt(point)
			if err != nil {return ctx, err}
			return setTuple(ctx, dest, normal), nil
		})

	sc.Then(
		`^(\w+).minimum = -infinity$`,
		func(ctx context.Context, name string) error {
			shape, err := getShape(ctx, name)
			if err != nil {return err}
			if shape.Trait.(*rt.Cylinder).Minimum != math.Inf(-1) {
				return fmt.Errorf("expected minimum to be -infinity, got %f", shape.Trait.AsCapped().Minimum)
			}
			return nil
		})

    sc.Then(
		`^(\w+).maximum = infinity$`,
		func(ctx context.Context, name string) error {
			shape, err := getShape(ctx, name)
			if err != nil {return err}
			if shape.Trait.(*rt.Cylinder).Maximum != math.Inf(1) {
				return fmt.Errorf("expected maximum to be infinity, got %f", shape.Trait.AsCapped().Maximum)
			}
			return nil
		})
	sc.Then(
		`^(\w+).closed = false$`,
		func(ctx context.Context, name string) error {
			shape, err := getShape(ctx, name)
			if err != nil {return err}
			if shape.Trait.(*rt.Cylinder).Closed != false {
				return fmt.Errorf("expected closed to be false, got %t", shape.Trait.AsCapped().Closed)
			}
			return nil
		})
}