package rt

import (
	"math"
)

type Triangle struct {
	P1, P2, P3 *Tuple
	E1, E2 *Tuple
	Normal *Tuple
}

func NewTriangle(p1, p2, p3 *Tuple) *Shape {
	e1 := p2.Subtract(p1)
    e2 := p3.Subtract(p1)
	normal := e2.Cross(e1).Normalize()
	return NewShape(
		&Triangle {
			P1:p1,
			P2:p2,
			P3:p3,
			E1:e1,
			E2:e2,
			Normal:normal,
		})
}

func (t *Triangle) LocalIntersect(ray Ray) Intersections {
	return Intersections{}
}

func (t *Triangle) Includes (s, other *Shape) bool{
	return s == other
}

func (t *Triangle) Equal(other any) bool {
	ot, ok := other.(*Triangle)
	if !ok {return false}

	return t.P1.Equal(ot.P1) &&
		t.P2.Equal(ot.P2) &&
		t.P3.Equal(ot.P3)
}

func (t *Triangle) String() string {
	return "triangle{}"
}

func (t *Triangle) AsCapped() *Capped {
	return nil
}

func (t *Triangle) Bounds() *BoundingBox {
	bounds := NewBoundingBox(t.P1, t.P1)
    bounds.AddPoint(t.P2)
	bounds.AddPoint(t.P3)
    return bounds
}

func (tri *Triangle) Intersect(shape *Shape, ray *Ray) *Intersections {
	dir_cross_e2 := ray.Direction.Cross(tri.E2)
    det := tri.E1.Dot(dir_cross_e2)
    if math.Abs(det) < EPSILON {
        return NewIntersections()
	}
    
	f := 1.0 / det
    p1_to_origin := ray.Origin.Subtract(tri.P1)
    u := f * p1_to_origin.Dot(dir_cross_e2)
    if u < 0 || u > 1 {
    	return NewIntersections()
	}
            
    origin_cross_e1 := p1_to_origin.Cross(tri.E1)
    v := f * ray.Direction.Dot(origin_cross_e1)
    if v < 0 || (u+v) > 1 {
    	return NewIntersections()
	}
            
    t := f * tri.E2.Dot(origin_cross_e1)
    return NewIntersections(
    		NewIntersectionWithUV(t, shape, u, v))
}

func (t *Triangle) LocalNormalAt(point *Tuple, i *Intersection) (*Tuple, error) {
	return t.Normal, nil
}