package rt

type SmoothTriangle struct {
	Triangle
	N1	*Tuple
	N2	*Tuple
	N3	*Tuple
}

func NewSmoothTriangle (p1, p2, p3, n1, n2, n3 *Tuple) *Shape {
	t := NewTriangle(p1,p2,p3)
	return NewShape(&SmoothTriangle{
		Triangle: *t.Trait.(*Triangle),
		N1: n1,
		N2: n2,
		N3: n3,
	})
}

func (t *SmoothTriangle) String() string {
	return "smmothtriangle{}"
}

func (t *SmoothTriangle) Equal (other ShapeTrait) bool {
	ot, ok := other.(*SmoothTriangle)
	if !ok {return false}

	return t.P1.Equal(ot.P1) &&
		t.P2.Equal(ot.P2) &&
		t.P3.Equal(ot.P3) &&
		t.N1.Equal(ot.N1) &&
		t.N2.Equal(ot.N2) &&
		t.N3.Equal(ot.N3)
}

func (t *SmoothTriangle) LocalNormalAt(point *Tuple, i *Intersection) (*Tuple, error) {
	a := t.N2.MultiplyScalar(i.U)
	b := t.N3.MultiplyScalar(i.V)
	c := t.N1.MultiplyScalar(1 - i.U - i.V)
	return a.Add(b).Add(c), nil
}