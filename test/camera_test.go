package test

import (
	"context"
	"fmt"
	"math"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func InitializeCameraScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w+) ← (\d+)$`,
		func (ctx context.Context, name string, val int) context.Context {
			return setInt(ctx, name, val)
		})
    sc.Given(
		`^(\w+) ← π/2$`,
		func (ctx context.Context, name string) context.Context {
			return setFloat(ctx, name, math.Pi/2.0)
		})
	sc.Given(
		`^(\w+) ← camera\((\d+), (\d+), π/2\)$`,
		func (ctx context.Context, dest string, h, v int) context.Context {
			c := rt.NewCamera(h,v, math.Pi/2.0)
			return setCamera(ctx, dest, c)
		})

    sc.When(
		`^(\w+) ← camera\((\w+), (\w+), (\w+)\)$`,
		func (ctx context.Context, dest, h, v, f string) (context.Context, error) {
			hsize, err := getInt(ctx, h)
			if err != nil {return ctx, err}
			vsize, err := getInt(ctx, v)
			if err != nil {return ctx, err}
			fov, err := getFloat(ctx, f)
			if err != nil {return ctx, err}

			 camera := rt.NewCamera (hsize, vsize, fov)
			 return setCamera(ctx, dest, camera), nil
		})
	sc.When(
		`^(\w+) ← ray_for_pixel\((\w+), (\d+), (\d+)\)$`,
		func (ctx context.Context, dest, source string, x, y int) (context.Context, error) {
			camera, err := getCamera(ctx, source)
			if err != nil {return ctx, err}

			ray, err := camera.RayForPixel(x, y)
			if err != nil {return ctx, err}

			return setRay (ctx, dest, ray), nil
		})

	sc.Then(
		`^(\w+).hsize = (\d+)$`,
		func (ctx context.Context, dest string, val int) error {
			camera, err := getCamera(ctx, dest)
			if err != nil {return err}

			if camera.HSize != val {
				return fmt.Errorf("expected %d, got %d", val, camera.HSize)
			}
			return nil
		})
    sc.Then(
		`^(\w+).vsize = (\d+)$`,
		func (ctx context.Context, dest string, val int) error {
			camera, err := getCamera(ctx, dest)
			if err != nil {return err}

			if camera.VSize != val {
				return fmt.Errorf("expected %d, got %d", val, camera.VSize)
			}
			return nil
		})

    sc.Then(
		`^(\w+).field_of_view = π/2$`, 
		func (ctx context.Context, dest string) error {
			camera, err := getCamera(ctx, dest)
			if err != nil {return err}

			expect := math.Pi/2.0
			if expect != camera.FieldOfView {
				return fmt.Errorf("expected %f, got %f", expect, camera.FieldOfView)
			}
			return nil
		})

    sc.Then(
		`^c.transform = (\w+)$`,
		func (ctx context.Context, source string) error {
			camera, err := getCamera(ctx, "c")
			if err != nil {return err}
			m, err := getMatrix(ctx, source)
			if err != nil {return err}

			if !camera.Transform.Equal(m) {
				return fmt.Errorf("expected %s, got %s", m, camera.Transform)
			}
			return nil
		})

	sc.Then(
		`^(\w+).pixel_size = (\-?\d+\.?\d*)$`,
		func (ctx context.Context, dest string, val float64) error {
			camera, err := getCamera(ctx, dest)
			if err != nil {return err}

			if camera.PixelSize != val {
				return fmt.Errorf("expected %f, got %f", val, camera.PixelSize)
			}
			return nil
		})
}