package rt

import (
	"math"
	"sort"
)

type World struct{
	Light *Light
	Shapes []*Shape
}

func NewWorld() *World {
	return &World{
		Light: nil,
		Shapes: make([]*Shape, 0),
	}
}

func defaultShape1() *Shape {
	s := NewSphere()
	s.Material.Color = NewColor(0.8, 1.0, 0.6)
	s.Material.Diffuse = 0.7
	s.Material.Specular = 0.2
	return s
}

func defaultShape2() *Shape {
	s := NewSphere()
	s.Transform = Scaling(0.5, 0.5, 0.5)
	return s
}

func DefaultWorld() *World {
	world := NewWorld()
	world.Light = NewPointLight(
		NewPoint(-10, 10, -10),
		NewColor(1, 1, 1))
	world.Add (defaultShape1())
	world.Add (defaultShape2())

	return world
}

func (w *World) Add (s *Shape) {
	w.Shapes = append(w.Shapes, s)
}

func (w *World) Contains (s *Shape) bool {
	for _, e := range w.Shapes {
		if e.Equal(s) {
			return true
		}
	}
	return false
}

func (w *World) Intersect (r *Ray) *Intersections {
	intersections := make ([]*Intersection, 0) 
	for _, s := range w.Shapes {
		xs := s.Intersect(r)
		for _, x := range xs.intersections {
			intersections = append(intersections, x)
		}
	}
	sort.Slice (intersections, func (x,y int) bool {
		return intersections[x].T < intersections[y].T
	})
	
	return NewIntersections(intersections...)
}

func (w *World) ShadeHit (comps *Computations) (*Color, error) {
	return w.ShadeHitWithLimit(comps, 5)
}
	
func (w *World) ShadeHitWithLimit (comps *Computations, limit int) (*Color, error) {
	surface, err := Lighting(
		comps.Object.Material,
		w.Light, 
		comps.OverPoint, comps.Eye, comps.Normal,
		w.IsShadowed(comps.OverPoint))
	if err != nil {return nil, err}

	reflected, err := w.ReflectedColorWithLimit(comps, limit-1)
	if err != nil {return nil, err}

	return surface.Add(reflected), nil
}

func (w *World) ColorAt (r *Ray) (*Color, error) {
	return w.ColorAtWithLimit(r, 5)
}

func (w *World) ColorAtWithLimit (r *Ray, limit int) (*Color, error) {
	intersections := w.Intersect(r)
	hit := intersections.Hit()
	if hit == nil {return NewColor(0,0,0), nil}

	comps, err := hit.PrepareComputations(r, intersections)
	if err != nil {return nil, err}

	return w.ShadeHitWithLimit(comps, limit)
}

func (w *World) IsShadowed (p *Tuple) bool {
	v := w.Light.Position.Subtract(p)
	distance := v.Magnitude()
	direction := v.Normalize()

	r := NewRay(p, direction)
	intersections := w.Intersect(r)

	h := intersections.Hit()
	return h != nil && h.T < distance
}

func (w *World) ReflectedColor (comps *Computations) (*Color, error) {
	return w.ReflectedColorWithLimit(comps, 5)
}

func (w *World) ReflectedColorWithLimit (comps *Computations, limit int) (*Color, error) {
	if limit <= 0 {
		return NewColor(0,0,0), nil
	}
	  
	if comps.Object.Material.Reflective == 0 {
        return NewColor(0,0,0), nil
	}
            
    reflectRay := NewRay(comps.OverPoint, comps.Reflect)
    color, err := w.ColorAtWithLimit(reflectRay, limit)
	if err != nil {return nil, err}

    return color.MultiplyScalar(comps.Object.Material.Reflective), nil
}

func (w *World) RefractedColor (comps *Computations, remaining int) (*Color, error) {
            
        if comps.Object.Material.Transparency == 0 || remaining <= 0 {
            return NewColor(0,0,0), nil
		}
            
        n_ratio := comps.N1 / comps.N2
        cosI := comps.Eye.Dot(comps.Normal)
        sin2t := (n_ratio * n_ratio) * (1 - (cosI * cosI))
        if sin2t > 1 {
            return NewColor(0,0,0), nil
		}
            
        cosT := math.Sqrt(1.0 - sin2t)
        direction := comps.Normal.MultiplyScalar(n_ratio * cosI - cosT).Subtract(comps.Eye.MultiplyScalar(n_ratio))
        refract_ray := NewRay(comps.UnderPoint, direction)
        
		color, err := w.ColorAtWithLimit(refract_ray, remaining - 1)
		if err != nil {return nil, err}
		
		return color.MultiplyScalar(comps.Object.Material.Transparency), nil
        
}