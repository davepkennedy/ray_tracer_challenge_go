package test

import (
	"context"
	"errors"
	"fmt"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

var (
	ErrTraitNotCSG = errors.New("trait not CSG")
)

func InitializeCubeScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w+) ← cube\(\)$`,
		func(ctx context.Context, name string) context.Context {
			cube := rt.NewCube()
			return setShape(ctx, name, cube)
		})

	sc.When(
		`^(\w+) ← local_normal_at\((\w+), (\w+)\)$`,
		func(ctx context.Context, dest, source, pointName string) (context.Context, error) {
			shape, err := getShape(ctx, source)
			if err != nil {
				return ctx, err
			}

			point, err := getTuple(ctx, pointName)
			if err != nil {
				return ctx, err
			}

			normal, err := shape.Trait.LocalNormalAt(point, nil)
			if err != nil {
				return ctx, err
			}
			return setTuple(ctx, dest, normal), nil
		})

	sc.Then(
		`^(\w+).operation = "(\w+)"$`,
		func(ctx context.Context, dest, operation string) error {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return err
			}

			trait, ok := shape.Trait.(*rt.CSG)
			if !ok {
				return ErrTraitNotCSG
			}

			if trait.Operation != operation {
				return fmt.Errorf("expected %s, got %s", operation, trait.Operation)
			}
			return nil
		})
	sc.Then(
		`^(\w+).left = (\w+)$`,
		func(ctx context.Context, destName, shapeName string) error {
			dest, err := getShape(ctx, destName)
			if err != nil {
				return err
			}

			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return err
			}

			trait, ok := dest.Trait.(*rt.CSG)
			if !ok {
				return ErrTraitNotCSG
			}

			if !trait.Left.Equal(shape) {
				return fmt.Errorf("expected %s, got %s", shape, trait.Left)
			}
			return nil
		})
	sc.Then(
		`^(\w+).right = (\w+)$`,
		func(ctx context.Context, destName, shapeName string) error {
			dest, err := getShape(ctx, destName)
			if err != nil {
				return err
			}

			shape, err := getShape(ctx, shapeName)
			if err != nil {
				return err
			}

			trait, ok := dest.Trait.(*rt.CSG)
			if !ok {
				return ErrTraitNotCSG
			}

			if !trait.Right.Equal(shape) {
				return fmt.Errorf("expected %s, got %s", shape, trait.Left)
			}
			return nil
		})
	sc.Then(
		`^(\w+).parent = (\w+)$`,
		func(ctx context.Context, childName, parentName string) error {
			child, err := getShape(ctx, childName)
			if err != nil {
				return err
			}

			parent, err := getShape(ctx, parentName)
			if err != nil {
				return err
			}

			if !parent.Equal(child.Parent) {
				return fmt.Errorf("expected %s, got %s", parent, child.Parent)
			}
			return nil
		})
}
