package full_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/full"
)

// Config: 1MB total data chunk; N Pieces, where 16 <= N <= 256
// First code & generate N coded pieces, then recode N coded pieces
// -- starts

func BenchmarkRecoderWith1M_16Pieces(t *testing.B) {
	recode(t, 16, 1<<20)
}

func BenchmarkRecoderWith1M_32Pieces(t *testing.B) {
	recode(t, 32, 1<<20)
}

func BenchmarkRecoderWith1M_64Pieces(t *testing.B) {
	recode(t, 64, 1<<20)
}

func BenchmarkRecoderWith1M_128Pieces(t *testing.B) {
	recode(t, 128, 1<<20)
}

func BenchmarkRecoderWith1M_256Pieces(t *testing.B) {
	recode(t, 256, 1<<20)
}

// -- ends

// Config: 16MB total data chunk; N Pieces, where 16 <= N <= 256
// First code & generate N coded pieces, then recode N coded pieces

// -- starts

func BenchmarkRecoderWith16M_16Pieces(t *testing.B) {
	recode(t, 16, 16*1<<20)
}

func BenchmarkRecoderWith16M_32Pieces(t *testing.B) {
	recode(t, 32, 16*1<<20)
}

func BenchmarkRecoderWith16M_64Pieces(t *testing.B) {
	recode(t, 64, 16*1<<20)
}

func BenchmarkRecoderWith16M_128Pieces(t *testing.B) {
	recode(t, 128, 16*1<<20)
}

func BenchmarkRecoderWith16M_256Pieces(t *testing.B) {
	recode(t, 256, 16*1<<20)
}

// -- ends

// Config: 32MB total data chunk; N Pieces, where 16 <= N <= 256
// First code & generate N coded pieces, then recode N coded pieces

// -- starts

func BenchmarkRecoderWith32M_16Pieces(t *testing.B) {
	recode(t, 16, 32*1<<20)
}

func BenchmarkRecoderWith32M_32Pieces(t *testing.B) {
	recode(t, 32, 32*1<<20)
}

func BenchmarkRecoderWith32M_64Pieces(t *testing.B) {
	recode(t, 64, 32*1<<20)
}

func BenchmarkRecoderWith32M_128Pieces(t *testing.B) {
	recode(t, 128, 32*1<<20)
}

func BenchmarkRecoderWith32M_256Pieces(t *testing.B) {
	recode(t, 256, 32*1<<20)
}

// -- ends

func recode(t *testing.B, pieceCount uint, total uint) {
	// non-reproducible sequence
	rand.Seed(time.Now().UnixNano())

	// -- encode
	data := generateData(total)
	enc, err := full.NewFullRLNCEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	pieces := make([]*kodr.CodedPiece, 0, pieceCount)
	for i := 0; i < int(pieceCount); i++ {
		pieces = append(pieces, enc.CodedPiece())
	}
	// -- encoding ends

	// -- recode
	rec := full.NewFullRLNCRecoder(pieces)

	t.ReportAllocs()
	t.SetBytes(int64((pieceCount+total/pieceCount)*pieceCount) + int64(pieceCount+total/pieceCount))
	t.ResetTimer()

	for i := 0; i < t.N; i++ {
		if _, err := rec.CodedPiece(); err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}
	}
	// -- recoding ends
}
