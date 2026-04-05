package rt

import (
	"math"
)

type BoundingBox struct {
	Min *Tuple
	Max *Tuple
}

func NewBoundingBox(min, max *Tuple) *BoundingBox {
	return &BoundingBox{
		Min: min,
		Max: max,
	}
}

func (b *BoundingBox) AddPoint(point *Tuple) {
	if point.X < b.Min.X { b.Min = NewPoint(point.X, b.Min.Y, b.Min.Z) }
	if point.X > b.Max.X { b.Max = NewPoint(point.X, b.Max.Y, b.Max.Z) }
	if point.Y < b.Min.Y { b.Min = NewPoint(b.Min.X, point.Y, b.Min.Z) }
	if point.Y > b.Max.Y { b.Max = NewPoint(b.Max.X, point.Y, b.Max.Z) }
	if point.Z < b.Min.Z { b.Min = NewPoint(b.Min.X, b.Min.Y, point.Z) }
    if point.Z > b.Max.Z { b.Max = NewPoint(b.Max.X, b.Max.Y, point.Z) }
}

func (b *BoundingBox) AddBox(box *BoundingBox) {
	b.AddPoint(box.Min)
	b.AddPoint(box.Max)
}

func (b *BoundingBox) ContainsTuple(other *Tuple) bool {
	contains_x := b.Min.X <= other.X && other.X <= b.Max.X
	contains_y := b.Min.Y <= other.Y && other.Y <= b.Max.Y
	contains_z := b.Min.Z <= other.Z && other.Z <= b.Max.Z
	return contains_x && contains_y && contains_z
}

func (b *BoundingBox) ContainsBox(box *BoundingBox) bool {
	contains_min := b.ContainsTuple(box.Min)
	contains_max := b.ContainsTuple(box.Max)
	return contains_min && contains_max
}

func (b *BoundingBox) Transform(matrix *Matrix) (*BoundingBox, error) {
	points := []*Tuple{
		b.Min,
		NewPoint(b.Min.X, b.Min.Y, b.Max.Z),
		NewPoint(b.Min.X, b.Max.Y, b.Min.Z),
		NewPoint(b.Min.X, b.Max.Y, b.Max.Z),
		NewPoint(b.Max.X, b.Min.Y, b.Min.Z),
		NewPoint(b.Max.X, b.Min.Y, b.Max.Z),
		NewPoint(b.Max.X, b.Max.Y, b.Min.Z),
		b.Max,
	}

	box := NewBoundingBox(b.Min, b.Max)
	for _, pt := range points {
		transformed, err := matrix.MultiplyTuple(pt)
		if err != nil {return nil, err}
		box.AddPoint(transformed)
	}

	return box, nil
}

func (b *BoundingBox) checkAxis(min, max, origin, direction float64) (float64, float64) {
    tmin_numerator := (min - origin)
	tmax_numerator := (max - origin)

	tmin, tmax := 0.0, 0.0
        
    if math.Abs(direction) >= EPSILON {
        tmin = tmin_numerator / direction
        tmax = tmax_numerator / direction
    } else {
        tmin = tmin_numerator * math.Inf(1)
		tmax = tmax_numerator * math.Inf(1)
	}

    if tmin > tmax {
		tmin, tmax = tmax, tmin
	}

    return tmin, tmax
}
        
func (b *BoundingBox) Intersects(ray *Ray) bool {
        xtmin, xtmax := b.checkAxis(b.Min.X, b.Max.X, ray.Origin.X, ray.Direction.X)
        ytmin, ytmax := b.checkAxis(b.Min.Y, b.Max.Y, ray.Origin.Y, ray.Direction.Y)
        ztmin, ztmax := b.checkAxis(b.Min.Z, b.Max.Z, ray.Origin.Z, ray.Direction.Z)

        tmin := math.Max(xtmin, math.Max(ytmin, ztmin))
        tmax := math.Min(xtmax, math.Min(ytmax, ztmax))

		if tmin > tmax { return false } 
		return true
}

func (b *BoundingBox) Split() (*BoundingBox, *BoundingBox) {
        // figure out the box's largest dimension
        dx := b.Max.X - b.Min.X
        dy := b.Max.Y - b.Min.Y
        dz := b.Max.Z - b.Min.Z

        greatest := MaxA(dx, dy, dz)

    	// variables to help construct the points on
     	// the dividing plane
        x0 := b.Min.X
        y0 := b.Min.Y
        z0 := b.Min.Z

        x1 := b.Max.X
        y1 := b.Max.Y
        z1 := b.Max.Z

        // adjust the points so that they lie on the
        // dividing plane
        if greatest == dx {
            x0 = x0 + dx / 2.0
			x1 = x0
		} else if greatest == dy {
            y0 = y0 + dy / 2.0
            y1 = y0
        } else {
            z0 = z0 + dz / 2.0
			z1 = z0
        }

        mid_min := NewPoint(x0, y0, z0)
        mid_max := NewPoint(x1, y1, z1)

        // construct and return the two halves of
        // the bounding box
        return NewBoundingBox(b.Min, mid_max), NewBoundingBox(mid_min, b.Max)
}