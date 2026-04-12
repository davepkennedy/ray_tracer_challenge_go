package test

import (
	"reflect"
	"context"
	"fmt"

	"raytracer/pkg/rt"
)

type TestState = context.Context

type valueKey[T any] struct{ name string }

func getValue[T any](state TestState, name string) (T, error) {
	a := state.Value(name)
	val, ok := a.(T)
	if !ok {
		if a != nil {
			return val, fmt.Errorf("no %s named %s in state; did find %s (%s)",
			 	reflect.TypeFor[T](), name, a, reflect.TypeOf(a))	
		}
		return val, fmt.Errorf("no %s named %s in state", reflect.TypeFor[T](), name)
	}
	return val, nil
}

func setValue[T any](state TestState, name string, value T) TestState {
	state = context.WithValue(state, name, value)
	return state
}

var (
	getBoolean = getValue[bool]
	setBoolean = setValue[bool]

	getCamera = getValue[*rt.Camera]
	setCamera = setValue[*rt.Camera]

	getCanvas = getValue[*rt.Canvas]
	setCanvas = setValue[*rt.Canvas]

	getColor = getValue[*rt.Color]
	setColor = setValue[*rt.Color]

	getComps = getValue[*rt.Computations]
	setComps = setValue[*rt.Computations]

	getFloat = getValue[float64]
	setFloat = setValue[float64]

	getInt = getValue[int]
	setInt = setValue[int]

	getIntersection = getValue[*rt.Intersection]
	setIntersection = setValue[*rt.Intersection]

	getIntersections = getValue[*rt.Intersections]
	setIntersections = setValue[*rt.Intersections]

	getLight = getValue[*rt.Light]
	setLight = setValue[*rt.Light]

	getMaterial = getValue[*rt.Material]
	setMaterial = setValue[*rt.Material]

	getMatrix = getValue[*rt.Matrix]
	setMatrix = setValue[*rt.Matrix]

	getObjFile = getValue[*rt.ObjFile]
	setObjFile = setValue[*rt.ObjFile]

	getPattern = getValue[*rt.Pattern]
	setPattern = setValue[*rt.Pattern]

	getRay = getValue[*rt.Ray]
	setRay = setValue[*rt.Ray]

	getShape = getValue[*rt.Shape]
	setShape = setValue[*rt.Shape]

	getString = getValue[string]
	setString = setValue[string]

	getTuple = getValue[*rt.Tuple]
	setTuple = setValue[*rt.Tuple]

	getWorld = getValue[*rt.World]
	setWorld = setValue[*rt.World]
)
