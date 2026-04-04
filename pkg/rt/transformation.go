package rt

import "math"

func Translation(x float64, y float64, z float64) *Matrix {
	t := IdentityMatrix()
	t.Set(0, 3, x)
	t.Set(1, 3, y)
	t.Set(2, 3, z)
	return t
}

func Scaling(x float64, y float64, z float64) *Matrix {
	s := IdentityMatrix()
	s.Set(0, 0, x)
	s.Set(1, 1, y)
	s.Set(2, 2, z)
	return s
}

func RotationX(radians float64) *Matrix {
	r := IdentityMatrix()
	r.Set(1, 1, math.Cos(radians))
	r.Set(1, 2, -math.Sin(radians))
	r.Set(2, 1, math.Sin(radians))
	r.Set(2, 2, math.Cos(radians))
	return r
}

func RotationY(radians float64) *Matrix {
	r := IdentityMatrix()
	r.Set(0, 0, math.Cos(radians))
	r.Set(0, 2, math.Sin(radians))
	r.Set(2, 0, -math.Sin(radians))
	r.Set(2, 2, math.Cos(radians))
	return r
}

func RotationZ(radians float64) *Matrix {
	r := IdentityMatrix()
	r.Set(0, 0, math.Cos(radians))
	r.Set(0, 1, -math.Sin(radians))
	r.Set(1, 0, math.Sin(radians))
	r.Set(1, 1, math.Cos(radians))
	return r
}

func Shearing(xy float64, xz float64, yx float64, yz float64, zx float64, zy float64) *Matrix {
	s := IdentityMatrix()
	s.Set(0, 1, xy)
	s.Set(0, 2, xz)
	s.Set(1, 0, yx)
	s.Set(1, 2, yz)
	s.Set(2, 0, zx)
	s.Set(2, 1, zy)
	return s
}

func ViewTransformation(from, to, up *Tuple) (*Matrix, error) {
	forward := to.Subtract(from).Normalize()
	left := forward.Cross(up.Normalize())
	trueUp := left.Cross(forward)

	orientation := NewMatrix(4, []float64{
		left.X, left.Y, left.Z, 0,
		trueUp.X, trueUp.Y, trueUp.Z, 0,
		-forward.X, -forward.Y, -forward.Z, 0,
		0, 0, 0, 1,
	})
	return orientation.MultiplyMatrix(Translation(-from.X, -from.Y, -from.Z))
}
