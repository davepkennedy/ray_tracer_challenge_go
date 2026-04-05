package rt

import (
	"math"
	"reflect"
)

type PlaneTrait struct {
}

func NewPlane() *Shape {
	return NewShape(&PlaneTrait{})
}

func (p *PlaneTrait) Equal(other ShapeTrait) bool {
	return reflect.TypeOf(p) == reflect.TypeOf(other)
}

func (p *PlaneTrait) String() string {
	return "plane{}"
}

func (p *PlaneTrait) Intersect(s *Shape, ray *Ray) *Intersections {
	if math.Abs(ray.Direction.Y) < EPSILON {
		return NewIntersections()
	}
	t := -ray.Origin.Y / ray.Direction.Y
	return NewIntersections(NewIntersection(t, s))
}

func (p *PlaneTrait) LocalNormalAt(point *Tuple) (*Tuple, error) {
	return NewVector(0, 1, 0), nil
}

func (p *PlaneTrait) AsCapped() *Capped {
	return nil
}

func (p *PlaneTrait) Bounds() *BoundingBox {
	return NewBoundingBox(
		NewPoint(math.Inf(-1), 0, math.Inf(-1)),
		NewPoint(math.Inf(1), 0, math.Inf(1)))
}