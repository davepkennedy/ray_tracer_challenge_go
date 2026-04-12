package test

import (
	"context"
	"fmt"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func InitializeMaterialsScenario(sc *godog.ScenarioContext) {

	sc.Given(
		`^(\w) ← material\(\)$`,
		func(ctx context.Context, dest string) context.Context {
			return setMaterial(ctx, dest, rt.NewMaterial())
		})

	sc.Given(
		`^(\w+).pattern ← stripe_pattern\(color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\), color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		func(ctx context.Context, target string, r1, g1, b1, r2, g2, b2 float64) (context.Context, error) {
			material, err := getMaterial(ctx, target)
			if err != nil {
				return ctx, err
			}

			material.Pattern = rt.NewStripePattern(rt.NewColor(r1, g1, b1), rt.NewColor(r2, g2, b2))

			return ctx, nil
		})

	sc.Given(
		`^(\w).ambient ← (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, val float64) (context.Context, error) {
			material, err := getMaterial(ctx, dest)
			if err != nil {
				return ctx, err
			}

			material.Ambient = val
			return ctx, nil
		})
	sc.Given(
		`^(\w).diffuse ← (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, val float64) (context.Context, error) {
			material, err := getMaterial(ctx, dest)
			if err != nil {
				return ctx, err
			}

			material.Diffuse = val
			return ctx, nil
		})
	sc.Given(
		`^(\w).specular ← (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, val float64) (context.Context, error) {
			material, err := getMaterial(ctx, dest)
			if err != nil {
				return ctx, err
			}
			material.Specular = val
			return ctx, nil
		})

	sc.Then(
		`^(\w+).color = color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, dest string, r, g, b float64) error {
			material, err := getMaterial(ctx, dest)
			if err != nil {
				return err
			}
			color := rt.NewColor(r, g, b)
			if !material.Color.Equal(color) {
				return fmt.Errorf("expected %s, got %s", color, material.Color)
			}
			return nil
		})
	sc.Then(
		`(\w+).ambient = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, f float64) error {
			material, err := getMaterial(ctx, dest)
			if err != nil {
				return err
			}

			if material.Ambient != f {
				return fmt.Errorf("expected %f, got %f", f, material.Ambient)
			}
			return nil
		})
	sc.Then(
		`(\w+).diffuse = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, f float64) error {
			material, err := getMaterial(ctx, dest)
			if err != nil {
				return err
			}

			if material.Diffuse != f {
				return fmt.Errorf("expected %f, got %f", f, material.Diffuse)
			}
			return nil
		})
	sc.Then(
		`(\w+).specular = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, f float64) error {
			material, err := getMaterial(ctx, dest)
			if err != nil {
				return err
			}

			if material.Specular != f {
				return fmt.Errorf("expected %f, got %f", f, material.Specular)
			}
			return nil
		})
	sc.Then(
		`^(\w+).shininess = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, f float64) error {
			material, err := getMaterial(ctx, dest)
			if err != nil {
				return err
			}

			if material.Shininess != f {
				return fmt.Errorf("expected %f, got %f", f, material.Shininess)
			}
			return nil
		})
	sc.Then(
		`^(\w).reflective = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, f float64) error {
			material, err := getMaterial(ctx, dest)
			if err != nil {
				return err
			}

			if material.Reflective != f {
				return fmt.Errorf("expected %f, got %f", f, material.Reflective)
			}
			return nil
		})
	sc.Then(
		`^(\w).transparency = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, f float64) error {
			material, err := getMaterial(ctx, dest)
			if err != nil {
				return err
			}

			if material.Transparency != f {
				return fmt.Errorf("expected %f, got %f", f, material.Transparency)
			}
			return nil
		})
	sc.Then(
		`^(\w).refractive_index = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, dest string, f float64) error {
			material, err := getMaterial(ctx, dest)
			if err != nil {
				return err
			}

			if material.RefractiveIndex != f {
				return fmt.Errorf("expected %f, got %f", f, material.RefractiveIndex)
			}
			return nil
		})

	sc.Then(
		`(\w+) = material\(\)`,
		func(ctx context.Context, dest string) error {
			material, err := getMaterial(ctx, dest)
			if err != nil {
				return err
			}

			expect := rt.NewMaterial()
			if !material.Equal(expect) {
				return fmt.Errorf("expected %s, got %s", expect, material)
			}
			return nil
		})
}
