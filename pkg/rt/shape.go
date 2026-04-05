package rt

import (
	"fmt"
	"math"
)

type Capped struct{
	Closed bool
	Minimum float64
	Maximum float64
}

func NewCapped () Capped {
	return Capped{
		Closed: false,
		Minimum: math.Inf(-1),
		Maximum: math.Inf(1),
	}
}

type ShapeTrait interface {
	Equal(t ShapeTrait) bool
	String() string
	Intersect(s* Shape, ray *Ray) *Intersections
	LocalNormalAt(point *Tuple) (*Tuple, error)
	AsCapped() *Capped
	Bounds() *BoundingBox
}

type Shape struct {
	Trait     ShapeTrait
	Transform *Matrix
	Material  *Material
	Parent    *Shape
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

	return s.Trait.Intersect(s, ray)}

func (s *Shape) WorldToObject(point *Tuple) (*Tuple, error) {
    if s.Parent != nil {
		var err error
		point, err = s.Parent.WorldToObject(point)
		if err != nil {return nil, err}
	}

	inv, err := s.Transform.Inverse()
	if err != nil {
		return nil, err
	}
	return inv.MultiplyTuple(point)
}

func (s *Shape) NormalToWorld(normal *Tuple) (*Tuple, error) {
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

	
    if s.Parent != nil {
		var err error
		normal, err = s.Parent.NormalToWorld(normal)
		if err != nil {return nil, err}
	}
	return normal, nil
}

func (s *Shape) NormalAt(point *Tuple) (*Tuple, error) {

	localPoint, err := s.WorldToObject(point)
	if err != nil {
		return nil, err
	}

	localNormal, err := s.Trait.LocalNormalAt(localPoint)
	if err != nil {
		return nil, err
	}

	return s.NormalToWorld(localNormal)
}

func (s *Shape) AddChild(child *Shape) {
	child.Parent = s
	trait, ok := s.Trait.(*Group)
	if !ok {return}
	trait.Children = append(trait.Children, child)
}

func (s *Shape) Bounds() *BoundingBox {
	return s.Trait.Bounds()
}

func (s *Shape) ParentSpaceBounds() (*BoundingBox, error) {
	b := s.Bounds()
	return b.Transform(s.Transform)
}