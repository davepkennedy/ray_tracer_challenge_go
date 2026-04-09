package test

import (
	"context"
	"fmt"
	"raytracer/pkg/rt"
)

type boolKey struct{ Name string }
type worldKey struct{ Name string }
type compsKey struct{ Name string }
type intKey struct{ Name string }
type floatKey struct{ Name string }
type cameraKey struct{ Name string }
type shapeKey struct{ Name string }
type patternKey struct{ Name string }
type objFileKey struct{ Name string }
type groupKey struct { Name string }

func getBoolean(ctx context.Context, name string) (bool, error) {
	val, ok := ctx.Value(boolKey{name}).(bool)
	if !ok {
		return false, fmt.Errorf("no bool named %s found", name)
	}
	return val, nil
}
func setBoolean(ctx context.Context, name string, val bool) context.Context {
	return context.WithValue(ctx, boolKey{name}, val)
}

func getCamera(ctx context.Context, name string) (*rt.Camera, error) {
	val, ok := ctx.Value(cameraKey{name}).(*rt.Camera)
	if !ok {
		return nil, fmt.Errorf("no camera named %s found", name)
	}
	return val, nil
}

func setCamera(ctx context.Context, name string, camera *rt.Camera) context.Context {
	return context.WithValue(ctx, cameraKey{name}, camera)
}

func getCanvas(ctx context.Context, name string) (*rt.Canvas, error) {
	val, ok := ctx.Value(canvasKey{name}).(*rt.Canvas)
	if !ok {
		return nil, fmt.Errorf("no canvas named %s found", name)
	}
	return val, nil
}

func getColor(ctx context.Context, name string) (*rt.Color, error) {
	val, ok := ctx.Value(colorKey{name}).(*rt.Color)
	if !ok {
		return nil, fmt.Errorf("no color named %s found", name)
	}
	return val, nil
}

func setColor(ctx context.Context, name string, color *rt.Color) context.Context {
	return context.WithValue(ctx, colorKey{name}, color)
}

func getComps(ctx context.Context, name string) (*rt.Computations, error) {
	val, ok := ctx.Value(compsKey{name}).(*rt.Computations)
	if !ok {
		return nil, fmt.Errorf("no computations name %s found", name)
	}
	return val, nil
}

func setComps(ctx context.Context, name string, comps *rt.Computations) context.Context {
	return context.WithValue(ctx, compsKey{name}, comps)
}

func getFloat(ctx context.Context, name string) (float64, error) {
	val, ok := ctx.Value(floatKey{name}).(float64)
	if !ok {
		return 0, fmt.Errorf("no float named %s found", name)
	}
	return val, nil
}

func setFloat(ctx context.Context, name string, val float64) context.Context {
	return context.WithValue(ctx, floatKey{name}, val)
}

func getInt(ctx context.Context, name string) (int, error) {
	val, ok := ctx.Value(intKey{name}).(int)
	if !ok {
		return 0, fmt.Errorf("no int named %s found", name)
	}
	return val, nil
}

func setInt(ctx context.Context, name string, val int) context.Context {
	return context.WithValue(ctx, intKey{name}, val)
}

func getIntersection(ctx context.Context, name string) (*rt.Intersection, error) {
	val, ok := ctx.Value(intersectionKey{name}).(*rt.Intersection)
	if !ok {
		return nil, fmt.Errorf("no intersection named %s found", name)
	}
	return val, nil
}
func setIntersection (ctx context.Context, name string, val *rt.Intersection) context.Context{
	return context.WithValue(ctx, intersectionKey{name}, val)
}

func getIntersections(ctx context.Context, name string) (*rt.Intersections, error) {
	val, ok := ctx.Value(intersectionsKey{name}).(*rt.Intersections)
	if !ok {
		return nil, fmt.Errorf("no intersections named %s found", name)
	}
	return val, nil
}

func setIntersections(ctx context.Context, name string, xs *rt.Intersections) context.Context {
	return context.WithValue(ctx, intersectionsKey{name}, xs)
}

func getLight(ctx context.Context, name string) (*rt.Light, error) {
	val, ok := ctx.Value(lightKey{name}).(*rt.Light)
	if !ok {
		return nil, fmt.Errorf("no light named %s found", name)
	}
	return val, nil
}

func getMaterial(ctx context.Context, name string) (*rt.Material, error) {
	val, ok := ctx.Value(materialKey{name}).(*rt.Material)
	if !ok {
		return nil, fmt.Errorf("no material named %s found", name)
	}
	return val, nil
}

func getMatrix(ctx context.Context, name string) (*rt.Matrix, error) {
	val, ok := ctx.Value(matrixKey{name}).(*rt.Matrix)
	if !ok {
		return nil, fmt.Errorf("no matrix named %s found", name)
	}
	return val, nil
}

func setMatrix(ctx context.Context, name string, matrix *rt.Matrix) context.Context {
	return context.WithValue(ctx, matrixKey{name}, matrix)
}

func getObjFile(ctx context.Context, name string) (*rt.ObjFile, error) {
	val, ok := ctx.Value(objFileKey{name}).(*rt.ObjFile)
	if !ok {
		return nil, fmt.Errorf("no obj file named %s found", name)
	}
	return val, nil
}

func setObjFile(ctx context.Context, name string, objFile *rt.ObjFile) context.Context {
	return context.WithValue(ctx, objFileKey{name}, objFile)
}

func getPattern(ctx context.Context, name string) (*rt.Pattern, error) {
	val, ok := ctx.Value(patternKey{name}).(*rt.Pattern)
	if !ok {
		return nil, fmt.Errorf("no pattern named %s found", name)
	}
	return val, nil
}

func setPattern(ctx context.Context, name string, pattern *rt.Pattern) context.Context {
	return context.WithValue(ctx, patternKey{name}, pattern)
}

func getRay(ctx context.Context, name string) (*rt.Ray, error) {
	val, ok := ctx.Value(rayKey{name}).(*rt.Ray)
	if !ok {
		return nil, fmt.Errorf("no ray named %s found", name)
	}
	return val, nil
}

func setRay(ctx context.Context, name string, ray *rt.Ray) context.Context {
	return context.WithValue(ctx, rayKey{name}, ray)
}

func getShape(ctx context.Context, name string) (*rt.Shape, error) {
	val, ok := ctx.Value(shapeKey{name}).(*rt.Shape)
	if !ok {
		return nil, fmt.Errorf("no shape named %s found", name)
	}
	return val, nil
}
func setShape(ctx context.Context, name string, shape *rt.Shape) context.Context {
	return context.WithValue(ctx, shapeKey{name}, shape)
}

func getString(ctx context.Context, name string) (string, error) {
	val, ok := ctx.Value(stringKey{name}).(string)
	if !ok {
		return "", fmt.Errorf("no string named %s found", name)
	}
	return val, nil
}
func setString (ctx context.Context, name string, val string) context.Context {
	return context.WithValue(ctx, stringKey{name}, val)
}

func getTuple(ctx context.Context, name string) (*rt.Tuple, error) {
	val, ok := ctx.Value(tupleKey{name}).(*rt.Tuple)
	if !ok {
		return nil, fmt.Errorf("no tuple named %s found", name)
	}
	return val, nil
}

func setTuple(ctx context.Context, name string, t *rt.Tuple) context.Context {
	return context.WithValue(ctx, tupleKey{name}, t)
}

func getWorld(ctx context.Context, name string) (*rt.World, error) {
	val, ok := ctx.Value(worldKey{name}).(*rt.World)
	if !ok {
		return nil, fmt.Errorf("no world named %s found", name)
	}
	return val, nil
}

func setWorld(ctx context.Context, name string, world *rt.World) context.Context {
	return context.WithValue(ctx, worldKey{name}, world)
}
