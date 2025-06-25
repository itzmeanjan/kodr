package systematic_test

import (
	"bytes"
	"errors"
	"math/rand"
	"testing"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/systematic"
)

func TestNewSystematicRLNCDecoder(t *testing.T) {
	var (
		pieceCount  uint                              = 128
		pieceLength uint                              = 8192
		pieces      []kodr_internals.Piece            = generatePieces(pieceCount, pieceLength)
		enc         *systematic.SystematicRLNCEncoder = systematic.NewSystematicRLNCEncoder(pieces)
		dec         *systematic.SystematicRLNCDecoder = systematic.NewSystematicRLNCDecoder(pieceCount)
	)

	for {
		c_piece := enc.CodedPiece()

		// simulate random coded_piece drop/ loss
		if rand.Intn(2) == 0 {
			continue
		}

		err := dec.AddPiece(c_piece)
		if errors.Is(err, kodr.ErrAllUsefulPiecesReceived) {
			if v := dec.Required(); v != 0 {
				t.Fatalf("required piece count should be: %d\n", v)
			}
			break
		}
	}

	d_pieces, err := dec.GetPieces()
	if err != nil {
		t.Fatalf("Error: %s\n", err.Error())
	}

	if len(d_pieces) != len(pieces) {
		t.Fatal("didn't decode all !")
	}

	for i := 0; i < int(pieceCount); i++ {
		if !bytes.Equal(pieces[i], d_pieces[i]) {
			t.Fatal("decoded data doesn't match !")
		}
	}
}
