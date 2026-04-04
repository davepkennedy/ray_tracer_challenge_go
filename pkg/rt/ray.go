package rt

type Ray struct {
	Origin    *Tuple
	Direction *Tuple
}

func NewRay(origin *Tuple, direction *Tuple) *Ray {
	return &Ray{
		Origin:    origin,
		Direction: direction,
	}
}

func (r *Ray) Position(t float64) *Tuple {
	ofs := r.Direction.MultiplyScalar(t)
	return r.Origin.Add(ofs)
}

func (r *Ray) Transform(matrix *Matrix) (*Ray, error) {
	newOrigin, error := matrix.MultiplyTuple(r.Origin)
	if error != nil {
		return nil, error
	}

	newDirection, error := matrix.MultiplyTuple(r.Direction)
	if error != nil {
		return nil, error
	}

	return NewRay(newOrigin, newDirection), nil

}
