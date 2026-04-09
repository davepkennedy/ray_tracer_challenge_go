package rt

import (
	"math"
	"reflect"
)

type Cube struct {}

func MaxA [T int | float64](a ...T) T {
	var max T
	
	if len(a) == 0 {
		return max
	}
	max = a[0]
	for _, v := range a {
		if v > max {
			max = v
		}
	}
	return max
}

func MinA [T int | float64](a ...T) T {
	var min T

	if len(a) == 0 {
		return min
	}
	min = a[0]

	for _, v := range a {
		if v < min {
			min = v
		}
	}
	return min
}

func NewCube() *Shape {
	return NewShape(&Cube{})
}

func (c *Cube) Equal(other ShapeTrait) bool {
	return reflect.TypeOf(c) == reflect.TypeOf(other)
}

func (s *Cube) String() string {
	return "cube{}"
}

func checkAxis(origin, direction float64) (float64, float64) {
	tminNumerator := (-1 - origin)
	tmaxNumerator := (1 - origin)

	var tmin, tmax float64

	if math.Abs(direction) >= EPSILON {
		tmin = tminNumerator / direction
        tmax = tmaxNumerator / direction
	} else {
        tmin = tminNumerator * math.Inf(1)
        tmax = tmaxNumerator * math.Inf(1)
	}
            
	if tmin > tmax {
		return tmax, tmin
	}
    return tmin, tmax
}

func (s *Cube) Intersect(shape* Shape, ray *Ray) *Intersections {
	xtmin, xtmax := checkAxis(ray.Origin.X, ray.Direction.X)
    ytmin, ytmax := checkAxis(ray.Origin.Y, ray.Direction.Y)
    ztmin, ztmax := checkAxis(ray.Origin.Z, ray.Direction.Z)
        
    tmin := MaxA(xtmin, ytmin, ztmin)
    tmax := MinA(xtmax, ytmax, ztmax)

	if (tmin > tmax) {
        return NewIntersections()
	} 
	
	return NewIntersections(
		NewIntersection(tmin, shape), 
		NewIntersection(tmax, shape))
}

func (s *Cube) LocalNormalAt(point *Tuple, _ *Intersection) (*Tuple, error) {
	maxc := MaxA(math.Abs(point.X), math.Abs(point.Y), math.Abs(point.Z))
        
    if maxc == math.Abs(point.X) {
            return NewVector(point.X, 0, 0), nil
	} else if maxc == math.Abs(point.Y){
            return NewVector(0, point.Y, 0), nil
	}
    return NewVector(0, 0, point.Z), nil
}

func (s *Cube) AsCapped() *Capped {
	return nil
}

func (c *Cube) Bounds() *BoundingBox {
	return NewBoundingBox(
		NewPoint(-1, -1, -1),
	    NewPoint(1, 1, 1))
	}