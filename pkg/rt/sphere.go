package rt

import (
	//"log"
	"math"
	"reflect"
)

type Sphere struct {
}

func NewSphere() *Shape {
	return NewShape(&Sphere{})
}

func GlassSphere() *Shape {
	s := NewSphere()
	s.Material.Transparency = 1.0
	s.Material.RefractiveIndex = 1.5
	return s
}

func (s *Sphere) Equal(other any) bool {
	return reflect.TypeOf(s) == reflect.TypeOf(other)
}

func (s *Sphere) String() string {
	return "sphere{}"
}

func (_ *Sphere) Includes (s, other *Shape) bool{
	return s == other
}

func (s *Sphere) Intersect(shape *Shape, ray *Ray) *Intersections {
	sphereToRay := ray.Origin.Subtract(NewPoint(0.0, 0.0, 0.0))

	a := ray.Direction.Dot(ray.Direction)
	b := 2.0 * ray.Direction.Dot(sphereToRay)
	c := sphereToRay.Dot(sphereToRay) - 1.0

	discriminant := (b * b) - 4.0*a*c

	if discriminant < 0.0 {
		return NewIntersections()
	}

	t1 := (-b - math.Sqrt(discriminant)) / (2.0 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2.0 * a)
	return NewIntersections(NewIntersection(t1, shape), NewIntersection(t2, shape))
}

func (s *Sphere) LocalNormalAt(point *Tuple, _ *Intersection) (*Tuple, error) {
	return point.Subtract(NewPoint(0, 0, 0)), nil
}

func (s *Sphere) AsCapped() *Capped {
	return nil
}

func (s *Sphere) Bounds() *BoundingBox {
	return NewBoundingBox(
		NewPoint(-1, -1, -1),
		NewPoint(1, 1, 1))
}