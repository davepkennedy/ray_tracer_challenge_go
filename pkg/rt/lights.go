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

func (l *Light) Equal (other *Light) bool {
	return l.Position.Equal(other.Position) && l.Intensity.Equal(other.Intensity)
}

func (l *Light) String () string {
	return fmt.Sprintf("light {p: %s i: %s}", l.Position, l.Intensity)
}

func Lighting(material *Material, light *Light, point, eye, normal *Tuple, inShadow bool) (*Color, error) {

	color := material.Color
	if material.Pattern != nil {
		color = material.Pattern.ColorAt(point)
	}

	effectiveColor := color.Multiply(light.Intensity)
	lightv := light.Position.Subtract(point).Normalize()
	ambient := effectiveColor.MultiplyScalar(material.Ambient)
	
	if inShadow {
		return ambient, nil
	}

	lightDotNormal := lightv.Dot(normal)

	var diffuse *Color
	var specular *Color
	if lightDotNormal < 0 {
		diffuse = NewColor(0, 0, 0)
		specular = NewColor(0, 0, 0)
	} else {
		diffuse = effectiveColor.MultiplyScalar(material.Diffuse).MultiplyScalar(lightDotNormal)
		reflect := lightv.Negate().Reflect(normal)
		reflectDotEye := reflect.Dot(eye)

		if reflectDotEye <= 0 {
			specular = NewColor(0, 0, 0)
		} else {
			factor := math.Pow(reflectDotEye, material.Shininess)
			specular = light.Intensity.MultiplyScalar(material.Specular).MultiplyScalar(factor)
		}
	}
	return ambient.Add(diffuse).Add(specular), nil
}