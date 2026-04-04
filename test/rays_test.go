package test

import (
	"context"
	"fmt"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

type rayKey struct{ Name string }

func InitializeRaysScenario(sc *godog.ScenarioContext) {

	sc.Given(
		`^(\w+) ← ray\(point\((\S+), (\S+), (\S+)\), vector\((\S+), (\S+), (\S+)\)\)$`,
		func(ctx context.Context, name string, x1, y1, z1, x2, y2, z2 string) (context.Context, error) {
			point, err := pointFromStrings(x1, y1, z1)
			if err != nil {
				return ctx, err
			}
			vector, err := vectorFromStrings(x2, y2, z2)
			if err != nil {
				return ctx, err
			}
			return setRay(ctx, name, rt.NewRay(point, vector)), nil
		})

	sc.When(
		`^(\w+) ← ray\((\w+), (\w+)\)$`,
		func(ctx context.Context, name string, origin, direction string) (context.Context, error) {
			point, ok := ctx.Value(tupleKey{origin}).(*rt.Tuple)
			if !ok {
				return ctx, fmt.Errorf("no tuple named %s found", origin)
			}
			vector, ok := ctx.Value(tupleKey{direction}).(*rt.Tuple)
			if !ok {
				return ctx, fmt.Errorf("no tuple named %s found", direction)
			}
			ray := rt.NewRay(point, vector)
			return context.WithValue(ctx, rayKey{name}, ray), nil
		})

	sc.Then(
		`^(\w+).origin = (\w+)$`,
		func(ctx context.Context, name, origin string) error {
			point, ok := ctx.Value(tupleKey{origin}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", origin)
			}
			ray, ok := ctx.Value(rayKey{name}).(*rt.Ray)
			if !ok {
				return fmt.Errorf("no ray named %s found", name)
			}
			if !ray.Origin.Equal(point) {
				return fmt.Errorf("expected %s, got %s", point, ray.Origin)
			}
			return nil
		})
	sc.Then(
		`^(\w+).origin = point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) error {
			ray, ok := ctx.Value(rayKey{name}).(*rt.Ray)
			if !ok {
				return fmt.Errorf("no ray named %s found", name)
			}
			point := rt.NewPoint(x, y, z)
			if !ray.Origin.Equal(point) {
				return fmt.Errorf("expected %s, got %s", point, ray.Origin)
			}
			return nil
		})

	sc.Then(
		`^(\w+).direction = (\w+)$`,
		func(ctx context.Context, name, direction string) error {

			vector, ok := ctx.Value(tupleKey{direction}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", direction)
			}
			ray, ok := ctx.Value(rayKey{name}).(*rt.Ray)
			if !ok {
				return fmt.Errorf("no ray named %s found", name)
			}
			if !ray.Direction.Equal(vector) {
				return fmt.Errorf("expected %s, got %s", vector, ray.Direction)
			}
			return nil
		})
	sc.Then(
		`^(\w+).direction = vector\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) error {
			ray, ok := ctx.Value(rayKey{name}).(*rt.Ray)
			if !ok {
				return fmt.Errorf("no ray named %s found", name)
			}
			vector := rt.NewVector(x, y, z)
			if !ray.Direction.Equal(vector) {
				return fmt.Errorf("expected %s, got %s", vector, ray.Direction)
			}
			return nil
		})

	sc.Then(
		`^position\((\w+), (\-?\d+\.?\d*)\) = point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, t float64, x, y, z float64) error {
			ray, ok := ctx.Value(rayKey{name}).(*rt.Ray)
			if !ok {
				return fmt.Errorf("no ray named %s found", name)
			}
			point := rt.NewPoint(x, y, z)
			position := ray.Position(t)
			if !position.Equal(point) {
				return fmt.Errorf("expected %s, got %s", point, position)
			}
			return nil
		})
}
