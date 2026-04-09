package test

import (
	"context"
	"fmt"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func InitializeSmoothTrianglesScenario(sc *godog.ScenarioContext) {
	sc.When(
		`^(\w+) ← smooth_triangle\((\w+), (\w+), (\w+), (\w+), (\w+), (\w+)\)$`,
		func (ctx context.Context, dest, p1n, p2n, p3n, n1n, n2n, n3n string) (context.Context, error) {
			p1, err := getTuple(ctx, p1n)
			if err != nil {return ctx, err}
			p2, err := getTuple(ctx, p2n)
			if err != nil {return ctx, err}
			p3, err := getTuple(ctx, p3n)
			if err != nil {return ctx, err}
			
			n1, err := getTuple(ctx, n1n)
			if err != nil {return ctx, err}
			n2, err := getTuple(ctx, n2n)
			if err != nil {return ctx, err}
			n3, err := getTuple(ctx, n3n)
			if err != nil {return ctx, err}

			return setShape(ctx, dest, rt.NewSmoothTriangle(p1,p2,p3,n1,n2,n3)), nil
		})
	sc.When(
		`^(\w+) ← normal_at\((\w+), point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\), (\w+)\)$`,
		func (ctx context.Context, dest, shapeName string, x, y, z float64, intersectionName string) (context.Context, error) {
			shape, err := getShape(ctx, shapeName)
			if err != nil {return ctx, err}

			intersection, err := getIntersection(ctx, intersectionName)
			if err != nil {return ctx, err}

			point := rt.NewPoint(x,y,z)

			normal, err := shape.NormalAt(point, intersection)
			if err != nil {return ctx, err}

			return setTuple(ctx, dest, normal), nil
		})

	sc.Then(
		`^(\w+).n1 = (\w+)$`,
		func (ctx context.Context, dest, source string) error {
			shape, err := getShape(ctx, dest)
			if err != nil { return err }
			
			normal, err := getTuple(ctx, source)
			if err != nil { return err }

			trait, ok := shape.Trait.(*rt.SmoothTriangle)
			if !ok {return fmt.Errorf("%s is not a shape with a SmoothTriangle", dest)}

			if !trait.N1.Equal(normal) {
				return fmt.Errorf("expected %s, got %s", normal, trait.N1)
			}
			return nil
		})
	sc.Then(
		`^(\w+).n2 = (\w+)$`,
		func (ctx context.Context, dest, source string) error {
			shape, err := getShape(ctx, dest)
			if err != nil { return err }
			
			normal, err := getTuple(ctx, source)
			if err != nil { return err }

			trait, ok := shape.Trait.(*rt.SmoothTriangle)
			if !ok {return fmt.Errorf("%s is not a shape with a SmoothTriangle", dest)}

			if !trait.N2.Equal(normal) {
				return fmt.Errorf("expected %s, got %s", normal, trait.N2)
			}
			return nil
		})
	sc.Then(
		`^(\w+).n3 = (\w+)$`,
		func (ctx context.Context, dest, source string) error {
			shape, err := getShape(ctx, dest)
			if err != nil { return err }
			
			normal, err := getTuple(ctx, source)
			if err != nil { return err }

			trait, ok := shape.Trait.(*rt.SmoothTriangle)
			if !ok {return fmt.Errorf("%s is not a shape with a SmoothTriangle", dest)}

			if !trait.N3.Equal(normal) {
				return fmt.Errorf("expected %s, got %s", normal, trait.N3)
			}
			return nil
		})
}