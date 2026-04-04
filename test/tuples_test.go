package test

import (
	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

type tupleKey struct{ Name string }
type colorKey struct{ Name string }

func parseFloat(s string) (float64, error) {
	re1, err := regexp.Compile(`√(\-?\d+\.?\d*)/(\-?\d+\.?\d*)`)
	if err != nil {
		return 0, fmt.Errorf("failed to compile regex: %w", err)
	}
	matches := re1.FindStringSubmatch(s)
	if len(matches) >= 3 {
		n, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, err
		}
		d, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			return 0, err
		}
		result := math.Sqrt(n) / d
		if s[0] == '-' {
			result = -result
		}
		return result, nil
	}

	re2, err := regexp.Compile(`√(\-?\d+\.?\d*)`)
	if err != nil {
		return 0, fmt.Errorf("failed to compile regex: %w", err)
	}
	matches = re2.FindStringSubmatch(s)
	if len(matches) >= 2 {
		n, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return 0, err
		}
		return math.Sqrt(n), nil
	}

	return strconv.ParseFloat(s, 64)
}

func parseXYZ(x, y, z string) (xf, yf, zf float64, err error) {
	xf, err = parseFloat(x)
	if err != nil {
		return
	}
	yf, err = parseFloat(y)
	if err != nil {
		return
	}
	zf, err = parseFloat(z)
	if err != nil {
		return
	}
	return
}

func pointFromStrings(x, y, z string) (*rt.Tuple, error) {
	xf, yf, zf, err := parseXYZ(x, y, z)
	if err != nil {
		return nil, fmt.Errorf("could not build point: %w", err)
	}
	return rt.NewPoint(xf, yf, zf), nil
}

func vectorFromStrings(x, y, z string) (*rt.Tuple, error) {
	xf, yf, zf, err := parseXYZ(x, y, z)
	if err != nil {
		return nil, fmt.Errorf("could not build point: %w", err)
	}
	return rt.NewVector(xf, yf, zf), nil
}

func colorFromString(s string) (*rt.Color, error) {
	re, err := regexp.Compile(`\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)`)
	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}
	matches := re.FindStringSubmatch(s)
	if len(matches) >= 4 {
		parts := make([]float64, len(matches))
		for i := 1; i < len(matches); i++ {
			val, err := strconv.ParseFloat(matches[i], 64)
			if err != nil {
				return nil, err
			}
			parts[i] = val
		}
		log.Printf("Parsed color: R=%f, G=%f, B=%f", parts[1], parts[2], parts[3])
		return rt.NewColor(parts[1], parts[2], parts[3]), nil
	}

	return nil, fmt.Errorf("failed to parse %s to color", s)
}

// (\-?\d+\.?\d*)
func InitializeTuplesScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w+) ← tuple\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z, w float64) (context.Context, error) {
			t := rt.NewTuple(x, y, z, w)
			return context.WithValue(ctx, tupleKey{Name: name}, t), nil
		})
	sc.Given(
		`^(\w+) ← point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) (context.Context, error) {
			p := rt.NewPoint(x, y, z)
			return context.WithValue(ctx, tupleKey{Name: name}, p), nil
		})
	sc.Given(
		`^(\w+) ← vector\((\S+), (\S+), (\S+)\)$`,
		func(ctx context.Context, name string, x, y, z string) (context.Context, error) {
			v, err := vectorFromStrings(x, y, z)
			if err != nil {
				return ctx, err
			}
			return context.WithValue(ctx, tupleKey{Name: name}, v), nil
		})
	sc.Given(
		`^(\w+) ← color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) context.Context {
			return setColor(ctx, name, rt.NewColor(x, y, z))
		})
	sc.Given(
		`^(\w+) ← color\((\-?\d+\.?\d*), result ← lighting(m, light, position, eyev, normalv), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, r, g, b float64) (context.Context, error) {
			c := rt.NewColor(r, g, b)
			return context.WithValue(ctx, colorKey{Name: name}, c), nil
		})

	sc.When(
		`(\w) ← reflect\((\w), (\w)\)`,
		func(ctx context.Context, dest, tupleName, normalName string) (context.Context, error) {
			tuple, ok := ctx.Value(tupleKey{tupleName}).(*rt.Tuple)
			if !ok {
				return ctx, fmt.Errorf("no tuple named %s found", tupleName)
			}
			normal, ok := ctx.Value(tupleKey{normalName}).(*rt.Tuple)
			if !ok {
				return ctx, fmt.Errorf("no tuple named %s found", normalName)
			}
			r := tuple.Reflect(normal)
			return context.WithValue(ctx, tupleKey{dest}, r), nil
		})
	sc.When(
		`^(\w+) ← normalize\((\w+)\)$`,
		func(ctx context.Context, name1, name2 string) (context.Context, error) {
			tuple, ok := ctx.Value(tupleKey{Name: name2}).(*rt.Tuple)
			if !ok {
				return ctx, fmt.Errorf("no tuple named %s found", name2)
			}
			return context.WithValue(ctx, tupleKey{Name: name1}, tuple.Normalize()), nil
		})

	sc.Then(
		`^(\w).x = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, x float64) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			if tuple.X != x {
				return fmt.Errorf("expected tuple.X = %f, got %f", x, tuple.X)
			}
			return nil
		})
	sc.Then(
		`^(\w).y = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, y float64) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			if tuple.Y != y {
				return fmt.Errorf("expected tuple.X = %f, got %f", y, tuple.Y)
			}
			return nil
		})
	sc.Then(
		`^(\w).z = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, z float64) error {

			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			if tuple.Z != z {
				return fmt.Errorf("expected tuple.X = %f, got %f", z, tuple.Z)
			}
			return nil
		})
	sc.Then(
		`^(\w).w = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, w float64) error {

			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			if tuple.W != w {
				return fmt.Errorf("expected tuple.X = %f, got %f", w, tuple.W)
			}
			return nil
		})

	sc.Then(
		`^(\w).red = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, r float64) error {
			color, ok := ctx.Value(colorKey{Name: name}).(*rt.Color)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			if color.Red != r {
				return fmt.Errorf("expected tuple.X = %f, got %f", r, color.Red)
			}
			return nil
		})
	sc.Then(
		`^(\w).green = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, g float64) error {
			color, ok := ctx.Value(colorKey{Name: name}).(*rt.Color)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			if color.Green != g {
				return fmt.Errorf("expected tuple.X = %f, got %f", g, color.Green)
			}
			return nil
		})
	sc.Then(
		`^(\w).blue = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, b float64) error {
			color, ok := ctx.Value(colorKey{Name: name}).(*rt.Color)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			if color.Blue != b {
				return fmt.Errorf("expected tuple.X = %f, got %f", b, color.Blue)
			}
			return nil
		})

	sc.Then(
		`^(\w) = tuple\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z, w float64) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			expect := rt.NewTuple(x, y, z, w)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s, got %s", expect, tuple)
			}
			return nil
		})
	sc.Then(
		`^-(\w) = tuple\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z, w float64) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			tuple = tuple.Negate()
			expect := rt.NewTuple(x, y, z, w)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s, got %s", expect, tuple)
			}
			return nil
		})
	sc.Then(
		`^(\w) \* (\-?\d+\.?\d*) = tuple\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, m, x, y, z, w float64) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			tuple = tuple.MultiplyScalar(m)
			expect := rt.NewTuple(x, y, z, w)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s, got %s", expect, tuple)
			}
			return nil
		})
	sc.Then(
		`^(\w) / (\-?\d+\.?\d*) = tuple\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, d, x, y, z, w float64) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			tuple = tuple.DivideScalar(d)
			expect := rt.NewTuple(x, y, z, w)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s, got %s", expect, tuple)
			}
			return nil
		})

	sc.Then(
		`^(\w+) \+ (\w+) = tuple\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name1, name2 string, x, y, z, w float64) error {
			tuple1, ok := ctx.Value(tupleKey{Name: name1}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name1)
			}
			tuple2, ok := ctx.Value(tupleKey{Name: name2}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name2)
			}
			tuple := tuple1.Add(tuple2)
			expect := rt.NewTuple(x, y, z, w)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s, got %s", expect, tuple)
			}
			return nil
		})

	sc.Then(
		`^(\w) is a point$`,
		func(ctx context.Context, name string) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			if !tuple.IsPoint() {
				return fmt.Errorf("expected tuple %s to be a point", tuple)
			}
			return nil
		})
	sc.Then(
		`^(\w) is a vector$`,
		func(ctx context.Context, name string) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			if !tuple.IsVector() {
				return fmt.Errorf("expected tuple %s to be a vector", tuple)
			}
			return nil
		})
	sc.Then(
		`^(\w) is not a point$`,
		func(ctx context.Context, name string) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			if tuple.IsPoint() {
				return fmt.Errorf("expected tuple %s to not be a point", tuple)
			}
			return nil
		})
	sc.Then(
		`^(\w) is not a vector$`,
		func(ctx context.Context, name string) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			if tuple.IsVector() {
				return fmt.Errorf("expected tuple %s to not be a vector", tuple)
			}
			return nil
		})

	sc.Then(
		`^magnitude\((\w+)\) = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, magnitude float64) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			m := tuple.Magnitude()
			if m != magnitude {
				return fmt.Errorf("expected tuple %s to have magnitude %f, was %f", tuple, magnitude, m)
			}
			return nil
		})

	sc.Then(
		`^magnitude\((\w)\) = √(\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, magnitude float64) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			magnitude = math.Sqrt(magnitude)
			m := tuple.Magnitude()
			if m != magnitude {
				return fmt.Errorf("expected tuple %s to have magnitude %f, was %f", tuple, magnitude, m)
			}
			return nil
		})

	sc.Then(
		`^(\w+) - (\w+) = vector\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name1, name2 string, x, y, z float64) error {
			tuple1, ok := ctx.Value(tupleKey{Name: name1}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name1)
			}
			tuple2, ok := ctx.Value(tupleKey{Name: name2}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name2)
			}
			tuple := tuple1.Subtract(tuple2)
			expect := rt.NewVector(x, y, z)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s - %s to equal %s, was %s", tuple1, tuple2, expect, tuple)
			}
			return nil
		})
	sc.Then(
		`^(\w+) - (\w+) = point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name1, name2 string, x, y, z float64) error {
			tuple1, ok := ctx.Value(tupleKey{Name: name1}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name1)
			}
			tuple2, ok := ctx.Value(tupleKey{Name: name2}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name2)
			}
			tuple := tuple1.Subtract(tuple2)
			expect := rt.NewPoint(x, y, z)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s - %s to equal %s, was %s", tuple1, tuple2, expect, tuple)
			}
			return nil
		})
	sc.Then(
		`^(\w+) = point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			expect := rt.NewPoint(x, y, z)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s to equal %s", tuple, expect)
			}
			return nil
		})

	sc.Then(
		`^(\w+) \+ (\w+) = color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name1, name2 string, r, g, b float64) error {
			color1, ok := ctx.Value(colorKey{Name: name1}).(*rt.Color)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name1)
			}
			color2, ok := ctx.Value(colorKey{Name: name2}).(*rt.Color)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name1)
			}
			expect := rt.NewColor(r, g, b)
			color := color1.Add(color2)
			if !color.Equal(expect) {
				return fmt.Errorf("expected %s to equal %s", color, expect)
			}
			return nil
		})
	sc.Then(
		`^(\w+) - (\w+) = color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name1, name2 string, r, g, b float64) error {
			color1, ok := ctx.Value(colorKey{Name: name1}).(*rt.Color)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name1)
			}
			color2, ok := ctx.Value(colorKey{Name: name2}).(*rt.Color)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name1)
			}
			expect := rt.NewColor(r, g, b)
			color := color1.Subtract(color2)
			if !color.Equal(expect) {
				return fmt.Errorf("expected %s to equal %s", color, expect)
			}
			return nil
		})
	sc.Then(
		`^(\w+) \* (\d+) = color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name1 string, m, r, g, b float64) error {
			color, ok := ctx.Value(colorKey{Name: name1}).(*rt.Color)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name1)
			}
			expect := rt.NewColor(r, g, b)
			color = color.MultiplyScalar(m)
			if !color.Equal(expect) {
				return fmt.Errorf("expected %s to equal %s", color, expect)
			}
			return nil
		})
	sc.Then(
		`^(\w+) \* (\w\d+) = color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name1, name2 string, r, g, b float64) error {
			color1, ok := ctx.Value(colorKey{Name: name1}).(*rt.Color)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name1)
			}
			color2, ok := ctx.Value(colorKey{Name: name2}).(*rt.Color)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name2)
			}
			expect := rt.NewColor(r, g, b)
			color := color1.Multiply(color2)
			if !color.Equal(expect) {
				return fmt.Errorf("expected %s to equal %s", color, expect)
			}
			return nil
		})
	sc.Then(
		`^(\w+) = color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, r, g, b float64) error {
			color, err := getColor(ctx, name)
			if err != nil {
				return err
			}

			expect := rt.NewColor(r, g, b)
			if !color.Equal(expect) {
				return fmt.Errorf("expected %s to equal %s", color, expect)
			}
			return nil
		})

	sc.Then(
		`^normalize\((\w+)\) = vector\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			tuple = tuple.Normalize()
			expect := rt.NewVector(x, y, z)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s to equal %s", tuple, expect)
			}
			return nil
		})

	sc.Then(
		`^normalize\((\w+)\) = approximately vector\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name string, x, y, z float64) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			tuple = tuple.Normalize()
			expect := rt.NewVector(x, y, z)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s to equal %s", tuple, expect)
			}
			return nil
		})

	sc.Then(
		`^cross\((\w+), (\w+)\) = vector\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, name1, name2 string, x, y, z float64) error {
			tuple1, ok := ctx.Value(tupleKey{Name: name1}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name1)
			}
			tuple2, ok := ctx.Value(tupleKey{Name: name2}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name2)
			}
			tuple := tuple1.Cross(tuple2)
			expect := rt.NewVector(x, y, z)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s to equal %s", tuple, expect)
			}
			return nil
		})
	sc.Then(
		`^dot\((\w+), (\w+)\) = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name1, name2 string, expect float64) error {
			tuple1, ok := ctx.Value(tupleKey{Name: name1}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name1)
			}
			tuple2, ok := ctx.Value(tupleKey{Name: name2}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name2)
			}
			dot := tuple1.Dot(tuple2)
			if dot != expect {
				return fmt.Errorf("expected %f to equal %f", dot, expect)
			}
			return nil
		})

	sc.Then(
		`^(\w+) = vector\((\S+), (\S+), (\S+)\)$`,
		func(ctx context.Context, name string, x, y, z string) error {
			tuple, ok := ctx.Value(tupleKey{Name: name}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", name)
			}
			expect, err := vectorFromStrings(x, y, z)
			if err != nil {
				return err
			}
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s to equal %s", tuple, expect)
			}
			return nil
		})
	sc.Then(
		`^(\w+) = normalize\((\w+)\)$`,
		func(ctx context.Context, expectName, sourceName string) error {
			expect, ok := ctx.Value(tupleKey{Name: expectName}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", expectName)
			}
			source, ok := ctx.Value(tupleKey{Name: sourceName}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", sourceName)
			}
			normalized := source.Normalize()
			if !expect.Equal(normalized) {
				return fmt.Errorf("expected %s, got %s", expect, normalized)
			}
			return nil
		})
}
