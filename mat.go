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

	return *copied
}

func (m *Matrix) Rank(field *galoisfield.GF) uint {
	rref := m.Rref(field)
	var count uint
	for i := range rref {
		for j := range rref {
			if rref[i][j] == 1 {
				count += 1
				break
			}
		}
	}
	return count
}
