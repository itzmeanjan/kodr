package matrix

import (
	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr"
)

type Matrix [][]byte

func (m *Matrix) copy(from Matrix) {
	*m = make([][]byte, from.Rows())
	for i := range *m {
		(*m)[i] = make([]byte, from.Cols())
		copy((*m)[i], from[i])
	}
}

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

func (m *Matrix) Rows() uint {
	return uint(len(*m))
}

func (m *Matrix) Cols() uint {
	return uint(len((*m)[0]))
}

func (m *Matrix) scale(idx int, scalar byte, field *galoisfield.GF) []byte {
	nRow := make([]byte, len((*m)[idx]))
	for i := 0; i < len((*m)[idx]); i++ {
		nRow[i] = field.Mul((*m)[idx][i], scalar)
	}
	return nRow
}

func (m *Matrix) invert(idx int, field *galoisfield.GF) []byte {
	for i := 0; i < len((*m)[idx]); i++ {
		if (*m)[idx][i] != 0 {
			factor := field.Div(1, (*m)[idx][i])
			return m.scale(idx, factor, field)
		}
	}
	return (*m)[idx]
}

func (m *Matrix) pivot(idx int) int {
	for i := 0; i < len((*m)[idx]); i++ {
		if (*m)[idx][i] == 1 {
			return i
		}
	}
	return -1
}

func add(a, b []byte, field *galoisfield.GF) []byte {
	c := make([]byte, len(a))
	for i := range a {
		c[i] = field.Add(a[i], b[i])
	}
	return c
}

func (m *Matrix) swap(i, j int) {
	tmp := make([]byte, len((*m)[i]))
	copy(tmp, (*m)[i])
	(*m)[i] = (*m)[j]
	(*m)[j] = tmp
}

func (m *Matrix) reorder() {
	for i := 0; i < int(m.Rows()); i++ {
		pivot_i := m.pivot(i)

		for j := i + 1; j < int(m.Rows()); j++ {
			pivot_j := m.pivot(j)
			if pivot_i > pivot_j || pivot_i == -1 {
				m.swap(i, j)
				i = 0
			}
		}
	}
}

func (m *Matrix) zeroRow(row int) bool {
	yes := true
	for i := 0; i < int(m.Cols()); i++ {
		if (*m)[row][i] != 0 {
			return false
		}
	}
	return yes
}

func (m *Matrix) clean() {
	for i := 0; i < int(m.Rows()); i++ {
		if m.zeroRow(i) {
			(*m)[i] = nil
			copy((*m)[i:], (*m)[i+1:])
			*m = (*m)[:len(*m)-1]
			i = i - 1
		}
	}
}

// Rref - Get matrix into reduced row echelon form, where
// matrix elements are GF(2**8) element, which are good fit
// for representing in 1 byte
func (m *Matrix) Rref(field *galoisfield.GF) Matrix {
	copied := new(Matrix)
	copied.copy(*m)

	for i := range *copied {
		row := copied.invert(i, field)
		copy((*copied)[i], row)
		idx := copied.pivot(i)
		if idx == -1 {
			continue
		}

		for j := range *copied {
			if i == j || (*copied)[j][idx] == 0 {
				continue
			}

			copy((*copied)[j], add(copied.scale(i, (*copied)[j][idx], field), (*copied)[j], field))
		}
	}

	copied.clean()
	copied.reorder()
	return *copied
}

// Rank_ - Expected to be invoked on row reduced matrix
// so that rref step can be skipped
//
// If you've a matrix which is not yet rref-ed, you may want
// to invoke `Rank()`
func (m *Matrix) Rank_() uint {
	var count uint
	for i := range *m {
		for j := range *m {
			if (*m)[i][j] == 1 {
				count += 1
				break
			}
		}
	}
	return count
}

// Rank - Make use of this method when you've a
// matrix which is not yet rref-ed
func (m *Matrix) Rank(field *galoisfield.GF) uint {
	rref := m.Rref(field)
	return rref.Rank_()
}

// Multiply - Multiplies two matrices ( which can be multiplied )
// in order `m x with`
func (m *Matrix) Multiply(field *galoisfield.GF, with Matrix) (Matrix, error) {
	if m.Cols() != with.Rows() {
		return nil, kodr.ErrMatrixDimensionMismatch
	}

	mult := make([][]byte, m.Rows())
	for i := 0; i < len(*m); i++ {
		mult[i] = make([]byte, with.Cols())
	}

	for i := 0; i < int(m.Rows()); i++ {
		for j := 0; j < int(with.Cols()); j++ {

			for k := 0; k < int(m.Cols()); k++ {
				mult[i][j] = field.Add(mult[i][j], field.Mul((*m)[i][k], with[k][j]))
			}

		}
	}

	return mult, nil
}
