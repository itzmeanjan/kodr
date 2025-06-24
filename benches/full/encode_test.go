package full_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr/full"
)

func BenchmarkFullRLNCEncoder(t *testing.B) {
	t.Run("1M", func(b *testing.B) {
		b.Run("16 Pieces", func(b *testing.B) { encode(b, 1<<4, 1<<20) })
		b.Run("32 Pieces", func(b *testing.B) { encode(b, 1<<5, 1<<20) })
		b.Run("64 Pieces", func(b *testing.B) { encode(b, 1<<6, 1<<20) })
		b.Run("128 Pieces", func(b *testing.B) { encode(b, 1<<7, 1<<20) })
		b.Run("256 Pieces", func(b *testing.B) { encode(b, 1<<8, 1<<20) })
	})

	t.Run("16M", func(b *testing.B) {
		b.Run("16 Pieces", func(b *testing.B) { encode(b, 1<<4, 1<<24) })
		b.Run("32 Pieces", func(b *testing.B) { encode(b, 1<<5, 1<<24) })
		b.Run("64 Pieces", func(b *testing.B) { encode(b, 1<<6, 1<<24) })
		b.Run("128 Pieces", func(b *testing.B) { encode(b, 1<<7, 1<<24) })
		b.Run("256 Pieces", func(b *testing.B) { encode(b, 1<<8, 1<<24) })
	})

	t.Run("32M", func(b *testing.B) {
		b.Run("16 Pieces", func(b *testing.B) { encode(b, 1<<4, 1<<25) })
		b.Run("32 Pieces", func(b *testing.B) { encode(b, 1<<5, 1<<25) })
		b.Run("64 Pieces", func(b *testing.B) { encode(b, 1<<6, 1<<25) })
		b.Run("128 Pieces", func(b *testing.B) { encode(b, 1<<7, 1<<25) })
		b.Run("256 Pieces", func(b *testing.B) { encode(b, 1<<8, 1<<25) })
	})
}

// generate random data of N-bytes
func generateData(n uint) []byte {
	data := make([]byte, n)
	// can safely ignore error
	rand.Read(data)
	return data
}

func encode(t *testing.B, pieceCount uint, total uint) {
	// non-reproducible random number sequence
	rand.Seed(time.Now().UnixNano())

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
