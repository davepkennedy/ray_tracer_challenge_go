package rt

import (
	"fmt"
	"math"
)

type Light struct {
	Position  *Tuple
	Intensity *Color
}

func NewPointLight(position *Tuple, intensity *Color) *Light {
	return &Light{
		Position:  position,
		Intensity: intensity,
	}
}

func (l *Light) Equal(other *Light) bool {
	return l.Position.Equal(other.Position) && l.Intensity.Equal(other.Intensity)
}

func (l *Light) String() string {
	return fmt.Sprintf("light {p: %s i: %s}", l.Position, l.Intensity)
}

func Lighting(material *Material, shape *Shape, light *Light, pt, eye, normal *Tuple, inShadow bool) (*Color, error) {
	diffuse := NewColor(0, 0, 0)
	specular := NewColor(0, 0, 0)

	var err error
	color := material.Color
	if material.Pattern != nil && shape != nil {
		color, err = material.Pattern.ColorAtObject(shape, pt)
		if err != nil {
			return nil, err
		}
	}

	effective_color := color.Multiply(light.Intensity)
	lightv := light.Position.Subtract(pt).Normalize()
	ambient := effective_color.MultiplyScalar(material.Ambient)

	if inShadow {
		return ambient, nil
	}

	lightDotNormal := lightv.Dot(normal)
	if lightDotNormal >= 0 {
		diffuse = effective_color.MultiplyScalar(material.Diffuse).MultiplyScalar(lightDotNormal)
		reflectv := lightv.Negate().Reflect(normal)
		reflectDotEye := reflectv.Dot(eye)

		if reflectDotEye > 0 {
			factor := math.Pow(reflectDotEye, material.Shininess)
			specular = light.Intensity.MultiplyScalar(material.Specular).MultiplyScalar(factor)
		}
	}
	return ambient.Add(diffuse).Add(specular), nil
}
