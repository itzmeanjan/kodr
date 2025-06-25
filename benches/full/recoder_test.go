package full_test

import (
	"testing"

	"github.com/itzmeanjan/kodr/full"
	"github.com/itzmeanjan/kodr/kodr_internals"
)

func BenchmarkFullRLNCRecoder(t *testing.B) {
	t.Run("1M", func(b *testing.B) {
		b.Run("16 Pieces", func(b *testing.B) { recode(b, 1<<4, 1<<20) })
		b.Run("32 Pieces", func(b *testing.B) { recode(b, 1<<5, 1<<20) })
		b.Run("64 Pieces", func(b *testing.B) { recode(b, 1<<6, 1<<20) })
		b.Run("128 Pieces", func(b *testing.B) { recode(b, 1<<7, 1<<20) })
		b.Run("256 Pieces", func(b *testing.B) { recode(b, 1<<8, 1<<20) })
	})

	t.Run("16M", func(b *testing.B) {
		b.Run("16 Pieces", func(b *testing.B) { recode(b, 1<<4, 1<<24) })
		b.Run("32 Pieces", func(b *testing.B) { recode(b, 1<<5, 1<<24) })
		b.Run("64 Pieces", func(b *testing.B) { recode(b, 1<<6, 1<<24) })
		b.Run("128 Pieces", func(b *testing.B) { recode(b, 1<<7, 1<<24) })
		b.Run("256 Pieces", func(b *testing.B) { recode(b, 1<<8, 1<<24) })
	})

	t.Run("32M", func(b *testing.B) {
		b.Run("16 Pieces", func(b *testing.B) { recode(b, 1<<4, 1<<25) })
		b.Run("32 Pieces", func(b *testing.B) { recode(b, 1<<5, 1<<25) })
		b.Run("64 Pieces", func(b *testing.B) { recode(b, 1<<6, 1<<25) })
		b.Run("128 Pieces", func(b *testing.B) { recode(b, 1<<7, 1<<25) })
		b.Run("256 Pieces", func(b *testing.B) { recode(b, 1<<8, 1<<25) })
	})
}

func recode(t *testing.B, pieceCount uint, total uint) {
	// Encode
	data := generateRandomData(total)
	enc, err := full.NewFullRLNCEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	pieces := make([]*kodr_internals.CodedPiece, 0, pieceCount)
	for range pieceCount {
		pieces = append(pieces, enc.CodedPiece())
	}

	// Recode
	rec := full.NewFullRLNCRecoder(pieces)

	t.ReportAllocs()
	t.SetBytes(int64((pieceCount+total/pieceCount)*pieceCount) + int64(pieceCount+total/pieceCount))
	t.ResetTimer()

	for t.Loop() {
		rec.CodedPiece()
	}
}
