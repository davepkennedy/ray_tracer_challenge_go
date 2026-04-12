package rt

import (
	"fmt"
	"math"
)

var (
	_ Equality = (*Tuple)(nil)
	_ Equality = (*Color)(nil)
)

type Tuple struct {
	X, Y, Z, W float64
}

type Color struct {
	Red, Green, Blue float64
}

func (t *Tuple) String() string {
	return fmt.Sprintf("Tuple(X: %f, Y: %f, Z: %f, W: %f)", t.X, t.Y, t.Z, t.W)
}

func (c *Color) String() string {
	return fmt.Sprintf("Color(Red: %f, Green: %f, Blue: %f)", c.Red, c.Green, c.Blue)
}

func (t *Tuple) IsPoint() bool {
	return !t.IsVector()
}

func (t *Tuple) IsVector() bool {
	return t.W == 0.0
}

func NewTuple(x, y, z, w float64) *Tuple {
	return &Tuple{X: x, Y: y, Z: z, W: w}
}

func NewPoint(x, y, z float64) *Tuple {
	return NewTuple(x, y, z, 1)
}

func NewVector(x, y, z float64) *Tuple {
	return NewTuple(x, y, z, 0)
}

func NewColor(r, g, b float64) *Color {
	return &Color{Red: r, Green: g, Blue: b}
}

func (t *Tuple) Equal(other any) bool {
	ot, ok := other.(*Tuple)
	if !ok {return false}

	return math.Abs(t.X-ot.X) < EPSILON &&
		math.Abs(t.Y-ot.Y) < EPSILON &&
		math.Abs(t.Z-ot.Z) < EPSILON &&
		math.Abs(t.W-ot.W) < EPSILON
}

func (c *Color) Equal(other any) bool {
	oc, ok := other.(*Color)
	if !ok {return false}

	return math.Abs(c.Red-oc.Red) < EPSILON &&
		math.Abs(c.Green-oc.Green) < EPSILON &&
		math.Abs(c.Blue-oc.Blue) < EPSILON
}

func (t *Tuple) Negate() *Tuple {
	return NewTuple(-t.X, -t.Y, -t.Z, -t.W)
}

func (t *Tuple) Add(other *Tuple) *Tuple {
	return NewTuple(
		t.X+other.X,
		t.Y+other.Y,
		t.Z+other.Z,
		t.W+other.W,
	)
}

func (t *Tuple) Subtract(other *Tuple) *Tuple {
	return NewTuple(
		t.X-other.X,
		t.Y-other.Y,
		t.Z-other.Z,
		t.W-other.W,
	)
}

func (t *Tuple) MultiplyScalar(scalar float64) *Tuple {
	return NewTuple(
		t.X*scalar,
		t.Y*scalar,
		t.Z*scalar,
		t.W*scalar,
	)
}

func (t *Tuple) DivideScalar(scalar float64) *Tuple {
	return NewTuple(
		t.X/scalar,
		t.Y/scalar,
		t.Z/scalar,
		t.W/scalar,
	)
}

func (t *Tuple) Magnitude() float64 {
	return math.Sqrt(t.X*t.X + t.Y*t.Y + t.Z*t.Z + t.W*t.W)
}

func (t *Tuple) Normalize() *Tuple {
	mag := t.Magnitude()
	return t.DivideScalar(mag)
}

func (t *Tuple) Dot(other *Tuple) float64 {
	return t.X*other.X + t.Y*other.Y + t.Z*other.Z + t.W*other.W
}

func (t *Tuple) Cross(other *Tuple) *Tuple {
	return NewVector(
		t.Y*other.Z-t.Z*other.Y,
		t.Z*other.X-t.X*other.Z,
		t.X*other.Y-t.Y*other.X,
	)
}

func (t *Tuple) Reflect(normal *Tuple) *Tuple {
	return t.Subtract(normal.MultiplyScalar(2 * t.Dot(normal)))
}

func (c *Color) Add(other *Color) *Color {
	return NewColor(
		c.Red+other.Red,
		c.Green+other.Green,
		c.Blue+other.Blue)
}

func (c *Color) Subtract(other *Color) *Color {
	return NewColor(
		c.Red-other.Red,
		c.Green-other.Green,
		c.Blue-other.Blue)
}

func (c *Color) MultiplyScalar(scalar float64) *Color {
	return NewColor(
		c.Red*scalar,
		c.Green*scalar,
		c.Blue*scalar)
}

func (c *Color) Multiply(other *Color) *Color {
	return NewColor(
		c.Red*other.Red,
		c.Green*other.Green,
		c.Blue*other.Blue)
}
