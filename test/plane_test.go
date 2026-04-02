package test

import (
	"context"
	
	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func InitializePlaneScenario(sc *godog.ScenarioContext) {
	sc.Given (
		`^(\w+) ← plane\(\)$`,
		func (ctx context.Context, dest string) context.Context {
			return setShape(ctx, dest, rt.NewPlane())
		})

	sc.When(
		`^(\w+) ← local_normal_at\((\w+), point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)\)$`,
		func (ctx context.Context, dest, source string, x, y, z float64) (context.Context, error) {
			shape, err := getShape(ctx, source)
			if err != nil {return ctx, err}

			p := rt.NewPoint(x,y,z)

			m, err := shape.Trait.LocalNormalAt(p)
			if err != nil {return ctx, err}

			return setTuple(ctx, dest, m), nil
		})
}