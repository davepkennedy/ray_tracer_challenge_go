package rt

import (
	"fmt"
	"sort"
)

type CSG struct {
	Operation	string
	Left 		*Shape
	Right 		*Shape
}

func NewCSG (operation string, s1, s2 *Shape) *Shape {
	s := NewShape (
		&CSG {
			Operation: operation,
			Left: s1,
			Right: s2,
		})

	s1.Parent = s
	s2.Parent = s
	return s
}

func IntersectionAllowed (operation string, lhit, inl, inr bool) bool {
	switch operation {
	case "union":
		return (lhit && !inr) || (!lhit && !inl)
	case "intersection":
		return (lhit && inr) || (!lhit && inl)
	case "difference":
		return (lhit && !inr) || (!lhit && inl)
	}
	return false
}

func (c *CSG) Equal (other any) bool {
	otherShape, ok := other.(*CSG)
	if !ok {return false}
	return c.Operation == otherShape.Operation &&
		c.Left.Equal(otherShape.Left) &&
		c.Right.Equal(otherShape.Right)
}

func (s *CSG) String() string {
	return fmt.Sprintf("csg{%s: (%s, %s)}", s.Operation, s.Left, s.Right)
}

func (c *CSG) Includes (s, other *Shape) bool {
	if s == other {
		return true
	}

	if c.Left.Includes(other) {return true}
	if c.Right.Includes(other) {return true}

	return false;
}

func (s *CSG) Intersect(shape* Shape, ray *Ray) *Intersections {
	xl := s.Left.Intersect(ray)
	xr := s.Right.Intersect(ray)

	combine := append(xl.intersections, xr.intersections...)
	
	sort.Slice(combine, func(x, y int) bool {
		return combine[x].T < combine[y].T
	})

	return s.FilterIntersections(NewIntersections(combine...))
}

func (s *CSG) LocalNormalAt(point *Tuple, i *Intersection) (*Tuple, error) {
	return nil, fmt.Errorf ("LocalNormalAt not defined for CSG")
}

func (s *CSG) AsCapped() *Capped {
	return nil
}

func (c *CSG) Bounds() *BoundingBox {
	bounds := c.Left.Bounds()
	bounds.AddBox(c.Right.Bounds())
	return bounds
}

func (c *CSG) FilterIntersections (xs *Intersections) *Intersections {
	inl := false
    inr := false
	lhit := false
        
    result := make([]*Intersection, 0)

    for _, i := range xs.intersections {
        lhit = c.Left.Includes(i.Object) 
        if IntersectionAllowed (c.Operation, lhit, inl, inr) {
            result = append(result, i)
		}
            
        if lhit {
            inl = !inl
		} else {
            inr = !inr
		}
	}
    return NewIntersections(result...)
}