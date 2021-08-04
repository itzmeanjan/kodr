package matrix_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/matrix"
)

func TestMatrixRref(t *testing.T) {
	field := galoisfield.DefaultGF256

	m_1 := matrix.Matrix{{70, 137, 2, 152}, {223, 92, 234, 98}, {217, 141, 33, 44}, {145, 135, 71, 45}}
	m_1_rref := matrix.Matrix{{1, 0, 0, 105}, {0, 1, 0, 181}, {0, 0, 1, 42}}
	m_1.Rref(field)
	if !m_1.Cmp(m_1_rref) {
		t.Fatal("rref doesn't match !")
	}

	m_2 := matrix.Matrix{{68, 54, 6, 230}, {16, 56, 215, 78}, {159, 186, 146, 163}, {122, 41, 205, 133}}
	m_2_rref := matrix.Matrix{{1, 0, 0, 0}, {0, 1, 0, 0}, {0, 0, 1, 0}, {0, 0, 0, 1}}
	m_2.Rref(field)
	if !m_2.Cmp(m_2_rref) {
		t.Fatal("rref doesn't match !")
	}

	m_3 := matrix.Matrix{{100, 31, 76, 199, 119}, {207, 34, 207, 208, 18}, {62, 20, 54, 6, 187}, {66, 8, 52, 73, 54}, {122, 138, 247, 211, 165}}
	m_3_rref := matrix.Matrix{{1, 0, 0, 0, 0}, {0, 1, 0, 0, 0}, {0, 0, 1, 0, 0}, {0, 0, 0, 1, 0}, {0, 0, 0, 0, 1}}
	m_3.Rref(field)
	if !m_3.Cmp(m_3_rref) {
		t.Fatal("rref doesn't match !")
	}
}

func TestMatrixRank(t *testing.T) {
	field := galoisfield.DefaultGF256

	m_1 := matrix.Matrix{{70, 137, 2, 152}, {223, 92, 234, 98}, {217, 141, 33, 44}, {145, 135, 71, 45}}
	if rank := m_1.Rank(field); rank != 3 {
		t.Fatalf("expected rank 3, received %d", rank)
	}

	m_2 := matrix.Matrix{{68, 54, 6, 230}, {16, 56, 215, 78}, {159, 186, 146, 163}, {122, 41, 205, 133}}
	if rank := m_2.Rank(field); rank != 4 {
		t.Fatalf("expected rank 4, received %d", rank)
	}

	m_3 := matrix.Matrix{{100, 31, 76, 199, 119}, {207, 34, 207, 208, 18}, {62, 20, 54, 6, 187}, {66, 8, 52, 73, 54}, {122, 138, 247, 211, 165}}
	if rank := m_3.Rank(field); rank != 5 {
		t.Fatalf("expected rank 5, received %d", rank)
	}
}

func TestMatrixMultiplication(t *testing.T) {
	field := galoisfield.DefaultGF256

	m_1 := matrix.Matrix{{102, 82, 165, 0}}
	m_2 := matrix.Matrix{{157, 233, 247}, {160, 28, 233}, {149, 234, 117}, {200, 181, 55}}
	m_3 := matrix.Matrix{{1, 2, 3}}
	expected := matrix.Matrix{{186, 23, 11}}

	if _, err := m_3.Multiply(field, m_2); !(err != nil && errors.Is(err, kodr.ErrMatrixDimensionMismatch)) {
		t.Fatal("expected failed matrix multiplication error indication")
	}

	mult, err := m_1.Multiply(field, m_2)
	if err != nil {
		t.Fatal(err.Error())
	}

	for i := 0; i < int(expected.Rows()); i++ {
		if !bytes.Equal(expected[i], mult[i]) {
			t.Fatal("row mismatch !")
		}
	}
}
