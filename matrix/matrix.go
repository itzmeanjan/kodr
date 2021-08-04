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

func (m *Matrix) isSystematic(idx int, pieceCount uint) bool {
	row := (*m)[idx]
	c_piece := &kodr.CodedPiece{Vector: row[:pieceCount], Piece: row[pieceCount:]}
	return c_piece.IsSystematic()
}

func absSub(a, b int) int {
	c := a - b
	if c < 0 {
		return -c
	}
	return c
}

func (m *Matrix) SystematicReorder(pieceCount uint) [][]int {
	for i := 0; i < int(m.Rows()); i++ {
		pivot_i := m.pivot(i)
		if m.isSystematic(i, pieceCount) {
			// pivot row already placed in correct
			// place, don't touch
			if pivot_i == i {
				continue
			}
		}

		for j := i + 1; j < int(m.Rows()); j++ {
			pivot_j := m.pivot(j)
			if m.isSystematic(j, pieceCount) {
				// pivot row already placed in correct
				// place, don't touch
				if pivot_j == j {
					continue
				}
			}

			// both are non-pivot rows, so just ignore
			if pivot_i == -1 && pivot_j == -1 {
				continue
			}

			// does swapping take row far away from where it should
			// be placed ?
			//
			// if yes, then just ignore
			if absSub(j, pivot_j) <= absSub(i, pivot_j) {
				continue
			}

			m.swap(i, j)
			i -= 1
			break
		}
	}

	pivots := make([][]int, 0, m.Rows())
	for i := 0; i < int(m.Rows()); i++ {
		if !m.isSystematic(i, pieceCount) {
			continue
		}
		pivots = append(pivots, []int{i, m.pivot(i)})
	}
	return pivots
}

func (m *Matrix) SystematicRREF(field *galoisfield.GF, pieceCount uint) Matrix {
	copied := new(Matrix)
	copied.copy(*m)

	s_indices := copied.SystematicReorder(pieceCount)
	if len(s_indices) == 0 || len(s_indices) == int(copied.Rows()) {
		return *copied
	}

	isIn := func(v_ int) bool {
		for _, v := range s_indices {
			if v[0] == v_ {
				return true
			}
		}
		return false
	}

	for _, v := range s_indices {
		for j := v[0] + 1; j < int(copied.Rows()); j++ {
			// this row is constructed systematically, skip it
			if isIn(j) {
				continue
			}

			scaled_ := copied.scale(v[0], (*copied)[j][v[1]], field)
			added_ := add(scaled_, (*copied)[j], field)
			(*copied)[j] = nil
			(*copied)[j] = added_
		}
	}

	if !(int(copied.Rows())-len(s_indices) > 1) {
		return *copied
	}

	head := 0
	for _, v := range s_indices {
		if head != v[0] {
			break
		}

		head++
	}

	_slice_matrix := (*copied)[head:]
	_slice_row_idx := s_indices[head-1][1]
	_rows := make([][]byte, len(_slice_matrix))

	for i := 0; i < len(_slice_matrix); i++ {
		_row := _slice_matrix[i][_slice_row_idx+1:]
		_rows[i] = _row
	}

	sub_matrix := Matrix(_rows)
	_rref := sub_matrix.Rref(field)

	for i := head; i < int(copied.Rows()); i++ {
		_i := i - head
		if _i < len(_rref) {
			copy((*copied)[i][_slice_row_idx+1:], _rref[_i])
			continue
		}

		(*copied)[i] = nil
		copy((*copied)[i:], (*copied)[i+1:])
		*copied = (*copied)[:len(*copied)-1]
	}

	return *copied
}
