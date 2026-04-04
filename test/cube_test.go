package test

import (
	"context"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func InitializeCubeScenario(sc *godog.ScenarioContext) {
	sc.Given (
		`(\w) ← cube\(\)`,
		func(ctx context.Context, name string) context.Context {
			cube := rt.NewCube()
			return setShape(ctx, name, cube)
		})

	sc.When(
		`(\w+) ← local_normal_at\((\w+), (\w+)\)`,
		func(ctx context.Context, dest, source, pointName string) (context.Context, error) {
			shape, err := getShape(ctx, source)
			if err != nil {return ctx, err}

			point, err := getTuple(ctx, pointName)
			if err != nil {return ctx, err}

			normal, err := shape.Trait.LocalNormalAt(point)
			if err != nil {return ctx, err}
			return setTuple(ctx, dest, normal), nil
		})
}