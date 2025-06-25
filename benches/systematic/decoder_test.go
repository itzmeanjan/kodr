package systematic_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/systematic"
)

func BenchmarkSystematicRLNCDecoder(t *testing.B) {
	t.Run("1M", func(b *testing.B) {
		b.Run("16Pieces", func(b *testing.B) { decode(b, 1<<4, 1<<20) })
		b.Run("32Pieces", func(b *testing.B) { decode(b, 1<<5, 1<<20) })
		b.Run("64Pieces", func(b *testing.B) { decode(b, 1<<6, 1<<20) })
		b.Run("128Pieces", func(b *testing.B) { decode(b, 1<<7, 1<<20) })
		b.Run("256Pieces", func(b *testing.B) { decode(b, 1<<8, 1<<20) })
	})

	t.Run("2M", func(b *testing.B) {
		b.Run("16Pieces", func(b *testing.B) { decode(b, 1<<4, 1<<21) })
		b.Run("32Pieces", func(b *testing.B) { decode(b, 1<<5, 1<<21) })
		b.Run("64Pieces", func(b *testing.B) { decode(b, 1<<6, 1<<21) })
		b.Run("128Pieces", func(b *testing.B) { decode(b, 1<<7, 1<<21) })
		b.Run("256Pieces", func(b *testing.B) { decode(b, 1<<8, 1<<21) })
	})

	t.Run("16M", func(b *testing.B) {
		b.Run("16 Pieces", func(b *testing.B) { decode(b, 1<<4, 1<<24) })
		b.Run("32 Pieces", func(b *testing.B) { decode(b, 1<<5, 1<<24) })
		b.Run("64 Pieces", func(b *testing.B) { decode(b, 1<<6, 1<<24) })
		b.Run("128 Pieces", func(b *testing.B) { decode(b, 1<<7, 1<<24) })
		b.Run("256 Pieces", func(b *testing.B) { decode(b, 1<<8, 1<<24) })
	})

	t.Run("32M", func(b *testing.B) {
		b.Run("16 Pieces", func(b *testing.B) { decode(b, 1<<4, 1<<25) })
		b.Run("32 Pieces", func(b *testing.B) { decode(b, 1<<5, 1<<25) })
		b.Run("64 Pieces", func(b *testing.B) { decode(b, 1<<6, 1<<25) })
		b.Run("128 Pieces", func(b *testing.B) { decode(b, 1<<7, 1<<25) })
		b.Run("256 Pieces", func(b *testing.B) { decode(b, 1<<8, 1<<25) })
	})
}

func decode(t *testing.B, pieceCount uint, total uint) {
	data := generateRandomData(total)

	enc, err := systematic.NewSystematicRLNCEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	pieces := make([]*kodr_internals.CodedPiece, 0, 2*pieceCount)
	for range 2 * pieceCount {
		pieces = append(pieces, enc.CodedPiece())
	}

	t.ResetTimer()

	totalDuration := 0 * time.Second
	for t.Loop() {
		totalDuration += decode_internal(t, pieceCount, pieces)
	}

	t.ReportMetric(0, "ns/op")
	t.ReportMetric(float64(totalDuration.Seconds())/float64(t.N), "seconds/decode")
}

func decode_internal(t *testing.B, pieceCount uint, pieces []*kodr_internals.CodedPiece) time.Duration {
	dec := systematic.NewSystematicRLNCDecoder(pieceCount)

	// Random shuffle piece ordering
	rand.Shuffle(len(pieces), func(i, j int) {
		pieces[i], pieces[j] = pieces[j], pieces[i]
	})

	totalDuration := 0 * time.Second
	for j := range 2 * pieceCount {
		if j+1 >= pieceCount && dec.IsDecoded() {
			break
		}

		begin := time.Now()
		dec.AddPiece(pieces[j])
		totalDuration += time.Since(begin)
	}

	if !dec.IsDecoded() {
		t.Fatal("expected pieces to be already decoded")
	}

	return totalDuration
}
