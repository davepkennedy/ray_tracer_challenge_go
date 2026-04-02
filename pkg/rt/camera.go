package rt

import (
	"math"
)

type Camera struct {
	HSize			int
	VSize			int
	FieldOfView		float64
	Transform		*Matrix
	PixelSize		float64

	HalfWidth		float64
	HalfHeight		float64
}

func NewCamera (hsize, vsize int, fov float64) *Camera {
	c := &Camera{
		HSize: hsize,
		VSize: vsize,
		FieldOfView: fov,
		Transform: IdentityMatrix(),
	}

	half_view := math.Tan(fov / 2)
    aspect := float64(hsize) / float64(vsize)

    if aspect >= 1 {
            c.HalfWidth = half_view
            c.HalfHeight = half_view / aspect
	} else {
            c.HalfWidth = half_view * aspect
            c.HalfHeight = half_view
	}
    c.PixelSize = (c.HalfWidth * 2) / float64(hsize)

	return c
}

    func (c *Camera) RayForPixel(px, py int) (*Ray, error) {
        xoffset := (float64(px) + 0.5) * c.PixelSize
        yoffset := (float64(py) + 0.5) * c.PixelSize
        
        worldX := c.HalfWidth - xoffset
        worldY := c.HalfHeight - yoffset
        
		inv, err := c.Transform.Inverse()
		if err != nil {return nil, err}

        pixel, err 	:= inv.MultiplyTuple(NewPoint(worldX, worldY, -1))
		if err != nil {return nil, err}
        origin, err	:= inv.MultiplyTuple(NewPoint(0, 0, 0))
		if err != nil {return nil, err}

        direction := pixel.Subtract(origin).Normalize()
        
        return NewRay(origin, direction), nil
	}

    func (c *Camera) Render (w *World) (*Canvas, error) {
        image := NewCanvas(c.HSize, c.VSize)
    
        for y := range c.VSize {
            for x := range c.HSize {
                ray, err := c.RayForPixel(x, y)
				if err != nil {return image, err}

                color, err := w.ColorAt(ray)
				if err != nil {return image, err}

                image.SetPixelAt(x, y, color)
			}
		}
        return image, nil
	}