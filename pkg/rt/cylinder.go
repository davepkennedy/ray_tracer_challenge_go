package rt

import (
	"math"
	"reflect"
)

type Cylinder struct {
	Capped
}

func NewCylinder() *Shape {
	return NewShape(&Cylinder{
		Capped: NewCapped(),
	})
}

func (c *Cylinder) Equal(other any) bool {
	return reflect.TypeOf(c) == reflect.TypeOf(other)
}

func (c *Cylinder) String() string {
	return "cylinder{}"
}

func (c *Cylinder) Includes (s, other *Shape) bool{
	return s == other
}

func (cyl *Cylinder) Intersect(s* Shape, ray *Ray) *Intersections {
	a := (ray.Direction.X * ray.Direction.X) + (ray.Direction.Z * ray.Direction.Z)
    b := 2 * ray.Origin.X * ray.Direction.X + 2 * ray.Origin.Z * ray.Direction.Z
	c := (ray.Origin.X * ray.Origin.X) + (ray.Origin.Z * ray.Origin.Z) - 1

    disc := (b * b) - 4 * a * c
    if disc < 0 {
		return NewIntersections()
	}
	
	var t0, t1 float64
	if a == 0 {
		t0 = math.NaN()
		t1 = math.NaN()
	} else {
		t0 = (-b - math.Sqrt(disc)) / (2 * a)
		t1 = (-b + math.Sqrt(disc)) / (2 * a)
	}

    if t0 > t1 {
    	t1, t0 = t0, t1
	}

    xs := make([]*Intersection, 0)
    y0 := ray.Origin.Y + t0 * ray.Direction.Y
    if cyl.Minimum < y0 && y0 < cyl.Maximum {
    	xs = append(xs, NewIntersection(t0, s))
	}
            
    y1 := ray.Origin.Y + t1 * ray.Direction.Y
    if cyl.Minimum < y1 && y1 < cyl.Maximum {
        xs = append(xs, NewIntersection(t1, s))
	}
            
    return NewIntersections(cyl.intersectCaps(s, ray, xs)...)
}

func (c *Cylinder) LocalNormalAt(pt *Tuple, _ *Intersection) (*Tuple, error) {
	dist := (pt.X * pt.X) + (pt.Z * pt.Z)

    if dist < 1 && pt.Y >= (c.Maximum - EPSILON){
        return NewVector(0, 1, 0), nil
	}
      
	if dist < 1 && pt.Y <= c.Minimum + EPSILON {
		return NewVector(0, -1, 0), nil
	}
	
	return NewVector(pt.X, 0, pt.Z), nil
}

 func (c *Cylinder) intersectCaps(s *Shape, r *Ray, xs []*Intersection) []*Intersection {
    if !c.Closed || math.Abs(r.Direction.Y) < EPSILON {	
        return xs
	}

    for _, v := range []float64{c.Minimum, c.Maximum} {
        t := (v - r.Origin.Y) / r.Direction.Y
        if c.checkCap(r, t) {
            xs = append(xs, NewIntersection(t, s))
		}
	}
    return xs
}

func (c *Cylinder) checkCap(r *Ray, t float64) bool {
    x := r.Origin.X + t * r.Direction.X
    z := r.Origin.Z + t * r.Direction.Z
    return (x * x + z * z) <= 1
}

func (c *Cylinder) AsCapped() *Capped {
	return &c.Capped
}

func (c *Cylinder) Bounds() *BoundingBox {
	return NewBoundingBox(
		NewPoint(-1, c.Minimum, -1),
		NewPoint(1, c.Maximum, 1))
}