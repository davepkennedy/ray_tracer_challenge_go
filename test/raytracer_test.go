package test

import (
	"testing"

	"github.com/cucumber/godog"
)

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		Name:                "raytracer",
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,

			//ShowStepDefinitions: true,
			Strict: true,
			//StopOnFailure: true,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	InitializeSphereScenario(ctx)
	InitializeRaysScenario(ctx)
	InitializeTuplesScenario(ctx)
	InitializeCanvasScenario(ctx)
	InitializeMatrixScenario(ctx)
	InitializeIntersectionScenario(ctx)
	InitializeTransformationScenario(ctx)
	InitializeLightsScenario(ctx)
	InitializeMaterialsScenario(ctx)
	InitializeWorldScenario(ctx)
	InitializeCameraScenario(ctx)
	InitializeShapesScenario(ctx)
	InitializePlaneScenario(ctx)
	InitializePatternsScenario(ctx)
	InitializeCubeScenario(ctx)
	InitializeCylinderScenario(ctx)
	InitializeConesScenario(ctx)
	InitializeGroupsScenario(ctx)
	InitializeTrianglesScenario(ctx)
	InitializeObjFileScenario(ctx)
	InitializeSmoothTrianglesScenario(ctx)
	InitializeCSGScenario(ctx)
}
