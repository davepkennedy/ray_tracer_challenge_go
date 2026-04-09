package rt

import (
	"math"
	"reflect"
)

func NewCone() *Shape {
	return NewShape(&Cone{
		Capped: NewCapped(),
	})
}

type Cone struct {
	Capped
}

func (c *Cone) Equal(other ShapeTrait) bool {
	return reflect.TypeOf(c) == reflect.TypeOf(other)
}

func (c *Cone) String() string {
	return "cone{}"
}

func (cone *Cone) Intersect(s *Shape, ray *Ray) *Intersections {
	dx, dy, dz := ray.Direction.X, ray.Direction.Y, ray.Direction.Z
	ox, oy, oz := ray.Origin.X, ray.Origin.Y, ray.Origin.Z

	a := (dx * dx) - (dy * dy) + (dz * dz)
    b := 2 * ox * dx - 2 * oy * dy + 2 * oz * dz
	c := (ox * ox) - (oy * oy) + (oz * oz)

    disc := (b * b) - 4 * a * c
    if disc < 0 {
		return NewIntersections()
	}
	
	ts := make([]float64, 0)
	if a == 0 {
		ts = append(ts, -c / (2 * b))
	} else {
		t0 := (-b - math.Sqrt(disc)) / (2 * a)
		t1 := (-b + math.Sqrt(disc)) / (2 * a)
		if t0 > t1 {
    		t1, t0 = t0, t1
		}
		ts = append(ts, t0, t1)
	}

    xs := make([]*Intersection, 0)
	for _, t := range ts {
	    y0 := ray.Origin.Y + t * ray.Direction.Y
    	if cone.Minimum < y0 && y0 < cone.Maximum {
    		xs = append(xs, NewIntersection(t, s))
		}
	}
            
    return NewIntersections(cone.intersectCaps(s, ray, xs)...)
}
        
func (cone *Cone) intersectCaps(s *Shape, r *Ray, xs []*Intersection) []*Intersection{
    if !cone.Closed || math.Abs(r.Direction.Y) < EPSILON {
        return xs
	}
            
    for _, v := range []float64{cone.Minimum, cone.Maximum} {
        t := (v - r.Origin.Y) / r.Direction.Y
            
        if cone.checkCap(r, t, math.Abs(v)) {
            xs = append(xs, NewIntersection(t, s))
        }
    }
	return xs
}

func (cone *Cone) checkCap(r *Ray, t float64, radius float64) bool {
	x := r.Origin.X + t * r.Direction.X
    z := r.Origin.Z + t * r.Direction.Z
        
    return (x * x + z * z) <= radius
}

func (cone *Cone) LocalNormalAt(pt *Tuple, _ *Intersection) (*Tuple, error) {
	dist := (pt.X * pt.X) + (pt.Z * pt.Z)
        
	if dist < 1 && pt.Y >= (cone.Maximum - EPSILON) {
    	return NewVector(0, 1, 0), nil
	}
            
    if dist < 1 && pt.Y <= cone.Minimum + EPSILON{
    	return NewVector(0, -1, 0), nil
	}
            
    y := math.Sqrt(pt.X * pt.X + pt.Z * pt.Z)
    if pt.Y > 0{
    	y = -y
	}
            
    return NewVector(pt.X, y, pt.Z), nil
}

func (cone *Cone) AsCapped() *Capped {
	return &cone.Capped
}

func (cone *Cone) Bounds() *BoundingBox {
	limit := math.Max(math.Abs(cone.Minimum), math.Abs(cone.Maximum))
	return NewBoundingBox(
		NewPoint(-limit, cone.Minimum, -limit),
		NewPoint(limit, cone.Maximum, limit))
}