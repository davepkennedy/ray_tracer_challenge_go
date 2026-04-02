package test

import (
	"fmt"
	"context"
	"strings"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

type canvasKey struct {Name string}
type stringKey struct {Name string}

func InitializeCanvasScenario(sc *godog.ScenarioContext) {
	sc.Given (
		`(\w) ← canvas\((\d+), (\d+)\)`,
		func (ctx context.Context, name string, width, height int) context.Context {
			canvas := rt.NewCanvas (width, height)
			return context.WithValue (ctx, canvasKey{name}, canvas)
		})

	sc.When(
		`write_pixel\((\w+), (\d+), (\d+), (\w+)\)`,
		func (ctx context.Context, dest string, x, y int, source string) (context.Context, error) {
			canvas, ok := ctx.Value (canvasKey{dest}).(*rt.Canvas)
			if !ok {
				return ctx, fmt.Errorf("no canvas named %s found", dest)
			}
			color, ok := ctx.Value(colorKey{source}).(*rt.Color)
			if !ok {
				return ctx, fmt.Errorf("no color named %s found", source)
			}
			canvas.SetPixelAt(x, y, color)
			return ctx, nil
		})
	sc.When(
		`(\w+) ← canvas_to_ppm\((\w+)\)`,
		func (ctx context.Context, dest string, source string) (context.Context, error) {
			canvas, ok := ctx.Value (canvasKey{source}).(*rt.Canvas)
			if !ok {
				return ctx, fmt.Errorf("no canvas named %s found", source)
			}
			ppm := canvas.ToPPM()
			return context.WithValue(ctx, stringKey{dest}, ppm), nil
		})
	sc.When(
		`every pixel of (\w+) is set to color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)`,
		func (ctx context.Context, dest string, red, green, blue float64) (context.Context, error) {
			canvas, ok := ctx.Value (canvasKey{dest}).(*rt.Canvas)
			if !ok {
				return ctx, fmt.Errorf("no canvas named %s found", dest)
			}
			
			color := rt.NewColor(red, green, blue)
			for x := 0; x < canvas.Width; x++ {
				for y := 0; y < canvas.Height; y++ {
					canvas.SetPixelAt(x,y, color)
				}
			}
			return ctx, nil
		})

	sc.Then (
		`(\w).width = (\d+)`,
		func (ctx context.Context, name string, width int) error {
			canvas, ok := ctx.Value (canvasKey{name}).(*rt.Canvas)
			if !ok {
				return fmt.Errorf("no canvas named %s found", name)
			}
			if width != canvas.Width {
				return fmt.Errorf("expected width to be %d, was %d", width, canvas.Width)
			}
			return nil
		})
	sc.Then (
		`(\w).height = (\d+)`,
		func (ctx context.Context, name string, height int) error {
			canvas, ok := ctx.Value (canvasKey{name}).(*rt.Canvas)
			if !ok {
				return fmt.Errorf("no canvas named %s found", name)
			}
			if height != canvas.Height {
				return fmt.Errorf("expected width to be %d, was %d", height, canvas.Height)
			}
			return nil
		})

	sc.Then(
		`every pixel of (\w+) is color\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)`,
		func (ctx context.Context, name string, red, green, blue float64) error {
			canvas, ok := ctx.Value(canvasKey{name}).(*rt.Canvas)
			if !ok {
				return fmt.Errorf("canvas %s not found in context", name)
			}

			expectedColor := rt.NewColor(red, green, blue)
			for i, pixel := range canvas.Pixels {
				if !pixel.Equal(expectedColor) {
					return fmt.Errorf("pixel %d: expected color %s, got %s", i, expectedColor.String(), pixel.String())
				}
			}		
			return nil
		})

	sc.Then(
		`pixel_at\((\w+), (\d+), (\d+)\) = (\w+)`,
		func (ctx context.Context, dest string, x, y int, source string) error {
			canvas, ok := ctx.Value (canvasKey{dest}).(*rt.Canvas)
			if !ok {
				return fmt.Errorf("no canvas named %s found", dest)
			}
			expect, ok := ctx.Value(colorKey{source}).(*rt.Color)
			if !ok {
				return fmt.Errorf("no color named %s found", source)
			}
			color := canvas.PixelAt(x,y)
			if !color.Equal(expect) {
				return fmt.Errorf("expected color at %d, %d to be %s, was %s", x, y, expect, color)
			}
			return nil
		})

	sc.Then(
		`^lines (\d+)-(\d+) of (\w+) are$`, 
		func(ctx context.Context, startLine, endLine int, name string, docString *godog.DocString) error {
			ppm, ok := ctx.Value(stringKey{name}).(string)
			if !ok {
				return fmt.Errorf("PPM string for %s not found in context", name)
			}
			fmt.Printf("PPM output for %s:\n%s\n", name, ppm)
			lines := strings.Split(ppm, "\n")
			expectedLines := strings.Split(docString.Content, "\n")
		
			if (endLine-startLine) > len(lines) {
				return fmt.Errorf("canvas %s PPM has only %d lines, but expected up to line %d", name, len(lines), endLine)
			}

			for i, s := range expectedLines {
				actualLine := lines[startLine+i-1]
				if actualLine != s {
					return fmt.Errorf("line %d: expected '%s', got '%s'", startLine+i, s, actualLine)
				}
			} 
			return nil
		})

	sc.Then(
		`(\w+) ends with a newline character`,
		func (ctx context.Context, name string) error {
			ppm, ok := ctx.Value(stringKey{name}).(string)
			if !ok {
				return fmt.Errorf("PPM string for %s not found in context", name)
			}
			if !strings.HasSuffix(ppm, "\n") {
				return fmt.Errorf("PPM string for %s does not end with a newline character", name)
			}	
			return nil
		})
}