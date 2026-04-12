package test

import (
	"context"
	"fmt"
	"math"
	"regexp"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

// 2:A, 2.75:B, 3.25:C, 4.75:B, 5.25:C, 6:A
func parseIntersections(ctx context.Context, list string) (*rt.Intersections, error) {
	re, err := regexp.Compile(`(-?√?\d\S*\d?):(\w+)`)
	if err != nil {
		return nil, err
	}

	matches := re.FindAllStringSubmatch(list, -1)
	if matches == nil {
		return nil, fmt.Errorf("invalid intersections list: %s", list)
	}

	intersections := make([]*rt.Intersection, 0)
	for _, match := range matches {
		t, err := parseFloat(match[1])
		if err != nil {
			return nil, err
		}
		shape, err := getShape(ctx, match[2])
		if err != nil {
			return nil, err
		}

		intersection := rt.NewIntersection(t, shape)
		intersections = append(intersections, intersection)
	}
	return rt.NewIntersections(intersections...), nil
}

func createIntersections(ctx context.Context, dest, name string) (context.Context, error) {
	intersection, err := getIntersection(ctx, name)
	if err != nil {
		return ctx, err
	}

	intersections := rt.NewIntersections(intersection)

	return setIntersections(ctx, dest, intersections), nil
}

func InitializeIntersectionScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w+) ← intersection\((\S+), (\w+)\)$`,
		func(ctx context.Context, dest string, ts, source string) (context.Context, error) {
			shape, err := getShape(ctx, source)
			if err != nil {
				return ctx, err
			}

			t, err := parseFloat(ts)
			if err != nil {
				return ctx, err
			}
			return setIntersection(ctx, dest, rt.NewIntersection(t, shape)), nil
		})
	sc.Given(`^(\w+) ← intersections\((\w+)\)$`, createIntersections)
	sc.Given(
		`^(\w+) ← intersections\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, name1, name2 string) (context.Context, error) {
			intersection1, err := getIntersection(ctx, name1)
			if err != nil {
				return ctx, err
			}
			intersection2, err := getIntersection(ctx, name2)
			if err != nil {
				return ctx, err
			}

			intersections := rt.NewIntersections(intersection1, intersection2)

			return setIntersections(ctx, dest, intersections), nil
		})
	sc.Given(
		`^(\w+) ← intersections\((\w+), (\w+), (\w+), (\w+)\)$`,
		func(ctx context.Context, dest, name1, name2, name3, name4 string) (context.Context, error) {
			intersection1, err := getIntersection(ctx, name1)
			if err != nil {
				return ctx, err
			}
			intersection2, err := getIntersection(ctx, name2)
			if err != nil {
				return ctx, err
			}
			intersection3, err := getIntersection(ctx, name3)
			if err != nil {
				return ctx, err
			}
			intersection4, err := getIntersection(ctx, name4)
			if err != nil {
				return ctx, err
			}

			intersections := rt.NewIntersections(intersection1, intersection2, intersection3, intersection4)

			return setIntersections(ctx, dest, intersections), nil
		})

	sc.Given(
		`^(\w+) ← intersections\((-?√?\d\S*\d?:.+)\)$`,
		func(ctx context.Context, dest, list string) (context.Context, error) {
			intersections, err := parseIntersections(ctx, list)
			if err != nil {
				return ctx, err
			}
			return setIntersections(ctx, dest, intersections), nil
		})

	sc.When(`^(\w+) ← intersections\((\w+)\)$`, createIntersections)
	sc.When(
		`^(\w) ← intersection\((-?\d+\.?\d*), (\w+)\)$`,
		func(ctx context.Context, dest string, t float64, source string) (context.Context, error) {
			shape, err := getShape(ctx, source)
			if err != nil {
				return ctx, err
			}

			return setIntersection(ctx, dest, rt.NewIntersection(t, shape)), nil
		})
	sc.When(
		`^(\w+) ← intersections\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, name1, name2 string) (context.Context, error) {
			intersection1, err := getIntersection(ctx, name1)
			if err != nil {
				return ctx, err
			}

			intersection2, err := getIntersection(ctx, name2)
			if err != nil {
				return ctx, err
			}

			intersections := rt.NewIntersections(intersection1, intersection2)

			return setIntersections(ctx, dest, intersections), nil
		})
	sc.When(
		`(\w+) ← hit\((\w+)\)`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			intersections, err := getIntersections(ctx, source)
			if err != nil {
				return ctx, err
			}

			hit := intersections.Hit()
			return setIntersection(ctx, dest, hit), nil
		})
	sc.When(
		`^(\w+) ← prepare_computations\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, intersectionName, rayName string) (context.Context, error) {
			intersection, err := getIntersection(ctx, intersectionName)
			if err != nil {
				return ctx, nil
			}
			ray, err := getRay(ctx, rayName)
			if err != nil {
				return ctx, nil
			}

			comps, err := intersection.PrepareComputations(ray, nil)
			if err != nil {
				return ctx, nil
			}

			return setComps(ctx, dest, comps), nil
		})
	sc.When(
		`^(\w+) ← prepare_computations\((\w+)\[(\d+)\], (\w+), (\w+)\)$`,
		func(ctx context.Context, dest, intersectionsName string, idx int, rayName, _ string) (context.Context, error) {
			intersections, err := getIntersections(ctx, intersectionsName)
			if err != nil {
				return ctx, err
			}

			intersection := intersections.At(idx)
			ray, err := getRay(ctx, rayName)
			if err != nil {
				return ctx, err
			}

			comps, err := intersection.PrepareComputations(ray, intersections)
			if err != nil {
				return ctx, err
			}

			return setComps(ctx, dest, comps), nil
		})
	sc.When(
		`^(\w+) ← prepare_computations\((\w+), (\w+), (\w+)\)$`,
		func(ctx context.Context, dest, intersectionName, rayName, intersectionsName string) (context.Context, error) {
			intersection, err := getIntersection(ctx, intersectionName)
			if err != nil {
				return ctx, err
			}
			ray, err := getRay(ctx, rayName)
			if err != nil {
				return ctx, err
			}
			intersections, err := getIntersections(ctx, intersectionsName)
			if err != nil {
				return ctx, err
			}

			comps, err := intersection.PrepareComputations(ray, intersections)
			if err != nil {
				return ctx, err
			}

			return setComps(ctx, dest, comps), nil
		})
	sc.When(
		`^reflectance ← schlick\(comps\)$`,
		func(ctx context.Context) (context.Context, error) {
			comps, err := getComps(ctx, "comps")
			if err != nil {
				return ctx, err
			}

			reflectance := comps.Schlick()
			return setFloat(ctx, "reflectance", reflectance), nil
		})
	sc.When(
		`^(\w+) ← intersection_with_uv\((-?\d+\.?\d*), (\w+), (-?\d+\.?\d*), (-?\d+\.?\d*)\)`,
		func(ctx context.Context, dest string, t float64, shapeName string, u, v float64) (context.Context, error) {
			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return ctx, err
			}

			i := rt.NewIntersectionWithUV(t, shape, u, v)
			return setIntersection(ctx, dest, i), nil
		})

	sc.Then(
		`^(\w).t = (-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, value float64) error {
			intersection, err := getIntersection(ctx, name)
			if err != nil {
				return err
			}

			if intersection.T != value {
				return fmt.Errorf("expected %f, got %f", value, intersection.T)
			}
			return nil
		})
	sc.Then(
		`^(\w+).object = (\w+)$`,
		func(ctx context.Context, name string, value string) error {
			intersection, err := getIntersection(ctx, name)
			if err != nil {
				return err
			}

			object, err := getShape(ctx, value)
			if err != nil {
				return err
			}

			if !intersection.Object.Equal(object) {
				return fmt.Errorf("expected %s, got %s", object, intersection.Object)
			}
			return nil
		})
	sc.Then(
		`^(\w+).count = (\d+)$`,
		func(ctx context.Context, name string, c int) error {
			intersections, err := getIntersections(ctx, name)
			if err != nil {
				return err
			}

			if intersections.Len() != c {
				return fmt.Errorf("expected %d, got %d", c, intersections.Len())
			}
			return nil
		})
	sc.Then(
		`^(\w+)\[(\d+)\].t = (-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, idx int, val float64) error {
			intersections, err := getIntersections(ctx, name)
			if err != nil {
				return err
			}

			t := intersections.At(idx).T
			if math.Abs(t-val) > rt.EPSILON {
				return fmt.Errorf("expected %f, got %f", val, t)
			}
			return nil
		})
	sc.Then(
		`^(\w+)\[(\d+)\].object = (\w+)$`,
		func(ctx context.Context, name string, idx int, value string) error {
			intersections, err := getIntersections(ctx, name)
			if err != nil {
				return err
			}

			object, err := getShape(ctx, value)
			if err != nil {
				return err
			}

			expect := intersections.At(idx).Object
			if !expect.Equal(object) {
				return fmt.Errorf("expected %s, got %s", object, expect)
			}
			return nil
		})

	sc.Then(
		`^(\w+)\[(\d+)\].u = (-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, idx int, val float64) error {
			intersections, err := getIntersections(ctx, name)
			if err != nil {
				return err
			}

			u := intersections.At(idx).U
			if math.Abs(u-val) > rt.EPSILON {
				return fmt.Errorf("expected %f, got %f", val, u)
			}
			return nil
		})
	sc.Then(
		`^(\w+)\[(\d+)\].v = (-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, idx int, val float64) error {
			intersections, err := getIntersections(ctx, name)
			if err != nil {
				return err
			}

			v := intersections.At(idx).V
			if math.Abs(v-val) > rt.EPSILON {
				return fmt.Errorf("expected %f, got %f", val, v)
			}
			return nil
		})
	sc.Then(
		`^(\w) = (\w{1,2})$`,
		func(ctx context.Context, name1, name2 string) error {
			intersection1, err := getValue[rt.Equality](ctx, name1)
			if err != nil {
				return err
			}

			intersection2, err := getValue[rt.Equality](ctx, name2)
			if err != nil {
				return err
			}

			if !intersection1.Equal(intersection2) {
				return fmt.Errorf("expected %s got %s", intersection2, intersection1)
			}
			return nil
		})
	sc.Then(
		`^(\w) is nothing$`,
		func(ctx context.Context, name string) error {
			intersection, err := getIntersection(ctx, name)
			if err != nil {
				return err
			}

			if intersection != nil {
				return fmt.Errorf("expected nil value, got %s", intersection)
			}
			return nil
		})
	sc.Then(
		`^xs is empty$`,
		func(ctx context.Context) error {
			intersections, err := getIntersections(ctx, "xs")
			if err != nil {
				return err
			}

			if intersections.Len() > 0 {
				return fmt.Errorf("expected %s to be empty", "xs")
			}
			return nil
		})

	sc.Then(
		`^(\w+).t = (\w+).t$`,
		func(ctx context.Context, dest, source string) error {
			comps, err := getComps(ctx, dest)
			if err != nil {
				return err
			}
			intersection, err := getIntersection(ctx, source)
			if err != nil {
				return err
			}

			if comps.T != intersection.T {
				return fmt.Errorf("expected %f, got %f", intersection.T, comps.T)
			}
			return nil
		})
	sc.Then(
		`^(\w+).object = (\w+).object$`,
		func(ctx context.Context, dest, source string) error {
			comps, err := getComps(ctx, dest)
			if err != nil {
				return err
			}
			intersection, err := getIntersection(ctx, source)
			if err != nil {
				return err
			}

			if !comps.Object.Equal(intersection.Object) {
				return fmt.Errorf("expected %s, got %s", intersection.Object, comps.Object)
			}
			return nil
		})
	sc.Then(
		`^(\w+).point = point\((-?\d+\.?\d*), (-?\d+\.?\d*), (-?\d+\.?\d*)\)$`,
		func(ctx context.Context, dest string, x, y, z float64) error {
			comps, err := getComps(ctx, dest)
			if err != nil {
				return err
			}

			point := rt.NewPoint(x, y, z)

			if !comps.Point.Equal(point) {
				return fmt.Errorf("expected %s, got %s", point, comps.Point)
			}
			return nil
		})
	sc.Then(
		`^(\w+).eyev = vector\((-?\d+\.?\d*), (-?\d+\.?\d*), (-?\d+\.?\d*)\)$`,
		func(ctx context.Context, dest string, x, y, z float64) error {
			comps, err := getComps(ctx, dest)
			if err != nil {
				return err
			}

			vector := rt.NewVector(x, y, z)

			if !comps.Eye.Equal(vector) {
				return fmt.Errorf("expected %s, got %s", vector, comps.Eye)
			}
			return nil
		})
	sc.Then(
		`^(\w+).normalv = vector\((-?\d+\.?\d*), (-?\d+\.?\d*), (-?\d+\.?\d*)\)$`,
		func(ctx context.Context, dest string, x, y, z float64) error {
			comps, err := getComps(ctx, dest)
			if err != nil {
				return err
			}

			vector := rt.NewVector(x, y, z)

			if !comps.Normal.Equal(vector) {
				return fmt.Errorf("expected %s, got %s", vector, comps.Normal)
			}
			return nil
		})

	sc.Then(
		`^(\w+).inside = (\w+)$`,
		func(ctx context.Context, dest, val string) error {
			comps, err := getComps(ctx, dest)
			if err != nil {
				return err
			}
			expect := (val == "true")

			if comps.Inside != expect {
				return fmt.Errorf("expected %v, got %v", expect, comps.Inside)
			}
			return nil
		})

	sc.Then(
		`^comps.over_point.z < -EPSILON/2$`,
		func(ctx context.Context) error {
			comps, err := getComps(ctx, "comps")
			if err != nil {
				return err
			}

			if !(comps.OverPoint.Z < -rt.EPSILON/2) {
				return fmt.Errorf("expected comps.over_point.z < -EPSILON/2")
			}
			return nil
		})
	sc.Then(
		`^comps.under_point.z > EPSILON/2$`,
		func(ctx context.Context) error {
			comps, err := getComps(ctx, "comps")
			if err != nil {
				return err
			}

			if !(comps.UnderPoint.Z > -rt.EPSILON/2) {
				return fmt.Errorf("expected comps.under_point.z > -EPSILON/2")
			}
			return nil
		})
	sc.Then(
		`^comps.point.z > comps.over_point.z$`,
		func(ctx context.Context) error {
			comps, err := getComps(ctx, "comps")
			if err != nil {
				return err
			}

			if !(comps.Point.Z > comps.OverPoint.Z) {
				return fmt.Errorf("expected comps.point.z > comps.over_point.z")
			}
			return nil
		})
	sc.Then(
		`^comps.point.z < comps.under_point.z$`,
		func(ctx context.Context) error {
			comps, err := getComps(ctx, "comps")
			if err != nil {
				return err
			}

			if !(comps.Point.Z < comps.UnderPoint.Z) {
				return fmt.Errorf("expected comps.point.z < comps.under_point.z")
			}
			return nil
		})
	sc.Then(
		`^comps.reflectv = vector\((\S+), (\S+), (\S+)\)$`,
		func(ctx context.Context, x, y, z string) error {
			comps, err := getComps(ctx, "comps")
			if err != nil {
				return err
			}

			v, err := vectorFromStrings(x, y, z)
			if err != nil {
				return err
			}

			if !comps.Reflect.Equal(v) {
				return fmt.Errorf("expected %s, got %s", v, comps.Reflect)
			}
			return nil
		})

	sc.Then(
		`^comps.n1 = (-?\d+\.?\d*)$`,
		func(ctx context.Context, val float64) error {
			comps, err := getComps(ctx, "comps")
			if err != nil {
				return err
			}

			if comps.N1 != val {
				return fmt.Errorf("expected %f, got %f", val, comps.N1)
			}
			return nil
		})
	sc.Then(
		`^comps.n2 = (-?\d+\.?\d*)$`,
		func(ctx context.Context, val float64) error {
			comps, err := getComps(ctx, "comps")
			if err != nil {
				return err
			}

			if comps.N2 != val {
				return fmt.Errorf("expected %f, got %f", val, comps.N2)
			}
			return nil
		})

	sc.Then(
		`^reflectance = (-?\d+\.?\d*)$`,
		func(ctx context.Context, val float64) error {
			reflectance, err := getFloat(ctx, "reflectance")
			if err != nil {
				return err
			}

			if math.Abs(reflectance-val) > rt.EPSILON {
				return fmt.Errorf("expected %f, got %f", val, reflectance)
			}
			return nil
		})

	sc.Then(
		`^(\w).u = (-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, val float64) error {
			intersection, err := getIntersection(ctx, dest)
			if err != nil {
				return err
			}

			if intersection.U != val {
				return fmt.Errorf("expected %f, got %f", val, intersection.U)
			}
			return nil
		})

	sc.Then(
		`^(\w).v = (-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, val float64) error {
			intersection, err := getIntersection(ctx, dest)
			if err != nil {
				return err
			}

			if intersection.V != val {
				return fmt.Errorf("expected %f, got %f", val, intersection.V)
			}
			return nil
		})

}
