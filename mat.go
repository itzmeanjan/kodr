package kodr

import (
	"github.com/cloud9-tools/go-galoisfield"
)

type Matrix [][]byte

func (m *Matrix) copy(from Matrix) {
	*m = make([][]byte, from.rows())
	for i := range *m {
		(*m)[i] = make([]byte, from.cols())
		copy((*m)[i], from[i])
	}
}

func (m *Matrix) cmp(m_ Matrix) bool {
	if m.rows() != m_.rows() || m.cols() != m_.cols() {
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

func (m *Matrix) rows() uint {
	return uint(len(*m))
}

func (m *Matrix) cols() uint {
	return uint(len((*m)[0]))
}

func scale(row []byte, scalar byte, field *galoisfield.GF) []byte {
	nRow := make([]byte, len(row))
	for i := 0; i < len(row); i++ {
		nRow[i] = field.Mul(row[i], scalar)
	}
	return nRow
}

func invert(row []byte, field *galoisfield.GF) []byte {
	for i := 0; i < len(row); i++ {
		if row[i] != 0 {
			factor := field.Div(1, row[i])
			return scale(row, factor, field)
		}
	}
	return row
}

func pivot(row []byte) int {
	for i := 0; i < len(row); i++ {
		if row[i] == 1 {
			return i
		}
	}
	return -1
}

func add(a []byte, b []byte, field *galoisfield.GF) []byte {
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
	for i := 0; i < int(m.rows()); i++ {
		pivot_i := pivot((*m)[i])

		for j := i + 1; j < int(m.rows()); j++ {
			pivot_j := pivot((*m)[j])
			if pivot_i > pivot_j || pivot_i == -1 {
				m.swap(i, j)
				i = 0
			}
		}
	}
}

func (m *Matrix) zeroRow(row int) bool {
	yes := true
	for i := 0; i < int(m.cols()); i++ {
		if (*m)[row][i] != 0 {
			return false
		}
	}
	return yes
}

func (m *Matrix) clean() {
	for i := 0; i < int(m.rows()); i++ {
		if m.zeroRow(i) {
			(*m)[i] = nil
			copy((*m)[i:], (*m)[i+1:])
			*m = (*m)[:len(*m)-1]
			i = i - 1
		}
	}
}

func (m *Matrix) Rref(field *galoisfield.GF) Matrix {
	copied := new(Matrix)
	copied.copy(*m)

	for i := range *copied {
		row := invert((*copied)[i], field)
		copy((*copied)[i], row)
		idx := pivot((*copied)[i])
		if idx == -1 {
			continue
		}

		for j := range *copied {
			if i == j || (*copied)[j][idx] == 0 {
				continue
			}

			copy((*copied)[j], add(scale((*copied)[i], (*copied)[j][idx], field), (*copied)[j], field))
		}
	}

	copied.clean()
	copied.reorder()
	return *copied
}

// rank - Expected to be invoked on row reduced matrix
// so that ( sometimes ) rref step can be skipped
func (m *Matrix) rank() uint {
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

func (m *Matrix) Rank(field *galoisfield.GF) uint {
	rref := m.Rref(field)
	return rref.rank()
}

func (m *Matrix) Multiply(field *galoisfield.GF, with Matrix) Matrix {
	if m.cols() != with.rows() {
		return nil
	}

	mult := make([][]byte, m.rows())
	for i := 0; i < len(*m); i++ {
		mult[i] = make([]byte, with.cols())
	}

	for i := 0; i < int(m.rows()); i++ {
		for j := 0; j < int(with.cols()); j++ {

			for k := 0; k < int(m.cols()); k++ {
				mult[i][j] = field.Add(mult[i][j], field.Mul((*m)[i][k], with[k][j]))
			}

		}
	}

	return mult
}
