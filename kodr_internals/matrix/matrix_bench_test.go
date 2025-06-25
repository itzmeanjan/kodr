package matrix_test

import (
	"crypto/rand"
	"testing"

	"github.com/itzmeanjan/kodr/kodr_internals/matrix"
)

// Note: If fill_with_zero is set, it's not really a random matrix
func random_matrix(rows, cols int, fill_with_zero bool) [][]byte {
	mat := make([][]byte, 0, rows)

	for range rows {
		row := make([]byte, cols)
		// already filled with zero
		if !fill_with_zero {
			rand.Read(row)
		}
		mat = append(mat, row)
	}
	return mat
}

func BenchmarkMatrixRref(b *testing.B) {
	b.Run("2x2", func(b *testing.B) { rref(b, 1<<1) })
	b.Run("4x4", func(b *testing.B) { rref(b, 1<<2) })
	b.Run("8x8", func(b *testing.B) { rref(b, 1<<3) })
	b.Run("16x16", func(b *testing.B) { rref(b, 1<<4) })
	b.Run("32x32", func(b *testing.B) { rref(b, 1<<5) })
	b.Run("64x64", func(b *testing.B) { rref(b, 1<<6) })
	b.Run("128x128", func(b *testing.B) { rref(b, 1<<7) })
	b.Run("256x256", func(b *testing.B) { rref(b, 1<<8) })
	b.Run("512x512", func(b *testing.B) { rref(b, 1<<9) })
	b.Run("1024x1024", func(b *testing.B) { rref(b, 1<<10) })
}

func rref(b *testing.B, dim int) {
	b.SetBytes(int64(dim*dim) << 1)
	b.ReportAllocs()

	for b.Loop() {
		b.StopTimer()
		coeffs := random_matrix(dim, dim, false)
		coded := random_matrix(dim, dim, true)
		d_state := matrix.NewDecoderState(coeffs, coded)
		b.StartTimer()

		d_state.Rref()
	}
}
