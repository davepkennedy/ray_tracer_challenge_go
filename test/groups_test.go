package test

import (
	"fmt"
	//"math"
	"context"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func addChild(ctx context.Context, target, source string) (context.Context, error) {
	targetShape, err := getShape(ctx, target)
	if err != nil {
		return ctx, err
	}
	sourceShape, err := getShape(ctx, source)
	if err != nil {
		return ctx, err
	}

	targetShape.AddChildren(sourceShape)
	return ctx, nil
}

func InitializeGroupsScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w+) ← group\(\)$`,
		func(ctx context.Context, name string) context.Context {
			group := rt.NewGroup()
			return setShape(ctx, name, group)
		})

	sc.Given(`^add_child\((\w+), (\w+)\)$`, addChild)
	sc.When(`^add_child\((\w+), (\w+)\)$`, addChild)

	sc.When(
		`^(\w) ← obj_to_group\((\w+)\)$`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			objFile, err := getObjFile(ctx, source)
			if err != nil {
				return ctx, err
			}

			group := objFile.ToGroup()
			return setShape(ctx, dest, group), nil
		})

	sc.Then(
		`^(\w) is empty$`,
		func(ctx context.Context, name string) (context.Context, error) {
			shape, err := getShape(ctx, name)
			if err != nil {
				return ctx, err
			}

			groupTrait, ok := shape.Trait.(*rt.Group)
			if !ok {
				return ctx, fmt.Errorf("expected %s to be a group", name)
			}

			if len(groupTrait.Children) != 0 {
				return ctx, fmt.Errorf("expected group to be empty, but had %d children", len(groupTrait.Children))
			}
			return ctx, nil
		})

	sc.Then(
		`^(\w) is not empty$`,
		func(ctx context.Context, name string) (context.Context, error) {
			shape, err := getShape(ctx, name)
			if err != nil {
				return ctx, err
			}

			groupTrait, ok := shape.Trait.(*rt.Group)
			if !ok {
				return ctx, fmt.Errorf("expected %s to be a group", name)
			}

			if len(groupTrait.Children) == 0 {
				return ctx, fmt.Errorf("expected group to be not empty, but had %d children", len(groupTrait.Children))
			}
			return ctx, nil
		})

	sc.Then(
		`^g.transform = identity_matrix$`,
		func(ctx context.Context) error {
			shape, err := getShape(ctx, "g")
			if err != nil {
				return err
			}

			expect := rt.IdentityMatrix()
			if !shape.Transform.Equal(expect) {
				return fmt.Errorf("expected %s, got %s", expect, shape.Transform)
			}
			return nil
		})

	sc.Then(
		`^(\w) includes (\w)$`,
		func(ctx context.Context, groupName, childName string) error {
			groupShape, err := getShape(ctx, groupName)
			if err != nil {
				return err
			}
			childShape, err := getShape(ctx, childName)
			if err != nil {
				return err
			}

			groupTrait, ok := groupShape.Trait.(*rt.Group)
			if !ok {
				return fmt.Errorf("expected %s to be a group", groupName)
			}

			for _, child := range groupTrait.Children {
				if child == childShape {
					return nil
				}
			}
			return fmt.Errorf("expected %s to include %s", groupName, childName)
		})

	/*
		sc.Then (
			`^(\w).parent = (\w)$`,
			func(ctx context.Context, childName, parentName string) error {
				childShape, err := getShape(ctx, childName)
				if err != nil {return err}

				parentShape, err := getShape(ctx, parentName)
				if err != nil {return err}

				if childShape.Parent != parentShape {
					return fmt.Errorf("expected %s to have %s as parent", childName, parentName)
				}
				return nil
			})
	*/
}
