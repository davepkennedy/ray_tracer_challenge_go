package test

import (
	"context"
	"fmt"
	"log"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func createRay(ctx context.Context, name string, x1, y1, z1, x2, y2, z2 string) (context.Context, error) {
	point, err := pointFromStrings(x1, y1, z1)
	if err != nil {
		return ctx, err
	}
	vector, err := vectorFromStrings(x2, y2, z2)
	if err != nil {
		return ctx, err
	}
	return setRay(ctx, name, rt.NewRay(point, vector)), nil
}

func InitializeRaysScenario(sc *godog.ScenarioContext) {

	sc.Given(
		`^(\w+) ← ray\(point\((\S+), (\S+), (\S+)\), vector\((\S+), (\S+), (\S+)\)\)$`,
		createRay)

	sc.When(
		`^(\w+) ← ray\(point\((\S+), (\S+), (\S+)\), vector\((\S+), (\S+), (\S+)\)\)$`,
		createRay)

	sc.When(
		`^(\w+) ← ray\((\w+), (\w+)\)$`,
		func(ctx context.Context, name string, origin, direction string) (context.Context, error) {
			point, err := getTuple(ctx, origin)
			if err != nil {
				return ctx, err
			}

			vector, err := getTuple(ctx, direction)
			if err != nil {
				return ctx, err
			}

			ray := rt.NewRay(point, vector)
			return setRay(ctx, name, ray), nil
		})

	sc.Then(
		`^(\w+).origin = (\w+)$`,
		func(ctx context.Context, name, origin string) error {
			point, err := getTuple(ctx, origin)
			if err != nil {
				return err
			}

			ray, err := getRay(ctx, name)
			if err != nil {
				return err
			}

			if !ray.Origin.Equal(point) {
				return fmt.Errorf("expected %s, got %s", point, ray.Origin)
			}
			return nil
		})
	sc.Then(
		`^(\w+).origin = point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) error {
			ray, err := getRay(ctx, name)
			if err != nil {
				return err
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

			vector, err := getTuple(ctx, direction)
			if err != nil {
				return err
			}

			ray, err := getRay(ctx, name)
			if err != nil {
				return err
			}

			if !ray.Direction.Equal(vector) {
				return fmt.Errorf("expected %s, got %s", vector, ray.Direction)
			}
			return nil
		})
	sc.Then(
		`^(\w+).direction = vector\((\S+), (\S+), (\S+)\)$`,
		func(ctx context.Context, name, x, y, z string) error {
			ray, err := getRay(ctx, name)
			if err != nil {
				return err
			}

			vector, err := vectorFromStrings(x,y,z)
			log.Printf("Passed comparison vector %s", vector)
			if err != nil {return err}

			if !ray.Direction.Equal(vector) {
				return fmt.Errorf("expected %s, got %s", vector, ray.Direction)
			}
			return nil
		})

	sc.Then(
		`^position\((\w+), (\-?\d+\.?\d*)\) = point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, t float64, x, y, z float64) error {
			ray, err := getRay(ctx, name)
			if err != nil {
				return err
			}

			point := rt.NewPoint(x, y, z)
			position := ray.Position(t)
			if !position.Equal(point) {
				return fmt.Errorf("expected %s, got %s", point, position)
			}
			return nil
		})
}
