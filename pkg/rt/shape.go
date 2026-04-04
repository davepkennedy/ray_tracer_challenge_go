package rt

import (
	"fmt"
)

type ShapeTrait interface {
	Equal(t ShapeTrait) bool
	String() string
	Intersect(ray *Ray) []float64
	LocalNormalAt(point *Tuple) (*Tuple, error)
}

type Shape struct {
	Trait     ShapeTrait
	Transform *Matrix
	Material  *Material
}

func NewShape(trait ShapeTrait) *Shape {
	return &Shape{
		Trait:     trait,
		Transform: IdentityMatrix(),
		Material:  NewMaterial(),
	}
}

func (s *Shape) Equal(other *Shape) bool {
	return s.Transform.Equal(other.Transform) && s.Material.Equal(other.Material) && s.Trait.Equal(other.Trait)
}

func (s *Shape) String() string {
	return fmt.Sprintf("shape{%s %f}", s.Trait, s.Material.RefractiveIndex)
}

func (s *Shape) Intersect(ray *Ray) *Intersections {
	m, err := s.Transform.Inverse()
	if err != nil {
		return NewIntersections()
	}
	ray, err = ray.Transform(m)
	if err != nil {
		return NewIntersections()
	}

	ts := s.Trait.Intersect(ray)
	intersections := make([]*Intersection, 0)
	for _, t := range ts {
		intersections = append(intersections, NewIntersection(t, s))
	}
	return NewIntersections(intersections...)
}

/*
def normal_at(self, pt, i):
        local_point = self.world_to_object(pt)
        local_normal = self.local_normal_at(local_point, i)
        return self.normal_to_world(local_normal)

    def world_to_object(self, pt):
        if self.has_parent:
            pt = self.parent.world_to_object(pt)
        return self.transform.inverse() * pt

    def normal_to_world (self, normal):
        normal = self.transform.inverse().transpose() * normal
        normal = Vector(normal.x, normal.y, normal.z)
        normal = normal.normalize()

        if self.has_parent:
            normal = self.parent.normal_to_world(normal)

        return normal
*/

func (s *Shape) worldToObject(point *Tuple) (*Tuple, error) {
	inv, err := s.Transform.Inverse()
	if err != nil {
		return nil, err
	}
	return inv.MultiplyTuple(point)
}

func (s *Shape) normalToWorld(normal *Tuple) (*Tuple, error) {
	inv, err := s.Transform.Inverse()
	if err != nil {
		return nil, err
	}

	normal, err = inv.Transpose().MultiplyTuple(normal)
	if err != nil {
		return nil, err
	}

	normal = NewVector(normal.X, normal.Y, normal.Z)
	normal = normal.Normalize()

	/*
		if self.has_parent:
	            normal = self.parent.normal_to_world(normal)
	*/
	return normal, nil
}

func (s *Shape) NormalAt(point *Tuple) (*Tuple, error) {

	localPoint, err := s.worldToObject(point)
	if err != nil {
		return nil, err
	}

	localNormal, err := s.Trait.LocalNormalAt(localPoint)
	if err != nil {
		return nil, err
	}

	return s.normalToWorld(localNormal)
}
