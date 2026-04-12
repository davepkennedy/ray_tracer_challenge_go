package rt

import (
	"fmt"
	"math"
	"strings"
)

var (
	_ Equality = (*Matrix)(nil)
)

type Matrix struct {
	dim  int
	data []float64
}

func NewMatrix(dim int, data []float64) *Matrix {
	return &Matrix{dim: dim, data: data}
}

func IdentityMatrix() *Matrix {
	m := NewMatrix(4, make([]float64, 16))
	for i := 0; i < 4; i++ {
		m.data[i*4+i] = 1.0
	}
	return m
}

func Multiply(ma ...*Matrix) (*Matrix, error) {
	t := IdentityMatrix()
	var err error
	for _, m := range ma {
		t, err = t.MultiplyMatrix(m)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func (m *Matrix) String() string {
	lines := make([]string, 0)
	for y := 0; y < m.dim; y++ {
		row := make([]string, 0)
		for x := 0; x < m.dim; x++ {
			row = append(row, fmt.Sprintf("%f", m.data[y*m.dim+x]))
		}
		lines = append(lines, strings.Join(row, " "))
	}
	return fmt.Sprintf("matrix {%dx%d\n%s}", m.dim, m.dim, strings.Join(lines, "\n"))
}

func (m *Matrix) Get(row, col int) float64 {
	return m.data[row*m.dim+col]
}

func (m *Matrix) Set(row, col int, value float64) {
	m.data[row*m.dim+col] = value
}

func (m *Matrix) Transpose() *Matrix {
	transposedData := make([]float64, len(m.data))

	for i := range len(m.data) {
		row := i / m.dim
		col := i % m.dim
		transposedData[i] = m.data[col*m.dim+row]
	}

	return NewMatrix(m.dim, transposedData)
}

func (m *Matrix) Determinant() float64 {
	if m.dim == 2 {
		return m.data[0]*m.data[3] - m.data[1]*m.data[2]
	}

	det := 0.0
	for col := range m.dim {
		det += m.data[col] * m.Cofactor(0, col)
	}
	return det
}

func (m *Matrix) Submatrix(rowToRemove, colToRemove int) *Matrix {
	newDim := m.dim - 1
	submatrixData := make([]float64, newDim*newDim)

	for i := range len(submatrixData) {
		row := i / newDim
		col := i % newDim

		var sourceRow, sourceCol int
		if row < rowToRemove {
			sourceRow = row
		} else {
			sourceRow = row + 1
		}
		if col < colToRemove {
			sourceCol = col
		} else {
			sourceCol = col + 1
		}

		submatrixData[row*newDim+col] = m.Get(sourceRow, sourceCol)
	}
	return NewMatrix(newDim, submatrixData)
}

func (m *Matrix) Minor(row, col int) float64 {
	return m.Submatrix(row, col).Determinant()
}

func (m *Matrix) Cofactor(row, col int) float64 {
	if (row+col)%2 == 0 {
		return m.Minor(row, col)
	}
	return -m.Minor(row, col)
}

func (m *Matrix) IsInvertible() bool {
	return m.Determinant() != 0.0
}

func (m *Matrix) Inverse() (*Matrix, error) {
	det := m.Determinant()
	if det == 0.0 {
		return nil, fmt.Errorf("matrix is not invertible")
	}

	inverseData := make([]float64, len(m.data))

	for row := 0; row < m.dim; row++ {
		for col := 0; col < m.dim; col++ {
			c := m.Cofactor(row, col)
			inverseData[col*m.dim+row] = c / det
		}
	}

	return NewMatrix(m.dim, inverseData), nil
}

func (m *Matrix) MultiplyMatrix(other *Matrix) (*Matrix, error) {
	if m.dim != other.dim {
		return nil, fmt.Errorf("matrix dimensions do not match for multiplication")
	}

	resultData := make([]float64, len(m.data))

	for pos := range len(m.data) {
		row := pos / m.dim
		col := pos % m.dim

		for i := range m.dim {
			resultData[pos] += m.Get(row, i) * other.Get(i, col)
		}
	}

	return NewMatrix(m.dim, resultData), nil
}

func (m *Matrix) MultiplyTuple(other *Tuple) (*Tuple, error) {
	if m.dim != 4 {
		return nil, fmt.Errorf("Matrix must be 4x4 to multiply with Tuple")
	}

	x := m.data[0]*other.X + m.data[1]*other.Y + m.data[2]*other.Z + m.data[3]*other.W
	y := m.data[4]*other.X + m.data[5]*other.Y + m.data[6]*other.Z + m.data[7]*other.W
	z := m.data[8]*other.X + m.data[9]*other.Y + m.data[10]*other.Z + m.data[11]*other.W
	w := m.data[12]*other.X + m.data[13]*other.Y + m.data[14]*other.Z + m.data[15]*other.W

	return NewTuple(x, y, z, w), nil
}

func (m *Matrix) equalDim(other *Matrix) bool {
	return m.dim == other.dim
}

func (m *Matrix) equalData(other *Matrix) bool {
	for i := 0; i < len(m.data); i++ {
		if math.Abs(m.data[i]-other.data[i]) > 1e-5 {
			return false
		}
	}
	return true
}

func (m *Matrix) Equal(other any) bool {
	om, ok := other.(*Matrix)
	if !ok {return false}

	return m.equalDim(om) && m.equalData(om)
}
