package rt

import "testing"

func TestN1AndN2(t *testing.T) {
	a := GlassSphere()
	a.Transform = Scaling(2, 2, 2)
	a.Material = NewMaterial()
	a.Material.RefractiveIndex = 1.5

	b := GlassSphere()
	b.Transform = Translation(0, 0, -0.25)
	b.Material = NewMaterial()
	b.Material.RefractiveIndex = 2.0

	c := GlassSphere()
	c.Transform = Translation(0, 0, 0.25)
	c.Material = NewMaterial()
	c.Material.RefractiveIndex = 2.5

	tests := []struct {
		index int
		n1    float64
		n2    float64
	}{
		{0, 1.0, 1.5},
		{1, 1.5, 2.0},
		{2, 2.0, 2.5},
		{3, 2.5, 2.5},
		{4, 2.5, 1.5},
		{5, 1.5, 1.0},
	}

	ray := NewRay(NewPoint(0, 0, -4), NewVector(0, 0, 1))
	intersections := NewIntersections(
		NewIntersection(2, a),
		NewIntersection(2.75, b),
		NewIntersection(3.25, c),
		NewIntersection(4.75, b),
		NewIntersection(5.25, c),
		NewIntersection(6, a),
	)

	for _, test := range tests {
		comps, err := intersections.At(test.index).PrepareComputations(ray, intersections)
		if err != nil {
			t.Errorf("error preparing computations at index %d: %v", test.index, err)
			continue
		}
		if comps.N1 != test.n1 {
			t.Errorf("n1 mismatch at index %d: expected %f, got %f", test.index, test.n1, comps.N1)
		}
		if comps.N2 != test.n2 {
			t.Errorf("n2 mismatch at index %d: expected %f, got %f", test.index, test.n2, comps.N2)
		}
	}
}
