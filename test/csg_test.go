package test

import (
	"context"
	"fmt"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func InitializeCSGScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w) ← csg\("(\w+)", (\w+), (\w+)\)$`,
		func(ctx context.Context, dest, operation, s1n, s2n string) (context.Context, error) {
			s1, err := getShape(ctx, s1n)
			if err != nil {
				return ctx, err
			}
			s2, err := getShape(ctx, s2n)
			if err != nil {
				return ctx, err
			}

			return setShape(ctx, dest, rt.NewCSG(operation, s1, s2)), nil
		})
	sc.Given(
		`^(\w+) ← csg\("union", sphere\(\), cube\(\)\)$`,
		func(ctx context.Context, dest string) context.Context {
			return setShape(ctx, dest, rt.NewCSG("union", rt.NewSphere(), rt.NewCube()))
		})

	sc.When(
		`^(\w+) ← csg\("(\w+)", (\w+), (\w+)\)$`,
		func(ctx context.Context, dest, operation, s1, s2 string) (context.Context, error) {
			left, err := getShape(ctx, s1)
			if err != nil {
				return ctx, err
			}

			right, err := getShape(ctx, s2)
			if err != nil {
				return ctx, err
			}

			csg := rt.NewCSG(operation, left, right)

			return setShape(ctx, dest, csg), nil
		})

	sc.When(
		`^result ← intersection_allowed\("(\w+)", (\w+), (\w+), (\w+)\)$`,
		func(ctx context.Context, op string, lhit, inl, inr string) context.Context {
			result := rt.IntersectionAllowed(op, lhit == "true", inl == "true", inr == "true")
			return setBoolean(ctx, "result", result)
		})

	sc.When(
		`^(\w+) ← filter_intersections\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, shapeName, intersectionsName string) (context.Context, error) {
			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return ctx, nil
			}

			intersections, err := getIntersections(ctx, intersectionsName)
			if err != nil {
				return ctx, nil
			}

			trait, ok := shape.Trait.(*rt.CSG)
			if !ok {
				return ctx, fmt.Errorf("Wanted CSG but got %s", shape)
			}

			result := trait.FilterIntersections(intersections)
			return setIntersections(ctx, dest, result), nil
		})

	sc.Then(
		`^result = (\w+)$`,
		func(ctx context.Context, value string) error {
			result, err := getBoolean(ctx, "result")
			if err != nil {
				return err
			}

			if result != (value == "true") {
				return fmt.Errorf("expected %v, get %v", value, result)
			}
			return nil
		})

	sc.Then(
		`^(\w+)\[(\d+)\] = (\w+)\[(\d+)\]$`,
		func(ctx context.Context, dest string, didx int, source string, sidx int) error {
			destIntersection, err := getIntersections(ctx, dest)
			if err != nil {
				return err
			}

			sourceIntersections, err := getIntersections(ctx, source)
			if err != nil {
				return err
			}

			left := destIntersection.At(didx)
			right := sourceIntersections.At(sidx)

			if !left.Equal(right) {
				return fmt.Errorf("expected %s, got %s", left, right)
			}
			return nil
		})
}
