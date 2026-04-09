package test

import (
	"context"
	"fmt"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

var (
	ErrNotTriangle = fmt.Errorf ("Shape Trait is not Triangular")
)

func getTriangle (shape *rt.Shape) (*rt.Triangle, error) {
	t, ok := shape.Trait.(*rt.Triangle)
	if ok {return t, nil}

	st, ok := shape.Trait.(*rt.SmoothTriangle)
	if !ok {return nil, ErrNotTriangle}

	return &st.Triangle, nil
}

func InitializeTrianglesScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w) ← triangle\((\w+), (\w+), (\w+)\)$`,
		func(ctx context.Context, dest, name1, name2, name3 string) (context.Context, error) {
			p1, err := getTuple(ctx, name1)
			if err != nil { return ctx, err }
			p2, err := getTuple(ctx, name2)
			if err != nil { return ctx, err }
			p3, err := getTuple(ctx, name3)
			if err != nil { return ctx, err }
			t := rt.NewTriangle(p1,p2,p3)
			return setShape(ctx, dest, t), nil
		})
	sc.Given(
		`^(\w) ← triangle\(point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\), point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\), point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		func(ctx context.Context, dest string, x1, y1, z1, x2, y2, z2, x3, y3, z3 float64) (context.Context, error) {
			p1 := rt.NewPoint(x1, y1, z1)
			p2 := rt.NewPoint(x2, y2, z2)
			p3 := rt.NewPoint(x3, y3, z3)
			t := rt.NewTriangle(p1,p2,p3)
			return setShape(ctx, dest, t), nil
		})

	sc.Then(
		`^(\w+).p1 = (\w+)$`,
		func(ctx context.Context, name, pointName string) error {
			shape, err := getShape(ctx, name)
			if err != nil { return err }
			point, err := getTuple(ctx, pointName)
			if err != nil { return err }

			trait, err := getTriangle(shape)
			if err != nil {return err}
			if !trait.P1.Equal(point) {
				return fmt.Errorf("expected %s.p1 to be %s, got %s", name, pointName, shape.Trait.(*rt.Triangle).P1)
			}
			return nil
		})	

	sc.Then(
		`^(\w+).p2 = (\w+)$`,
		func(ctx context.Context, name, pointName string) error {
			shape, err := getShape(ctx, name)
			if err != nil { return err }
			point, err := getTuple(ctx, pointName)
			if err != nil { return err }

			trait, err := getTriangle(shape)
			if err != nil {return err}
			if !trait.P2.Equal(point) {
				return fmt.Errorf("expected %s.p2 to be %s, got %s", name, pointName, shape.Trait.(*rt.Triangle).P2)
			}
			return nil
		})

	sc.Then(
		`^(\w+).p3 = (\w+)$`,
		func(ctx context.Context, name, pointName string) error {
			shape, err := getShape(ctx, name)
			if err != nil { return err }
			point, err := getTuple(ctx, pointName)
			if err != nil { return err }

			trait, err := getTriangle(shape)
			if err != nil {return err}
			if !trait.P3.Equal(point) {
				return fmt.Errorf("expected %s.p3 to be %s, got %s", name, pointName, shape.Trait.(*rt.Triangle).P3)
			}
			return nil
		})

	sc.Then(
		`^(\w).e1 = vector\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) error {
			shape, err := getShape(ctx, name)
			if err != nil { return err }
			edge := rt.NewVector(x, y, z)

			trait, err := getTriangle(shape)
			if err != nil {return err}
			if !trait.E1.Equal(edge) {
				return fmt.Errorf("expected %s.e1 to be %s, got %s", name, edge, shape.Trait.(*rt.Triangle).E1)
			}
			return nil
		})

	sc.Then(
		`^(\w).e2 = vector\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) error {
			shape, err := getShape(ctx, name)
			if err != nil { return err }
			edge := rt.NewVector(x, y, z)
			trait, err := getTriangle(shape)
			if err != nil {return err}
			if !trait.E2.Equal(edge) {
				return fmt.Errorf("expected %s.e2 to be %s, got %s", name, edge, shape.Trait.(*rt.Triangle).E2)
			}
			return nil
		})

	sc.Then(
		`^(\w).normal = vector\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) error {
			shape, err := getShape(ctx, name)
			if err != nil { return err }
			normal := rt.NewVector(x, y, z)

			trait, err := getTriangle(shape)
			if err != nil {return err}
			if !trait.Normal.Equal(normal) {
				return fmt.Errorf("expected %s.normal to be %s, got %s", name, normal, shape.Trait.(*rt.Triangle).Normal)
			}
			return nil
		})

	sc.Then(
		`^(\w+) = (\w).normal$`,
		func(ctx context.Context, dest, source string) error {
			shape, err := getShape(ctx, source)
			if err != nil { return err }
			
			normal, err := getTuple(ctx, dest)
			if err != nil { return err }

			trait, err := getTriangle(shape)
			if err != nil {return err}

			if !trait.Normal.Equal(normal) {
				return fmt.Errorf("expected %s to be %s, got %s", dest, normal, shape.Trait.(*rt.Triangle).Normal)
			}
			return nil
		})
}