package rt

import (
	"math"
	"reflect"
)

type PlaneTrait struct {
}

func NewPlane() *Shape {
	return NewShape(PlaneTrait{})
}

func (p PlaneTrait) Equal(other ShapeTrait) bool {
	return reflect.TypeOf(p) == reflect.TypeOf(other)
}

func (p PlaneTrait) String() string {
	return "plane{}"
}

func (p PlaneTrait) Intersect(ray *Ray) []float64 {
	if math.Abs(ray.Direction.Y) < EPSILON {
    	return []float64{}
	}
    t := -ray.Origin.Y / ray.Direction.Y
    return []float64{t}
}

func (p PlaneTrait) LocalNormalAt (point *Tuple)  (*Tuple, error) {
	return NewVector(0,1,0), nil
}