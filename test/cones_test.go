package test

import (
	"context"
	"github.com/cucumber/godog"
	"raytracer/pkg/rt"
)

func InitializeConesScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w+) ← cone\(\)$`,
		func(ctx context.Context, name string) context.Context {
			cone := rt.NewCone()
			return setShape(ctx, name, cone)
		})
}
