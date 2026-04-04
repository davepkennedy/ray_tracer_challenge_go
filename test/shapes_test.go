package test

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

type TestShape struct{}

var (
	savedRay *rt.Ray
)

func NewTestShape() *rt.Shape {
	return rt.NewShape(TestShape{})
}

func selectPattern(name string) (*rt.Pattern, error) {
	switch name {
	case "test_pattern()":
		return newTestPattern(), nil
	}
	return nil, fmt.Errorf("unknown pattern %s", name)
}

func InitializeShapesScenario(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		savedRay = nil
		return ctx, nil
	})

	sc.Given(
		`^(\w+) ← test_shape\(\)$`,
		func(ctx context.Context, dest string) context.Context {
			return setShape(ctx, dest, NewTestShape())
		})
	sc.Given(
		`^set_transform\((\w+), scaling\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		func(ctx context.Context, target string, x, y, z float64) (context.Context, error) {
			shape, err := getShape(ctx, target)
			if err != nil {
				return ctx, err
			}

			pt := rt.Scaling(x, y, z)
			shape.Transform = pt
			return ctx, nil
		})
	sc.Given(
		`^(\w+) has:$`,
		func(ctx context.Context, dest string, table *godog.Table) (context.Context, error) {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return ctx, err
			}

			for _, row := range table.Rows {
				switch row.Cells[0].Value {
				case "material.transparency":
					t, err := strconv.ParseFloat(row.Cells[1].Value, 64)
					if err != nil {
						return ctx, err
					}
					shape.Material.Transparency = t
				case "material.refractive_index":
					t, err := strconv.ParseFloat(row.Cells[1].Value, 64)
					if err != nil {
						return ctx, err
					}
					shape.Material.RefractiveIndex = t
				case "material.ambient":
					t, err := strconv.ParseFloat(row.Cells[1].Value, 64)
					if err != nil {
						return ctx, err
					}
					shape.Material.Ambient = t
				case "material.pattern":
					pattern, err := selectPattern(row.Cells[1].Value)
					if err != nil {
						return ctx, err
					}
					shape.Material.Pattern = pattern
				default:
					return ctx, fmt.Errorf("unknown property %s", row.Cells[0].Value)
				}
			}
			return ctx, nil
		})

	sc.When(
		`^(\w+) ← local_intersect\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, shapeName, rayName string) (context.Context, error) {
			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return ctx, err
			}

			ray, err := getRay(ctx, rayName)
			if err != nil {
				return ctx, err
			}

			ts := shape.Trait.Intersect(ray)
			is := make([]*rt.Intersection, 0)
			for _, t := range ts {
				is = append(is, rt.NewIntersection(t, shape))
			}
			xs := rt.NewIntersections(is...)
			return setIntersections(ctx, dest, xs), nil
		})

	sc.Then(
		`^(\w).transform = translation\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, dest string, x, y, z float64) error {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return err
			}

			m := rt.Translation(x, y, z)
			if !shape.Transform.Equal(m) {
				return fmt.Errorf("expected %s, got %s", m, shape.Transform)
			}
			return nil
		})
	sc.Then(
		`^s.saved_ray.origin = point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, x, y, z float64) error {
			p := rt.NewPoint(x, y, z)
			if !p.Equal(savedRay.Origin) {
				return fmt.Errorf("expected %s, gpt %s", p, savedRay.Origin)
			}
			return nil
		})
	sc.Then(
		`^s.saved_ray.direction = vector\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, x, y, z float64) error {
			v := rt.NewVector(x, y, z)
			if !v.Equal(savedRay.Direction) {
				return fmt.Errorf("expected %s, gpt %s", v, savedRay.Direction)
			}
			return nil
		})
}

func (t TestShape) Equal(other rt.ShapeTrait) bool {
	return reflect.TypeOf(t) == reflect.TypeOf(other)
}

func (t TestShape) String() string {
	return "test_shape{}"
}

func (t TestShape) Intersect(ray *rt.Ray) []float64 {
	savedRay = ray
	return []float64{}
}

func (t TestShape) LocalNormalAt(point *rt.Tuple) (*rt.Tuple, error) {
	return point, nil
}
