package test

import (
	"context"
	"fmt"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

type TestPattern struct {
}

func (t TestPattern) ColorAt(pt *rt.Tuple, a, b *rt.Color) *rt.Color {
	return rt.NewColor(pt.X, pt.Y, pt.Z)
}

func newTestPattern() *rt.Pattern {
	return rt.NewPattern(rt.NewColor(0, 0, 0), rt.NewColor(1, 1, 1), TestPattern{})
}

func getTransform(s string) func(x, y, z float64) *rt.Matrix {
	switch s {
	case "scaling":
		return rt.Scaling
	case "translation":
		return rt.Translation
	}

	return func(x, y, z float64) *rt.Matrix { return rt.IdentityMatrix() }
}

func setPatternTransform(ctx context.Context, target, transform string, x, y, z float64) (context.Context, error) {
	pattern, err := getPattern(ctx, target)
	if err != nil {
		return ctx, err
	}

	t := getTransform(transform)(x, y, z)
	pattern.Transform = t

	return ctx, nil
}

func patternAtShape(ctx context.Context, dest, patternName, objectName string, x, y, z float64) (context.Context, error) {
	pattern, err := getPattern(ctx, patternName)
	if err != nil {
		return ctx, err
	}

	shape, err := getShape(ctx, objectName)
	if err != nil {
		return ctx, err
	}

	pt := rt.NewPoint(x, y, z)
	color, err := pattern.ColorAtObject(shape, pt)
	if err != nil {
		return ctx, err
	}

	return setColor(ctx, dest, color), nil
}

func patternAtPoint(ctx context.Context, dest string, x, y, z float64, colorName string) error {
	pattern, err := getPattern(ctx, dest)
	if err != nil {
		return err
	}
	expect, err := getColor(ctx, colorName)
	if err != nil {
		return err
	}

	p := rt.NewPoint(x, y, z)
	color := pattern.ColorAt(p)

	if !expect.Equal(color) {
		return fmt.Errorf("expected %s, got %s", expect, color)
	}
	return nil
}

func InitializePatternsScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w+) ← stripe_pattern\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, firstColor, secondColor string) (context.Context, error) {
			c1, err := getColor(ctx, firstColor)
			if err != nil {
				return ctx, err
			}
			c2, err := getColor(ctx, secondColor)
			if err != nil {
				return ctx, err
			}

			return setPattern(ctx, dest, rt.NewStripePattern(c1, c2)), nil
		})
	sc.Given(
		`^(\w+) ← gradient_pattern\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, firstColor, secondColor string) (context.Context, error) {
			c1, err := getColor(ctx, firstColor)
			if err != nil {
				return ctx, err
			}
			c2, err := getColor(ctx, secondColor)
			if err != nil {
				return ctx, err
			}

			return setPattern(ctx, dest, rt.NewGradientPattern(c1, c2)), nil
		})
	sc.Given(
		`^(\w+) ← ring_pattern\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, firstColor, secondColor string) (context.Context, error) {
			c1, err := getColor(ctx, firstColor)
			if err != nil {
				return ctx, err
			}
			c2, err := getColor(ctx, secondColor)
			if err != nil {
				return ctx, err
			}

			return setPattern(ctx, dest, rt.NewRingPattern(c1, c2)), nil
		})
	sc.Given(
		`^(\w+) ← checkers_pattern\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, firstColor, secondColor string) (context.Context, error) {
			c1, err := getColor(ctx, firstColor)
			if err != nil {
				return ctx, err
			}
			c2, err := getColor(ctx, secondColor)
			if err != nil {
				return ctx, err
			}

			return setPattern(ctx, dest, rt.NewCheckersPattern(c1, c2)), nil
		})
	sc.Given(
		`^(\w+) ← test_pattern\(\)$`,
		func(ctx context.Context, dest string) context.Context {
			return setPattern(ctx, dest, newTestPattern())
		})

	sc.Given(
		`^set_pattern_transform\((\w+), (\w+)\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		setPatternTransform)

	sc.When(
		`^set_pattern_transform\((\w+), (\w+)\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		setPatternTransform)

	sc.When(
		`^(\w+) ← stripe_at_object\((\w+), (\w+), point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		patternAtShape)
	sc.When(
		`^(\w+) ← pattern_at_shape\((\w+), (\w+), point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		patternAtShape)

	sc.Then(
		`(\w+).transform = identity_matrix$`,
		func(ctx context.Context, dest string) error {
			pattern, err := getPattern(ctx, dest)
			if err != nil {
				return err
			}

			expect := rt.IdentityMatrix()
			if !pattern.Transform.Equal(expect) {
				return fmt.Errorf("expected %s, got %s", expect, pattern.Transform)
			}
			return nil
		})
	sc.Then(
		`(\w+).a = (\w+)$`,
		func(ctx context.Context, dest, colorName string) error {
			pattern, err := getPattern(ctx, dest)
			if err != nil {
				return err
			}

			color, err := getColor(ctx, colorName)
			if err != nil {
				return err
			}

			if !pattern.A.Equal(color) {
				return fmt.Errorf("expected %s, got %s", color, pattern.A)
			}
			return nil
		})

	sc.Then(
		`(\w+).b = (\w+)$`,
		func(ctx context.Context, dest, colorName string) error {
			pattern, err := getPattern(ctx, dest)
			if err != nil {
				return err
			}

			color, err := getColor(ctx, colorName)
			if err != nil {
				return err
			}

			if !pattern.B.Equal(color) {
				return fmt.Errorf("expected %s, got %s", color, pattern.B)
			}
			return nil
		})

	sc.Then(
		`^stripe_at\((\w+), point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\) = (\w+)$`,
		patternAtPoint)

	sc.Then(
		`^pattern_at\((\w+), point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\) = (\w+)$`,
		patternAtPoint)
	sc.Then(
		`^pattern_at\((\w+), point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\) = color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, dest string, x, y, z float64, r, g, b float64) error {
			pattern, err := getPattern(ctx, dest)
			if err != nil {
				return err
			}

			p := rt.NewPoint(x, y, z)
			color := pattern.ColorAt(p)

			expect := rt.NewColor(r, g, b)

			if !expect.Equal(color) {
				return fmt.Errorf("expected %s, got %s", expect, color)
			}
			return nil
		})

	sc.Then(
		`^c = (\w+)$`,
		func(ctx context.Context, colorName string) error {
			color, err := getColor(ctx, "c")
			if err != nil {
				return err
			}

			expect, err := getColor(ctx, colorName)
			if err != nil {
				return err
			}

			if !expect.Equal(color) {
				return fmt.Errorf("expected %s, got %s", expect, color)
			}
			return nil
		})
	sc.Then(
		`^(\w\w+).transform = translation\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, dest string, x, y, z float64) error {
			pattern, err := getPattern(ctx, dest)
			if err != nil {
				return err
			}

			m := rt.Translation(x, y, z)

			if !pattern.Transform.Equal(m) {
				return fmt.Errorf("expected %s, got %s", m, pattern.Transform)
			}
			return nil
		})
}
