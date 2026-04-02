package rt

import (
	"fmt"
)

type Material struct {
	Color 		*Color
	Ambient		float64
	Diffuse 	float64
	Specular	float64
	Shininess	float64
	Reflective		float64
	Transparency	float64
	RefractiveIndex	float64

	Pattern			*Pattern
}

func NewMaterial() *Material {
	return &Material {
		Color: NewColor(1,1,1),
		Ambient: 0.1,
		Diffuse: 0.9,
		Specular: 0.9,
		Shininess: 200.0,
		Reflective: 0.0,
		Transparency: 0.0,
		RefractiveIndex: 1.0,
	}
}

func (m *Material) Equal (other *Material) bool {
	return m.Color.Equal(other.Color) &&
		m.Ambient			== other.Ambient &&
		m.Diffuse 			== other.Diffuse &&
		m.Specular 			== other.Specular &&
		m.Shininess 		== other.Shininess &&
		m.Reflective 		== other.Reflective &&
		m.Transparency 		== other.Transparency &&
		m.RefractiveIndex	== other.RefractiveIndex
}

func (m *Material) String () string {
	return fmt.Sprintf("material {c: %s a: %f d:%f sp: %f sh: %f r: %f t: %f ri: %f}",
		m.Color, m.Ambient, m.Diffuse, m.Specular, m.Shininess,
		m.Reflective, m.Transparency, m.RefractiveIndex)
}