package test

import (
	"context"
	"fmt"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

type lightKey struct{ Name string }

func setLighting(ctx context.Context, dest, materialName, lightName, positionName, eyeName, normalName string, inShadow bool) (context.Context, error) {
	material, err := getMaterial(ctx, materialName)
	if err != nil {
		return ctx, err
	}
	light, err := getLight(ctx, lightName)
	if err != nil {
		return ctx, err
	}
	position, err := getTuple(ctx, positionName)
	if err != nil {
		return ctx, err
	}
	eye, err := getTuple(ctx, eyeName)
	if err != nil {
		return ctx, err
	}
	normal, err := getTuple(ctx, normalName)
	if err != nil {
		return ctx, err
	}

	color, err := rt.Lighting(material, nil, light, position, eye, normal, inShadow)
	if err != nil {
		return ctx, err
	}
	return context.WithValue(ctx, colorKey{dest}, color), nil
}

func InitializeLightsScenario(sc *godog.ScenarioContext) {

	sc.Given(
		`^(\w+) ← point_light\(point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\), color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		func(ctx context.Context, dest string, x, y, z, r, g, b float64) (context.Context, error) {
			tuple := rt.NewPoint(x, y, z)
			color := rt.NewColor(r, g, b)
			return context.WithValue(ctx, lightKey{dest}, rt.NewPointLight(tuple, color)), nil
		})
	sc.Given(
		`^(\w+) ← (true|false)$`,
		func(ctx context.Context, dest, val string) context.Context {
			return setBoolean(ctx, dest, val == "true")
		})

	sc.When(
		`^(\w+) ← point_light\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, position, intensity string) (context.Context, error) {
			tuple, ok := ctx.Value(tupleKey{position}).(*rt.Tuple)
			if !ok {
				return ctx, fmt.Errorf("no tuple named %s found", position)
			}
			color, ok := ctx.Value(colorKey{intensity}).(*rt.Color)
			if !ok {
				return ctx, fmt.Errorf("no color named %s found", intensity)
			}
			return context.WithValue(ctx, lightKey{dest}, rt.NewPointLight(tuple, color)), nil
		})
	sc.When(`^(\w+) ← lighting\((\w+), (\w+), (\w+), (\w+), (\w+)\)$`,
		func(ctx context.Context, dest, materialName, lightName, positionName, eyeName, normalName string) (context.Context, error) {
			return setLighting(ctx, dest, materialName, lightName, positionName, eyeName, normalName, false)
		})
	sc.When(`^(\w+) ← lighting\((\w+), (\w+), (\w+), (\w+), (\w+), (\w+)\)$`,
		func(ctx context.Context, dest, materialName, lightName, positionName, eyeName, normalName, shadowFlagName string) (context.Context, error) {
			shadowFlag, err := getBoolean(ctx, shadowFlagName)
			if err != nil {
				return ctx, err
			}
			return setLighting(ctx, dest, materialName, lightName, positionName, eyeName, normalName, shadowFlag)
		})
	sc.When(
		`^(\w+) ← lighting\((\w+), (\w+), point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\), (\w+), (\w+), false\)$`,
		func(ctx context.Context, dest, materialName, lightName string, x, y, z float64, eyeName, normalName string) (context.Context, error) {
			const posName = "testPlaceholderPoint"
			ctx = setTuple(ctx, posName, rt.NewPoint(x, y, z))
			return setLighting(ctx, dest, materialName, lightName, posName, eyeName, normalName, false)
		})

	sc.Then(
		`^(\w+).position = (\w+)$`,
		func(ctx context.Context, dest, position string) error {
			light, ok := ctx.Value(lightKey{dest}).(*rt.Light)
			if !ok {
				return fmt.Errorf("no light named %s found", dest)
			}
			tuple, ok := ctx.Value(tupleKey{position}).(*rt.Tuple)
			if !ok {
				return fmt.Errorf("no tuple named %s found", position)
			}
			if !light.Position.Equal(tuple) {
				return fmt.Errorf("expected %s, got %s", tuple, light.Position)
			}
			return nil
		})
	sc.Then(
		`^(\w+).intensity = (\w+)$`,
		func(ctx context.Context, dest, intensity string) error {
			light, ok := ctx.Value(lightKey{dest}).(*rt.Light)
			if !ok {
				return fmt.Errorf("no light named %s found", dest)
			}
			color, ok := ctx.Value(colorKey{intensity}).(*rt.Color)
			if !ok {
				return fmt.Errorf("no tuple named %s found", intensity)
			}
			if !light.Intensity.Equal(color) {
				return fmt.Errorf("expected %s, got %s", color, light.Position)
			}
			return nil
		})

}
