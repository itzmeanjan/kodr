package matrix_test

import (
	"bytes"
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/full"
	"github.com/itzmeanjan/kodr/matrix"
)

// Generates `N`-bytes of random data from default
// randomization source
func generateData(n uint) []byte {
	data := make([]byte, n)
	// can safely ignore error
	rand.Read(data)
	return data
}

// Generates N-many pieces each of M-bytes length, to be used
// for testing purposes
func generatePieces(pieceCount uint, pieceLength uint) []kodr.Piece {
	pieces := make([]kodr.Piece, 0, pieceCount)
	for i := 0; i < int(pieceCount); i++ {
		pieces = append(pieces, generateData(pieceLength))
	}
	return pieces
}

func TestParallelDecoderState(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var (
		pieceCount uint64 = 1 << 8
		pieceLen   uint64 = 1 << 12
	)

	original_pieces := generatePieces(uint(pieceCount), uint(pieceLen))
	enc := full.NewFullRLNCEncoder(original_pieces)

	ctx := context.Background()
	dec_state := matrix.NewParallelDecoderState(ctx, pieceCount, pieceLen)

	start := time.Now()
	for !dec_state.IsDecoded() {
		c_piece := enc.CodedPiece()
		// simulate pieces being dropped !
		if rand.Intn(2) == 0 {
			continue
		}

		if err := dec_state.AddPiece(c_piece); err != nil && errors.Is(err, kodr.ErrAllUsefulPiecesReceived) {
			break
		}
	}

	t.Logf("decoding completed in %s\n", time.Since(start))
	t.Logf("crunched %d bytes of data\n", enc.DecodableLen())

	for i := uint64(0); i < pieceCount; i++ {
		d_piece, err := dec_state.GetPiece(i)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		if !bytes.Equal(original_pieces[i], d_piece) {
			t.Logf("decoded one doesn't match with original one for piece %d\n", i)
		}
	}
}
