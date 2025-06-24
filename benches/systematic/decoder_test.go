package systematic_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr/full"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/systematic"
)

func BenchmarkFullRLNCDecoder(t *testing.B) {
	t.Run("1M", func(b *testing.B) {
		b.Run("16Pieces", func(b *testing.B) { decode(b, 1<<4, 1<<20) })
		b.Run("32Pieces", func(b *testing.B) { decode(b, 1<<5, 1<<20) })
		b.Run("64Pieces", func(b *testing.B) { decode(b, 1<<6, 1<<20) })
		b.Run("128Pieces", func(b *testing.B) { decode(b, 1<<7, 1<<20) })
		b.Run("256Pieces", func(b *testing.B) { decode(b, 1<<8, 1<<20) })
		b.Run("512Pieces", func(b *testing.B) { decode(b, 1<<9, 1<<20) })
	})

	t.Run("2M", func(b *testing.B) {
		b.Run("16Pieces", func(b *testing.B) { decode(b, 1<<4, 1<<21) })
		b.Run("32Pieces", func(b *testing.B) { decode(b, 1<<5, 1<<21) })
		b.Run("64Pieces", func(b *testing.B) { decode(b, 1<<6, 1<<21) })
		b.Run("128Pieces", func(b *testing.B) { decode(b, 1<<7, 1<<21) })
		b.Run("256Pieces", func(b *testing.B) { decode(b, 1<<8, 1<<21) })
		b.Run("512Pieces", func(b *testing.B) { decode(b, 1<<9, 1<<21) })
	})
}

func decode(t *testing.B, pieceCount uint, total uint) {
	rand.Seed(time.Now().UnixNano())

	data := generateData(total)
	enc, err := systematic.NewSystematicRLNCEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	pieces := make([]*kodr_internals.CodedPiece, 0, 2*pieceCount)
	for i := 0; i < int(2*pieceCount); i++ {
		pieces = append(pieces, enc.CodedPiece())
	}

	t.ResetTimer()

	totalDuration := 0 * time.Second
	for i := 0; i < t.N; i++ {
		totalDuration += decode_(t, pieceCount, pieces)
	}

	t.ReportMetric(0, "ns/op")
	t.ReportMetric(float64(totalDuration.Seconds())/float64(t.N), "second/decode")
}

func decode_(t *testing.B, pieceCount uint, pieces []*kodr_internals.CodedPiece) time.Duration {
	dec := full.NewFullRLNCDecoder(pieceCount)

	// randomly shuffle piece ordering
	rand.Shuffle(int(2*pieceCount), func(i, j int) {
		pieces[i], pieces[j] = pieces[j], pieces[i]
	})

	totalDuration := 0 * time.Second
	for j := 0; j < int(2*pieceCount); j++ {
		if j+1 >= int(pieceCount) && dec.IsDecoded() {
			break
		}

		begin := time.Now()
		dec.AddPiece(pieces[j])
		totalDuration += time.Since(begin)
	}

	if !dec.IsDecoded() {
		t.Fatal("expected pieces to be decoded")
	}

	return totalDuration
}
