package rt

import (
	"math"
)

type PatternTrait interface {
	ColorAt(pt *Tuple, a, b *Color) *Color
}

type Pattern struct {
	A         *Color
	B         *Color
	Transform *Matrix
	Trait     PatternTrait
}

type StripePattern struct{}
type GradientPattern struct{}
type RingPattern struct{}
type CheckersPattern struct{}

func NewPattern(a, b *Color, trait PatternTrait) *Pattern {
	return &Pattern{
		A:         a,
		B:         b,
		Transform: IdentityMatrix(),
		Trait:     trait,
	}
}

func (p *Pattern) ColorAt(pt *Tuple) *Color {
	return p.Trait.ColorAt(pt, p.A, p.B)
}

func (p *Pattern) ColorAtObject(shape *Shape, pt *Tuple) (*Color, error) {
	inv, err := shape.Transform.Inverse()
	if err != nil {
		return nil, err
	}

	objectPoint, err := inv.MultiplyTuple(pt)
	if err != nil {
		return nil, err
	}
	inv, err = p.Transform.Inverse()
	if err != nil {
		return nil, err
	}

	patternPoint, err := inv.MultiplyTuple(objectPoint)
	if err != nil {
		return nil, err
	}

	return p.ColorAt(patternPoint), err
}

func NewStripePattern(a, b *Color) *Pattern {
	return NewPattern(a, b, StripePattern{})
}

func (p StripePattern) ColorAt(pt *Tuple, a, b *Color) *Color {
	if int(math.Floor(pt.X))%2 == 0 {
		return a
	}
	return b
}

func NewGradientPattern(a, b *Color) *Pattern {
	return NewPattern(a, b, GradientPattern{})
}

func (p GradientPattern) ColorAt(pt *Tuple, a, b *Color) *Color {
	distance := b.Subtract(a)
	fraction := pt.X - math.Floor(pt.X)
	return a.Add(distance.MultiplyScalar(fraction))
}

func NewRingPattern(a, b *Color) *Pattern {
	return NewPattern(a, b, RingPattern{})
}

func (p RingPattern) ColorAt(pt *Tuple, a, b *Color) *Color {
	xs := pt.X * pt.X
	zs := pt.Z * pt.Z
	if (int(math.Floor(math.Sqrt(xs+zs))) % 2) == 0 {
		return a
	}
	return b
}

func NewCheckersPattern(a, b *Color) *Pattern {
	return NewPattern(a, b, CheckersPattern{})
}

func (p CheckersPattern) ColorAt(pt *Tuple, a, b *Color) *Color {
	if int(math.Floor(pt.X)+math.Floor(pt.Y)+math.Floor(pt.Z))%2 == 0 {
		return a
	}
	return b
}
