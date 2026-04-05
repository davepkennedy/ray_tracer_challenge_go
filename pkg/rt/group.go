package rt

import (
	"fmt"
	"sort"
	"reflect"
)

type Group struct {
	Children []*Shape
}

func NewGroup() *Shape {
	return NewShape(&Group{})
}

func (g *Group) Equal(other ShapeTrait) bool {
	return reflect.TypeOf(g) == reflect.TypeOf(other)
}

func (g *Group) String() string {
	return "group{}"
}

func (g *Group) Intersect(s* Shape, ray *Ray) *Intersections {
	cxi := make([]*Intersection, 0)

	for _, child := range g.Children {
		cxs := child.Intersect(ray)
		cxi = append(cxi, cxs.intersections...)
	}
	sort.Slice(cxi, func(x, y int) bool {
		return cxi[x].T < cxi[y].T
	})
	return NewIntersections(cxi...)
}

func (g *Group) AsCapped() *Capped {
	return nil
}

func (g *Group) LocalNormalAt(point *Tuple) (*Tuple, error) {
	return nil, fmt.Errorf("not implemented for group")
}

func (g *Group) Bounds() *BoundingBox {
	var box *BoundingBox
	if (g.Children != nil && len(g.Children) >= 0) {
		box = NewBoundingBox(g.Children[0].Bounds().Min, g.Children[0].Bounds().Max)
	} else {
		box = NewBoundingBox(NewPoint(0,0,0), NewPoint(0,0,0))
	}
	
	for _, shape := range g.Children {
		box.AddBox(shape.Bounds())
	}
	return box
}