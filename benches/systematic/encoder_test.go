package systematic_test

import (
	"crypto/rand"
	"testing"

	"github.com/itzmeanjan/kodr/systematic"
)

func BenchmarkSystematicRLNCEncoder(t *testing.B) {
	t.Run("1M", func(b *testing.B) {
		b.Run("16Pieces", func(b *testing.B) { encode(b, 1<<4, 1<<20) })
		b.Run("32Pieces", func(b *testing.B) { encode(b, 1<<5, 1<<20) })
		b.Run("64Pieces", func(b *testing.B) { encode(b, 1<<6, 1<<20) })
		b.Run("128Pieces", func(b *testing.B) { encode(b, 1<<7, 1<<20) })
		b.Run("256Pieces", func(b *testing.B) { encode(b, 1<<8, 1<<20) })
	})

	t.Run("2M", func(b *testing.B) {
		b.Run("16Pieces", func(b *testing.B) { encode(b, 1<<4, 1<<21) })
		b.Run("32Pieces", func(b *testing.B) { encode(b, 1<<5, 1<<21) })
		b.Run("64Pieces", func(b *testing.B) { encode(b, 1<<6, 1<<21) })
		b.Run("128Pieces", func(b *testing.B) { encode(b, 1<<7, 1<<21) })
		b.Run("256Pieces", func(b *testing.B) { encode(b, 1<<8, 1<<21) })
	})
}

func generateRandomData(n uint) []byte {
	data := make([]byte, n)
	rand.Read(data)

	return data
}

func encode(t *testing.B, pieceCount uint, total uint) {
	data := generateRandomData(total)

	enc, err := systematic.NewSystematicRLNCEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	t.ReportAllocs()
	t.SetBytes(int64(total+enc.Padding()) + int64(enc.CodedPieceLen()))
	t.ResetTimer()

	for t.Loop() {
		enc.CodedPiece()
	}
}
