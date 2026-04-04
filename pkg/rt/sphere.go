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

func (s *Sphere) Equal(other ShapeTrait) bool {
	return reflect.TypeOf(s) == reflect.TypeOf(other)
}

func (s *Sphere) String() string {
	return "sphere{}"
}

func (s *Sphere) Intersect(ray *Ray) []float64 {
	sphereToRay := ray.Origin.Subtract(NewPoint(0.0, 0.0, 0.0))

	a := ray.Direction.Dot(ray.Direction)
	b := 2.0 * ray.Direction.Dot(sphereToRay)
	c := sphereToRay.Dot(sphereToRay) - 1.0

	discriminant := (b * b) - 4.0*a*c

	if discriminant < 0.0 {
		return []float64{}
	}

	t1 := (-b - math.Sqrt(discriminant)) / (2.0 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2.0 * a)
	return []float64{t1, t2}
}

func (s *Sphere) LocalNormalAt(point *Tuple) (*Tuple, error) {
	return point.Subtract(NewPoint(0, 0, 0)), nil
}
