package matrix

import (
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals/gf256"
)

type Matrix [][]byte

// Cell by cell value comparision of two matrices, which
// returns `true` only if all cells are found to be equal
func (m *Matrix) Cmp(m_ Matrix) bool {
	if m.Rows() != m_.Rows() || m.Cols() != m_.Cols() {
		return false
	}

	for i := range *m {
		for j := range (*m)[i] {
			if (*m)[i][j] != m_[i][j] {
				return false
			}
		}
	}
	return true
}

// #-of rows in matrix
//
// This may change in runtime, when some rows are removed
// as they're found to be linearly dependent with some other
// row, after application of RREF
func (m *Matrix) Rows() uint {
	return uint(len(*m))
}

// #-of columns in matrix
//
// This isn't expected to change after initialised
func (m *Matrix) Cols() uint {
	return uint(len((*m)[0]))
}

// Multiplies two matrices ( which can be multiplied )
// in order `m x with`
func (m *Matrix) Multiply(with Matrix) (Matrix, error) {
	if m.Cols() != with.Rows() {
		return nil, kodr.ErrMatrixDimensionMismatch
	}

	mult := make([][]byte, m.Rows())
	for i := 0; i < len(*m); i++ {
		mult[i] = make([]byte, with.Cols())
	}

	for i := range m.Rows() {
		for j := range with.Cols() {
			for k := range m.Cols() {
				res := gf256.New(mult[i][j])

				l := gf256.New((*m)[i][k])
				r := gf256.New(with[k][j])

				res.AddAssign(l.Mul(r))
				mult[i][j] = res.Get()
			}
		}
	}

	return mult, nil
}
