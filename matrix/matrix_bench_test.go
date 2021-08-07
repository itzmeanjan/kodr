package matrix_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr/matrix"
)

func random_matrix(rows, cols int) [][]byte {
	mat := make([][]byte, 0, rows)
	for i := 0; i < rows; i++ {
		row := make([]byte, cols)
		rand.Read(row)
		mat = append(mat, row)
	}
	return mat
}

func BenchmarkMatrixRref(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	gf := galoisfield.DefaultGF256

	b.Run("2x2", func(b *testing.B) { rref(b, 1<<1, gf) })
	b.Run("4x4", func(b *testing.B) { rref(b, 1<<2, gf) })
	b.Run("8x8", func(b *testing.B) { rref(b, 1<<3, gf) })
	b.Run("16x16", func(b *testing.B) { rref(b, 1<<4, gf) })
	b.Run("32x32", func(b *testing.B) { rref(b, 1<<5, gf) })
	b.Run("64x64", func(b *testing.B) { rref(b, 1<<6, gf) })
	b.Run("128x128", func(b *testing.B) { rref(b, 1<<7, gf) })
	b.Run("256x256", func(b *testing.B) { rref(b, 1<<8, gf) })
	b.Run("512x512", func(b *testing.B) { rref(b, 1<<9, gf) })
	b.Run("1024x1024", func(b *testing.B) { rref(b, 1<<10, gf) })
}

func rref(b *testing.B, dim int, gf *galoisfield.GF) {
	b.ResetTimer()
	b.SetBytes(int64(dim * dim))
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		m := matrix.Matrix(random_matrix(dim, dim))
		m.Rref(gf)
	}
}
