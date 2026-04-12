package test

import (
	"context"
	"fmt"
	"os"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func parseFile(ctx context.Context, dest string, src string) (context.Context, error) {
	content, err := getString(ctx, src)
	if err != nil {
		return ctx, err
	}
	objFile, err := rt.ParseObjFileFromString(content)
	if err != nil {
		return ctx, err
	}
	return setObjFile(ctx, dest, objFile), nil
}

func InitializeObjFileScenario(sc *godog.ScenarioContext) {
	sc.Given(
		`^(\w+) ← a file containing:`,
		func(ctx context.Context, dest string, fileContent *godog.DocString) (context.Context, error) {
			return setString(ctx, dest, fileContent.Content), nil
		})
	sc.Given(
		`^(\w+) ← the file "(\S+)"$`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			content, err := os.ReadFile("files\\" + source)
			if err != nil {
				return ctx, err
			}

			return setString(ctx, dest, string(content)), nil
		})

	sc.Given(`^(\w+) ← parse_obj_file\((\w+)\)`, parseFile)
	sc.When(`^(\w+) ← parse_obj_file\((\w+)\)`, parseFile)

	sc.When(
		`^(\w) ← (\w+).default_group$`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			objFile, err := getObjFile(ctx, source)
			if err != nil {
				return ctx, err
			}

			return setShape(ctx, dest, objFile.DefaultGroup), nil
		})
	sc.When(
		`^(\w+) ← first child of (\w+)$`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			group, err := getShape(ctx, source)
			if err != nil {
				return ctx, err
			}

			if len(group.Trait.(*rt.Group).Children) < 1 {
				return ctx, fmt.Errorf("needed at least first child of group")
			}
			return setShape(ctx, dest, group.Trait.(*rt.Group).Children[0]), nil
		})
	sc.When(
		`^(\w+) ← second child of (\w+)$`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			group, err := getShape(ctx, source)
			if err != nil {
				return ctx, err
			}

			if len(group.Trait.(*rt.Group).Children) < 2 {
				return ctx, fmt.Errorf("needed at least second child of group")
			}
			return setShape(ctx, dest, group.Trait.(*rt.Group).Children[1]), nil
		})
	sc.When(
		`^(\w+) ← third child of (\w+)$`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			group, err := getShape(ctx, source)
			if err != nil {
				return ctx, err
			}

			if len(group.Trait.(*rt.Group).Children) < 3 {
				return ctx, fmt.Errorf("needed at least third child of group")
			}
			return setShape(ctx, dest, group.Trait.(*rt.Group).Children[2]), nil
		})
	sc.When(
		`^(\w+) ← "(\w+)" from (\w+)$`,
		func(ctx context.Context, dest, name, source string) (context.Context, error) {
			objFile, err := getObjFile(ctx, source)
			if err != nil {
				return ctx, err
			}

			group, ok := objFile.Groups[name]
			if !ok {
				return ctx, fmt.Errorf("no group %s found in %s", name, source)
			}

			return setShape(ctx, dest, group), nil
		})

	sc.Then(
		`^(\w+) should have ignored (\d+) lines$`,
		func(ctx context.Context, objFileName string, expectedIgnored int) error {
			objFile, err := getObjFile(ctx, objFileName)
			if err != nil {
				return err
			}
			if objFile.Ignored != expectedIgnored {
				return fmt.Errorf("expected %d lines to be ignored, but got %d", expectedIgnored, objFile.Ignored)
			}
			return nil
		})

	sc.Then(
		`^(\w+).vertices\[(\d+)\] = point\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, objFileName string, index int, x float64, y float64, z float64) error {
			objFile, err := getObjFile(ctx, objFileName)
			index -= 1
			if err != nil {
				return err
			}
			if index < 0 || index >= len(objFile.Vertices) {
				return fmt.Errorf("index %d is out of bounds", index)
			}
			vertex := objFile.Vertices[index]
			if vertex.X != x || vertex.Y != y || vertex.Z != z {
				return fmt.Errorf("expected vertex (%f, %f, %f), but got (%f, %f, %f)", x, y, z, vertex.X, vertex.Y, vertex.Z)
			}
			return nil
		})

	sc.Then(
		`^(\w+).normals\[(\d+)\] = vector\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, objFileName string, index int, x float64, y float64, z float64) error {
			objFile, err := getObjFile(ctx, objFileName)
			index -= 1
			if err != nil {
				return err
			}
			if index < 0 || index >= len(objFile.Normals) {
				return fmt.Errorf("index %d is out of bounds", index)
			}
			vertex := objFile.Normals[index]
			if vertex.X != x || vertex.Y != y || vertex.Z != z {
				return fmt.Errorf("expected vertex (%f, %f, %f), but got (%f, %f, %f)", x, y, z, vertex.X, vertex.Y, vertex.Z)
			}
			return nil
		})

	sc.Then(
		`^(\w+).p1 = (\w+).vertices\[(\d+)\]$`,
		func(ctx context.Context, dest, source string, idx int) error {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return err
			}

			objFile, err := getObjFile(ctx, source)
			if err != nil {
				return err
			}

			idx -= 1

			if idx < 0 && len(objFile.Vertices) <= idx {
				return fmt.Errorf("idx %d is out of bounds %d", idx, len(objFile.Vertices))
			}

			trait, err := getTriangle(shape)
			if err != nil {
				return err
			}
			if !trait.P1.Equal(objFile.Vertices[idx]) {
				return fmt.Errorf("expected p1 to be %s, got %s", objFile.Vertices[idx], trait.P1)
			}
			return nil
		})
	sc.Then(
		`^(\w+).p2 = (\w+).vertices\[(\d+)\]$`,
		func(ctx context.Context, dest, source string, idx int) error {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return err
			}

			objFile, err := getObjFile(ctx, source)
			if err != nil {
				return err
			}

			idx -= 1

			if idx < 0 && len(objFile.Vertices) <= idx {
				return fmt.Errorf("idx %d is out of bounds %d", idx, len(objFile.Vertices))
			}

			trait, err := getTriangle(shape)
			if err != nil {
				return err
			}
			if !trait.P2.Equal(objFile.Vertices[idx]) {
				return fmt.Errorf("expected p2 to be %s, got %s", objFile.Vertices[idx], trait.P2)
			}
			return nil
		})
	sc.Then(
		`^(\w+).p3 = (\w+).vertices\[(\d+)\]$`,
		func(ctx context.Context, dest, source string, idx int) error {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return err
			}

			objFile, err := getObjFile(ctx, source)
			if err != nil {
				return err
			}

			idx -= 1

			if idx < 0 && len(objFile.Vertices) <= idx {
				return fmt.Errorf("idx %d is out of bounds %d", idx, len(objFile.Vertices))
			}

			trait, err := getTriangle(shape)
			if err != nil {
				return err
			}
			if !trait.P3.Equal(objFile.Vertices[idx]) {
				return fmt.Errorf("expected p1 to be %s, got %s", objFile.Vertices[idx], trait.P3)
			}
			return nil
		})
	sc.Then(
		`^(\w+).n1 = (\w+).normals\[(\d+)\]$`,
		func(ctx context.Context, dest, source string, idx int) error {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return err
			}

			objFile, err := getObjFile(ctx, source)
			if err != nil {
				return err
			}

			idx -= 1

			if idx < 0 && len(objFile.Normals) <= idx {
				return fmt.Errorf("idx %d is out of bounds %d", idx, len(objFile.Vertices))
			}

			if !shape.Trait.(*rt.SmoothTriangle).N1.Equal(objFile.Normals[idx]) {
				return fmt.Errorf("expected p1 to be %s, got %s", objFile.Vertices[idx], shape.Trait.(*rt.SmoothTriangle).N1)
			}
			return nil
		})
	sc.Then(
		`^(\w+).n2 = (\w+).normals\[(\d+)\]$`,
		func(ctx context.Context, dest, source string, idx int) error {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return err
			}

			objFile, err := getObjFile(ctx, source)
			if err != nil {
				return err
			}

			idx -= 1

			if idx < 0 && len(objFile.Normals) <= idx {
				return fmt.Errorf("idx %d is out of bounds %d", idx, len(objFile.Vertices))
			}

			if !shape.Trait.(*rt.SmoothTriangle).N2.Equal(objFile.Normals[idx]) {
				return fmt.Errorf("expected p1 to be %s, got %s", objFile.Vertices[idx], shape.Trait.(*rt.SmoothTriangle).N2)
			}
			return nil
		})
	sc.Then(
		`^(\w+).n3 = (\w+).normals\[(\d+)\]$`,
		func(ctx context.Context, dest, source string, idx int) error {
			shape, err := getShape(ctx, dest)
			if err != nil {
				return err
			}

			objFile, err := getObjFile(ctx, source)
			if err != nil {
				return err
			}

			idx -= 1

			if idx < 0 && len(objFile.Normals) <= idx {
				return fmt.Errorf("idx %d is out of bounds %d", idx, len(objFile.Vertices))
			}

			if !shape.Trait.(*rt.SmoothTriangle).N3.Equal(objFile.Normals[idx]) {
				return fmt.Errorf("expected p1 to be %s, got %s", objFile.Vertices[idx], shape.Trait.(*rt.SmoothTriangle).N3)
			}
			return nil
		})

	sc.Then(
		`^(\w+) includes "(\w+)" from (\w+)$`,
		func(ctx context.Context, target, name, source string) error {
			shape, err := getShape(ctx, target)
			if err != nil {
				return err
			}

			objFile, err := getObjFile(ctx, source)
			if err != nil {
				return err
			}

			group := objFile.Groups[name]
			if !shape.Trait.(*rt.Group).Contains(group) {
				return fmt.Errorf("%s does not contain %s", shape, group)
			}
			return nil
		})

	sc.Then(
		`^t2 = t1$`,
		func(ctx context.Context) error {
			t1, err := getShape(ctx, "t1")
			if err != nil {
				return err
			}
			t2, err := getShape(ctx, "t2")
			if err != nil {
				return err
			}

			if !t2.Equal(t1) {
				return fmt.Errorf("expected %s to equal %s", t2, t1)
			}
			return nil
		})
}
