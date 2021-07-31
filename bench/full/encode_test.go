package full_test

import (
	"crypto/rand"
	"testing"

	"github.com/itzmeanjan/kodr/full"
)

// Config: 1MB total data chunk; N Pieces, where 16 <= N <= 256
// -- starts

func BenchmarkEncoderWith1M_16(t *testing.B) {
	encode(t, 16, 1<<20)
}

func BenchmarkEncoderWith1M_32(t *testing.B) {
	encode(t, 32, 1<<20)
}

func BenchmarkEncoderWith1M_64(t *testing.B) {
	encode(t, 64, 1<<20)
}

func BenchmarkEncoderWith1M_128(t *testing.B) {
	encode(t, 128, 1<<20)
}

func BenchmarkEncoderWith1M_256(t *testing.B) {
	encode(t, 256, 1<<20)
}

// -- ends

// Config: 16MB total data chunk; N Pieces, where 16 <= N <= 256
// -- starts

func BenchmarkEncoderWith16M_16(t *testing.B) {
	encode(t, 16, 16*1<<20)
}

func BenchmarkEncoderWith16M_32(t *testing.B) {
	encode(t, 32, 16*1<<20)
}

func BenchmarkEncoderWith16M_64(t *testing.B) {
	encode(t, 64, 16*1<<20)
}

func BenchmarkEncoderWith16M_128(t *testing.B) {
	encode(t, 128, 16*1<<20)
}

func BenchmarkEncoderWith16M_256(t *testing.B) {
	encode(t, 256, 16*1<<20)
}

// -- ends

// Config: 32MB total data chunk; N Pieces, where 16 <= N <= 256
// -- starts

func BenchmarkEncoderWith32M_16(t *testing.B) {
	encode(t, 16, 32*1<<20)
}

func BenchmarkEncoderWith32M_32(t *testing.B) {
	encode(t, 32, 32*1<<20)
}

func BenchmarkEncoderWith32M_64(t *testing.B) {
	encode(t, 64, 32*1<<20)
}

func BenchmarkEncoderWith32M_128(t *testing.B) {
	encode(t, 128, 32*1<<20)
}

func BenchmarkEncoderWith32M_256(t *testing.B) {
	encode(t, 256, 32*1<<20)
}

// -- ends

// generate random data of N-bytes
func generateData(n uint) []byte {
	data := make([]byte, n)
	// can safely ignore error
	rand.Read(data)
	return data
}

func encode(t *testing.B, pieceCount uint, total uint) {
	data := generateData(total)
	enc, err := full.NewFullRLNCEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	t.ReportAllocs()
	// because pieceSize = total / pieceCount
	// so each coded piece = pieceCount + pieceSize bytes
	t.SetBytes(int64(total) + int64(pieceCount+total/pieceCount))
	t.ResetTimer()

	// keep generating encoded pieces on-the-fly
	for i := 0; i < t.N; i++ {
		enc.CodedPiece()
	}
}
