package rt

import (
	"fmt"
	"slices"
	"sort"
)

type Intersection struct {
	T      float64
	Object *Shape
}

type Intersections struct {
	intersections []*Intersection
}

func NewIntersections(i ...*Intersection) *Intersections {
	return &Intersections{i}
}

func NewIntersection(t float64, object *Shape) *Intersection {
	return &Intersection{
		T:      t,
		Object: object,
	}
}

func filter(i []*Intersection) []*Intersection {
	c := make([]*Intersection, 0)
	for _, intersection := range i {
		if intersection.T >= 0 {
			c = append(c, intersection)
		}
	}
	return c
}

func (i *Intersection) Equal(other *Intersection) bool {
	return i.Object.Equal(other.Object) && i.T == other.T
}

func (i *Intersection) String() string {
	return fmt.Sprintf("intersection{s: %s, t: %f}", i.Object, i.T)
}

func (i *Intersections) At(idx int) *Intersection {
	if idx < 0 || idx >= len(i.intersections) {
		return nil
	}
	return i.intersections[idx]
}

func (i *Intersections) String() string {
	s := "intersections{"
	for _, intersection := range i.intersections {
		s += fmt.Sprintf("%s,\n", intersection)
	}
	s += "}"
	return s
}

func (i *Intersections) Len() int {
	return len(i.intersections)
}

func (i *Intersections) Hit() *Intersection {
	c := filter(i.intersections)
	if len(c) == 0 {
		return nil
	}

	sort.Slice(c, func(x, y int) bool {
		return c[x].T < c[y].T
	})

	return c[0]
}

type Computations struct {
	Object     *Shape
	T          float64
	Point      *Tuple
	Eye        *Tuple
	Normal     *Tuple
	Inside     bool
	OverPoint  *Tuple
	UnderPoint *Tuple
	Reflect    *Tuple
	N1         float64
	N2         float64
}

func (i *Intersection) PrepareComputations(r *Ray, xs *Intersections) (*Computations, error) {

	if xs == nil {
		xs = NewIntersections(i)
	}

	pos := r.Position(i.T)
	normal, err := i.Object.NormalAt(pos)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare computations: %w", err)
	}
	comps := &Computations{
		T:      i.T,
		Object: i.Object,
		Point:  pos,
		Eye:    r.Direction.Negate(),
		Normal: normal,
		Inside: false,
	}

	if comps.Normal.Dot(comps.Eye) < 0 {
		comps.Inside = true
		comps.Normal = comps.Normal.Negate()
	}

	comps.OverPoint = comps.Point.Add(comps.Normal.MultiplyScalar(EPSILON))
	comps.UnderPoint = comps.Point.Subtract(comps.Normal.MultiplyScalar(EPSILON))
	comps.Reflect = r.Direction.Reflect(comps.Normal)

	containers := make([]*Shape, 0)

	for _, inter := range xs.intersections {
		if inter.Equal(i) {
			if len(containers) == 0 {
				comps.N1 = 1.0
			} else {
				comps.N1 = containers[len(containers)-1].Material.RefractiveIndex
			}
		}

		if idx := slices.Index(containers, inter.Object); idx >= 0 {
			containers = slices.Delete(containers, idx, idx+1)
		} else {
			containers = append(containers, inter.Object)
		}

		if inter.Equal(i) {
			if len(containers) == 0 {
				comps.N2 = 1.0
			} else {
				comps.N2 = containers[len(containers)-1].Material.RefractiveIndex
			}
			break
		}
	}
	return comps, nil
}
