package test

import (
	"context"
	"fmt"
	"strconv"

	"raytracer/pkg/rt"

	"github.com/cucumber/godog"
)

func buildNByNMatrix(n int, table *godog.Table) (*rt.Matrix, error) {
	data := make([]float64, 0)

	for _, row := range table.Rows {
		for _, cell := range row.Cells {
			value, error := strconv.ParseFloat(cell.Value, 64)
			if error != nil {
				return nil, fmt.Errorf("failed to parse matrix value '%s' as float64: %v", cell.Value, error)
			}
			data = append(data, value)
		}
	}
	return rt.NewMatrix(n, data), nil
}

func createNByNMatrix(ctx context.Context, name string, n int, table *godog.Table) (context.Context, error) {
	matrix, err := buildNByNMatrix(n, table)
	if err != nil {
		return ctx, fmt.Errorf("failed to create matrix '%s': %v", name, err)
	}
	return setMatrix(ctx, name, matrix), nil
}

func InitializeMatrixScenario(sc *godog.ScenarioContext) {
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ctx = setMatrix(ctx, "identity_matrix", rt.IdentityMatrix())
		return ctx, nil
	})

	sc.Given(
		`^the following (\d+)x(\d+) matrix (\w+):$`,
		func(ctx context.Context, rows int, cols int, name string, table *godog.Table) (context.Context, error) {
			return createNByNMatrix(ctx, name, rows, table)
		})
	sc.Given(
		`^the following matrix (\w+):$`,
		func(ctx context.Context, name string, table *godog.Table) (context.Context, error) {
			return createNByNMatrix(ctx, name, 4, table)
		})

	sc.Given(
		`^(\w+) ← inverse\((\w+)\)$`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			matrix, err := getMatrix(ctx, source)
			if err != nil {
				return ctx, err
			}

			matrix, err = matrix.Inverse()
			if err != nil {
				return ctx, fmt.Errorf("matrix %s is not invertible", source)
			}
			return setMatrix(ctx, dest, matrix), nil
		})

	sc.Given(
		`^(\w+) ← transpose\((\w+)\)$`,
		func(ctx context.Context, dest, source string) (context.Context, error) {
			matrix, err := getMatrix(ctx, source)
			if err != nil {
				return ctx, err
			}

			matrix = matrix.Transpose()

			return setMatrix(ctx, dest, matrix), nil
		})
	sc.Given(
		`^(\w+) ← submatrix\((\w+), (\d+), (\d+)\)$`,
		func(ctx context.Context, dest, source string, row, col int) (context.Context, error) {
			matrix, err := getMatrix(ctx, source)
			if err != nil {
				return ctx, err
			}

			matrix = matrix.Submatrix(row, col)
			return setMatrix(ctx, dest, matrix), nil
		})
	sc.Given(
		`^(\w) ← (\w) \* (\w)$`,
		func(ctx context.Context, dest, name1, name2 string) (context.Context, error) {
			matrix1, err := getMatrix(ctx, name1)
			if err != nil {
				return ctx, fmt.Errorf("no matrix named %s found", name1)
			}
			matrix2, err := getMatrix(ctx, name2)
			if err != nil {
				return ctx, fmt.Errorf("no matrix named %s found", name2)
			}
			matrix, err := matrix1.MultiplyMatrix(matrix2)
			if err != nil {
				return ctx, fmt.Errorf("cannot multiply %s by %s", matrix1, matrix2)
			}
			return setMatrix(ctx, dest, matrix), nil
		})

	sc.When(
		`^(\w+) ← (\w) \* (\w+)$`,
		func(ctx context.Context, dest string, name1 string, name2 string) (context.Context, error) {
			matrix, err := getMatrix(ctx, name1)
			if err != nil {
				return ctx, err
			}
			tuple, err := getTuple(ctx, name2)
			if err != nil {
				return ctx, err
			}

			tuple, err = matrix.MultiplyTuple(tuple)
			if err != nil {
				return ctx, fmt.Errorf("cannot multiply %s by %s", matrix, tuple)
			}
			return setTuple(ctx, dest, tuple), nil
		})

	sc.Then(
		`^(\w)\[(\d+),(\d+)\] = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, row int, col int, expect float64) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			value := matrix.Get(row, col)
			if value != expect {
				return fmt.Errorf("expected %d,%d in %s to be %f but was %f", row, col, matrix, expect, value)
			}
			return nil
		})

	sc.Then(
		`^(\w)\[(\d+),(\d+)\] = (\-?\d+\.?\d*)/(\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, row int, col int, numer, denom float64) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			expect := numer / denom
			value := matrix.Get(row, col)
			if value != expect {
				return fmt.Errorf("expected %d,%d in %s to be %f but was %f", row, col, matrix, expect, value)
			}
			return nil
		})
	sc.Then(
		`^(\w) != (\w)$`,
		func(ctx context.Context, name1 string, name2 string) error {
			matrix1, err := getMatrix(ctx, name1)
			if err != nil {
				return fmt.Errorf("no matrix named %s found", name1)
			}
			matrix2, err := getMatrix(ctx, name2)
			if err != nil {
				return fmt.Errorf("no matrix named %s found", name2)
			}
			if matrix1.Equal(matrix2) {
				return fmt.Errorf("expected %s to not equal %s", matrix1, matrix2)
			}
			return nil
		})

	sc.Then(
		`^determinant\((\w+)\) = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, value float64) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			determinant := matrix.Determinant()
			if determinant != value {
				return fmt.Errorf("expected %f got %f", value, determinant)
			}
			return nil
		})
	sc.Then(
		`^cofactor\((\w+), (\d+), (\d+)\) = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, row int, col int, value float64) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			cofactor := matrix.Cofactor(row, col)
			if cofactor != value {
				return fmt.Errorf("expected %f got %f", value, cofactor)
			}
			return nil
		})
	sc.Then(
		`^minor\((\w+), (\d+), (\d+)\) = (\-?\d+\.?\d*)$`,
		func(ctx context.Context, name string, row int, col int, value float64) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			minor := matrix.Minor(row, col)
			if minor != value {
				return fmt.Errorf("expected %f got %f", value, minor)
			}
			return nil
		})

	sc.Then(
		`^(\w+) is invertible$`,
		func(ctx context.Context, name string) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			if !matrix.IsInvertible() {
				return fmt.Errorf("expected %s to be invertible but it is not", matrix)
			}
			return nil
		})
	sc.Then(
		`^(\w+) is not invertible$`,
		func(ctx context.Context, name string) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			if matrix.IsInvertible() {
				return fmt.Errorf("expected %s to not be invertible but it is", matrix)
			}
			return nil
		})

	sc.Then(
		`^(\w+) is the following (\d+)x(\d+) matrix:$`,
		func(ctx context.Context, name string, rows int, cols int, table *godog.Table) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			expect, err := buildNByNMatrix(rows, table)
			if err != nil {
				return fmt.Errorf("failed to build matrix from table")
			}
			if !matrix.Equal(expect) {
				return fmt.Errorf("expected %s, got %s", expect, matrix)
			}
			return nil
		})
	sc.Then(
		`^submatrix\((\w+), (\d+), (\d+)\) is the following (\d+)x(\d+) matrix:$`,
		func(ctx context.Context, name string, row int, col int, rows int, cols int, table *godog.Table) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			expect, err := buildNByNMatrix(rows, table)
			if err != nil {
				return fmt.Errorf("failed to build matrix from table")
			}
			submatrix := matrix.Submatrix(row, col)
			if !submatrix.Equal(expect) {
				return fmt.Errorf("expected %s, got %s", expect, submatrix)
			}
			return nil
		})
	sc.Then(
		`^inverse\((\w+)\) is the following (\d+)x(\d+) matrix:$`,
		func(ctx context.Context, name string, rows int, cols int, table *godog.Table) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			matrix, err = matrix.Inverse()
			if err != nil {
				return fmt.Errorf("expected %s to be invertible but it was not", matrix)
			}
			expect, err := buildNByNMatrix(rows, table)
			if err != nil {
				return fmt.Errorf("failed to build matrix from table")
			}
			if !matrix.Equal(expect) {
				return fmt.Errorf("expected %s, got %s", expect, matrix)
			}
			return nil
		})
	sc.Then(
		`^(\w) \* (\w) is the following (\d+)x(\d+) matrix:$`,
		func(ctx context.Context, name1, name2 string, rows int, cols int, table *godog.Table) error {
			matrix1, err := getMatrix(ctx, name1)
			if err != nil {
				return err
			}
			matrix2, err := getMatrix(ctx, name2)
			if err != nil {
				return err
			}

			expect, err := buildNByNMatrix(rows, table)
			if err != nil {
				return fmt.Errorf("failed to build matrix from table")
			}
			matrix, err := matrix1.MultiplyMatrix(matrix2)
			if err != nil {
				return fmt.Errorf("failed to multipled %s by %s", matrix1, matrix2)
			}
			if !matrix.Equal(expect) {
				return fmt.Errorf("multiplied %s by %s. expected %s, got %s", matrix1, matrix2, expect, matrix)
			}
			return nil
		})
	sc.Then(
		`^transpose\((\w+)\) is the following matrix:$`,
		func(ctx context.Context, name string, table *godog.Table) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			matrix = matrix.Transpose()
			expect, err := buildNByNMatrix(4, table)
			if err != nil {
				return fmt.Errorf("failed to build matrix from table")
			}
			if !matrix.Equal(expect) {
				return fmt.Errorf("expected %s, got %s", expect, matrix)
			}
			return nil
		})

	sc.Then(
		`^(\w) \* inverse\((\w)\) = (\w)$`,
		func(ctx context.Context, name1 string, name2 string, name3 string) error {
			matrix1, err := getMatrix(ctx, name1)
			if err != nil {
				return err
			}
			matrix2, err := getMatrix(ctx, name2)
			if err != nil {
				return err
			}
			matrix3, err := getMatrix(ctx, name3)
			if err != nil {
				return err
			}

			matrix2, err = matrix2.Inverse()
			if err != nil {
				return fmt.Errorf("matrix %s should be interible but it is not", matrix2)
			}
			matrix, err := matrix1.MultiplyMatrix(matrix2)
			if err != nil {
				return fmt.Errorf("%s cannot be multiplied by %s", matrix1, matrix2)
			}
			if !matrix.Equal(matrix3) {
				return fmt.Errorf("expected %s to equal %s", matrix, matrix3)
			}
			return nil
		})

	sc.Then(
		`^(\w) \* (\w) = tuple\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, matrixName, tupleName string, x, y, z, w float64) error {
			matrix, err := getMatrix(ctx, matrixName)
			if err != nil {
				return err
			}
			tuple, err := getTuple(ctx, tupleName)
			if err != nil {
				return err
			}

			tuple, err = matrix.MultiplyTuple(tuple)
			if err != nil {
				return fmt.Errorf("%s cannot be multiplied by %s", matrix, tuple)
			}
			expect := rt.NewTuple(x, y, z, w)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s got %s", expect, tuple)
			}
			return nil
		})
	sc.Then(
		`^(\w+) \* (\w) = point\((\S+), (\S+), (\S+)\)$`,
		func(ctx context.Context, matrixName, tupleName, x, y, z string) error {
			matrix, err := getMatrix(ctx, matrixName)
			if err != nil {
				return err
			}
			tuple, err := getTuple(ctx, tupleName)
			if err != nil {
				return err
			}

			tuple, err = matrix.MultiplyTuple(tuple)
			if err != nil {
				return fmt.Errorf("%s cannot be multiplied by %s", matrix, tuple)
			}
			xf, err := parseFloat(x)
			if err != nil {
				return fmt.Errorf("could not parse float from %s: %v", x, err)
			}
			yf, err := parseFloat(y)
			if err != nil {
				return fmt.Errorf("could not parse float from %s: %v", x, err)
			}
			zf, err := parseFloat(z)
			if err != nil {
				return fmt.Errorf("could not parse float from %s: %v", x, err)
			}
			expect := rt.NewPoint(xf, yf, zf)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s got %s", expect, tuple)
			}
			return nil
		})
	sc.Then(
		`^(\w+) \* (\w) = vector\((\-?\d+\.?\d*), (\-?\d+\.?\d*), (\-?\d+\.?\d*)\)$`,
		func(ctx context.Context, matrixName, tupleName string, x, y, z float64) error {
			matrix, err := getMatrix(ctx, matrixName)
			if err != nil {
				return err
			}
			tuple, err := getTuple(ctx, tupleName)
			if err != nil {
				return err
			}

			tuple, err = matrix.MultiplyTuple(tuple)
			if err != nil {
				return fmt.Errorf("%s cannot be multiplied by %s", matrix, tuple)
			}
			expect := rt.NewVector(x, y, z)
			if !tuple.Equal(expect) {
				return fmt.Errorf("expected %s got %s", expect, tuple)
			}
			return nil
		})

	sc.Then(
		`^(\w) \* identity_matrix = (\w)$`,
		func(ctx context.Context, name1, name2 string) error {
			matrix1, err := getMatrix(ctx, name1)
			if err != nil {
				return err
			}
			matrix2, err := getMatrix(ctx, name2)
			if err != nil {
				return err
			}

			matrix, err := matrix1.MultiplyMatrix(rt.IdentityMatrix())
			if err != nil {
				return fmt.Errorf("%s cannot be multiplied by identity matrix", matrix1)
			}
			if !matrix.Equal(matrix2) {
				return fmt.Errorf("expected %s * identity_matrix to produce %s, got %s", matrix1, matrix2, matrix)
			}
			return nil
		})
	sc.Then(
		`^(\w) = identity_matrix$`,
		func(ctx context.Context, name string) error {
			matrix, err := getMatrix(ctx, name)
			if err != nil {
				return err
			}

			identity := rt.IdentityMatrix()
			if !matrix.Equal(identity) {
				return fmt.Errorf("expected %s to equal identity matrix", matrix)
			}
			return nil
		})
}
