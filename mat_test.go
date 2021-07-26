package kodr

import (
	"testing"

	"github.com/cloud9-tools/go-galoisfield"
)

func TestRref(t *testing.T) {
	field := galoisfield.DefaultGF256

	m_1 := Matrix{{70, 137, 2, 152}, {223, 92, 234, 98}, {217, 141, 33, 44}, {145, 135, 71, 45}}
	m_1_rref := Matrix{{1, 0, 0, 105}, {0, 1, 0, 181}, {0, 0, 1, 42}, {0, 0, 0, 0}}
	rref := m_1.Rref(field)
	if !rref.cmp(m_1_rref) {
		t.Fatalf("rref doesn't match !")
	}

	m_2 := Matrix{{68, 54, 6, 230}, {16, 56, 215, 78}, {159, 186, 146, 163}, {122, 41, 205, 133}}
	m_2_rref := Matrix{{1, 0, 0, 0}, {0, 1, 0, 0}, {0, 0, 1, 0}, {0, 0, 0, 1}}
	rref = m_2.Rref(field)
	if !rref.cmp(m_2_rref) {
		t.Fatalf("rref doesn't match !")
	}

	m_3 := Matrix{{100, 31, 76, 199, 119}, {207, 34, 207, 208, 18}, {62, 20, 54, 6, 187}, {66, 8, 52, 73, 54}, {122, 138, 247, 211, 165}}
	m_3_rref := Matrix{{1, 0, 0, 0, 0}, {0, 1, 0, 0, 0}, {0, 0, 1, 0, 0}, {0, 0, 0, 1, 0}, {0, 0, 0, 0, 1}}
	rref = m_3.Rref(field)
	if !rref.cmp(m_3_rref) {
		t.Fatalf("rref doesn't match !")
	}
}
