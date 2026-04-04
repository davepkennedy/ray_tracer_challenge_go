package rt

import (
	"fmt"
	"math"
)

type Canvas struct {
	Width, Height int
	Pixels        []*Color
}

func NewCanvas(width, height int) *Canvas {
	pixels := make([]*Color, width*height)
	for i := range pixels {
		pixels[i] = NewColor(0, 0, 0) // Initialize to black
	}

	return &Canvas{
		Width:  width,
		Height: height,
		Pixels: pixels,
	}
}

func (c *Canvas) SetPixelAt(x, y int, color *Color) {
	c.Pixels[y*c.Width+x] = color
}

func (c *Canvas) PixelAt(x, y int) *Color {
	return c.Pixels[y*c.Width+x]
}

func clampColorValue(value float64) string {
	if value < 0 {
		return "0"
	}
	if value > 1 {
		return "255"
	}
	return fmt.Sprintf("%d", int(math.Ceil(value*255)))
}

func (c *Canvas) ToPPM() string {
	var ppm string
	ppm += "P3\n"
	ppm += fmt.Sprintf("%d %d\n", c.Width, c.Height)
	ppm += "255\n"

	for y := 0; y < c.Height; y++ {
		var line string
		for x := 0; x < c.Width; x++ {
			color := c.PixelAt(x, y)
			v := [3]string{clampColorValue(color.Red), clampColorValue(color.Green), clampColorValue(color.Blue)}
			for _, c := range v {
				if (len(line) + len(c)) > 69 {
					ppm += line + "\n"
					line = ""
				}
				if len(line) > 0 {
					line += " "
				}
				line += c
			}
		}
		ppm += line + "\n"
	}

	return ppm
}
