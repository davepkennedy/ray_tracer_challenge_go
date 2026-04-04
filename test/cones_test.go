package test

import (
	"context"
	"raytracer/pkg/rt"
	"github.com/cucumber/godog"
)

func InitializeConesScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w+) ← cone\(\)$`,
		func (ctx context.Context, name string) context.Context {
			cone := rt.NewCone()
			return setShape(ctx, name, cone)
		})
}