package matrix_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/full"
	"github.com/itzmeanjan/kodr/matrix"
)

func copyCodedPieces(c_pieces []*kodr.CodedPiece) []*kodr.CodedPiece {
	copied := make([]*kodr.CodedPiece, 0, len(c_pieces))

	for i := 0; i < len(c_pieces); i++ {
		v_len := len(c_pieces[i].Vector)
		flat := c_pieces[i].Flatten()
		_piece := kodr.CodedPiece{Vector: flat[:v_len], Piece: flat[v_len:]}
		copied = append(copied, &_piece)
	}

	return copied
}

func try_decode(b *testing.B, pieceCount, pieceLen uint64, coded_pieces []*kodr.CodedPiece) {
	ctx := context.Background()
	dec_state := matrix.NewParallelDecoderState(ctx, pieceCount, pieceLen)

	for idx := uint(0); ; idx++ {
		if err := dec_state.AddPiece(coded_pieces[idx]); errors.Is(err, kodr.ErrAllUsefulPiecesReceived) {
			break
		}
	}
}

func decoder_flow(b *testing.B, pieceCount, pieceLen, codedPieceCount uint64) {
	original_pieces := generatePieces(uint(pieceCount), uint(pieceLen))
	enc := full.NewFullRLNCEncoder(original_pieces)

	coded_pieces := make([]*kodr.CodedPiece, 0, codedPieceCount)
	for i := uint64(0); i < codedPieceCount; i++ {
		coded_pieces = append(coded_pieces, enc.CodedPiece())
	}

	b.ResetTimer()
	b.SetBytes(int64(pieceCount+pieceLen) * int64(pieceCount))
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		// don't record these allocations
		//
		// because these are not necessarily
		// required to be captured for this benchmark !
		b.StopTimer()
		copied := copyCodedPieces(coded_pieces)
		b.StartTimer()

		// record this !
		try_decode(b, pieceCount, pieceLen, copied)
	}
}

func BenchmarkParallelDecoderState(b *testing.B) {
	rand.Seed(time.Now().UnixNano())

	var (
		extraPieceCount uint64 = 1 << 4
	)

	b.Run("16Pieces", func(b *testing.B) {
		var (
			pieceCount      uint64 = 1 << 4
			codedPieceCount uint64 = pieceCount + extraPieceCount
		)

		b.Run("each_1kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<10, codedPieceCount) })
		b.Run("each_2kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<11, codedPieceCount) })
		b.Run("each_4kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<12, codedPieceCount) })
		b.Run("each_8kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<13, codedPieceCount) })
		b.Run("each_16kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<14, codedPieceCount) })
		b.Run("each_32kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<15, codedPieceCount) })
		b.Run("each_64kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<16, codedPieceCount) })
		b.Run("each_128kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<17, codedPieceCount) })
		b.Run("each_256kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<18, codedPieceCount) })
		b.Run("each_512kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<19, codedPieceCount) })
	})

	b.Run("32Pieces", func(b *testing.B) {
		var (
			pieceCount      uint64 = 1 << 5
			codedPieceCount uint64 = pieceCount + extraPieceCount
		)

		b.Run("each_1kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<10, codedPieceCount) })
		b.Run("each_2kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<11, codedPieceCount) })
		b.Run("each_4kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<12, codedPieceCount) })
		b.Run("each_8kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<13, codedPieceCount) })
		b.Run("each_16kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<14, codedPieceCount) })
		b.Run("each_32kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<15, codedPieceCount) })
		b.Run("each_64kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<16, codedPieceCount) })
		b.Run("each_128kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<17, codedPieceCount) })
		b.Run("each_256kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<18, codedPieceCount) })
		b.Run("each_512kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<19, codedPieceCount) })
	})

	b.Run("64Pieces", func(b *testing.B) {
		var (
			pieceCount      uint64 = 1 << 6
			codedPieceCount uint64 = pieceCount + extraPieceCount
		)

		b.Run("each_1kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<10, codedPieceCount) })
		b.Run("each_2kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<11, codedPieceCount) })
		b.Run("each_4kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<12, codedPieceCount) })
		b.Run("each_8kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<13, codedPieceCount) })
		b.Run("each_16kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<14, codedPieceCount) })
		b.Run("each_32kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<15, codedPieceCount) })
		b.Run("each_64kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<16, codedPieceCount) })
		b.Run("each_128kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<17, codedPieceCount) })
		b.Run("each_256kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<18, codedPieceCount) })
		b.Run("each_256kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<18, codedPieceCount) })
		b.Run("each_512kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<19, codedPieceCount) })
	})

	b.Run("128Pieces", func(b *testing.B) {
		var (
			pieceCount      uint64 = 1 << 7
			codedPieceCount uint64 = pieceCount + extraPieceCount
		)

		b.Run("each_1kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<10, codedPieceCount) })
		b.Run("each_2kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<11, codedPieceCount) })
		b.Run("each_4kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<12, codedPieceCount) })
		b.Run("each_8kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<13, codedPieceCount) })
		b.Run("each_16kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<14, codedPieceCount) })
		b.Run("each_32kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<15, codedPieceCount) })
		b.Run("each_64kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<16, codedPieceCount) })
		b.Run("each_128kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<17, codedPieceCount) })
		b.Run("each_256kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<18, codedPieceCount) })
		b.Run("each_256kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<18, codedPieceCount) })
		b.Run("each_512kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<19, codedPieceCount) })
	})

	b.Run("256Pieces", func(b *testing.B) {
		var (
			pieceCount      uint64 = 1 << 8
			codedPieceCount uint64 = pieceCount + extraPieceCount
		)

		b.Run("each_1kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<10, codedPieceCount) })
		b.Run("each_2kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<11, codedPieceCount) })
		b.Run("each_4kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<12, codedPieceCount) })
		b.Run("each_8kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<13, codedPieceCount) })
		b.Run("each_16kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<14, codedPieceCount) })
		b.Run("each_32kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<15, codedPieceCount) })
		b.Run("each_64kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<16, codedPieceCount) })
		b.Run("each_128kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<17, codedPieceCount) })
		b.Run("each_256kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<18, codedPieceCount) })
		b.Run("each_256kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<18, codedPieceCount) })
		b.Run("each_512kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<19, codedPieceCount) })
	})

	b.Run("512Pieces", func(b *testing.B) {
		var (
			pieceCount      uint64 = 1 << 9
			codedPieceCount uint64 = pieceCount + extraPieceCount
		)

		b.Run("each_1kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<10, codedPieceCount) })
		b.Run("each_2kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<11, codedPieceCount) })
		b.Run("each_4kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<12, codedPieceCount) })
		b.Run("each_8kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<13, codedPieceCount) })
		b.Run("each_16kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<14, codedPieceCount) })
		b.Run("each_32kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<15, codedPieceCount) })
		b.Run("each_64kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<16, codedPieceCount) })
		b.Run("each_128kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<17, codedPieceCount) })
		b.Run("each_256kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<18, codedPieceCount) })
		b.Run("each_256kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<18, codedPieceCount) })
		b.Run("each_512kB", func(b *testing.B) { decoder_flow(b, pieceCount, 1<<19, codedPieceCount) })
	})
}
