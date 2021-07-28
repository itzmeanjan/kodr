package full_test

import (
	"bytes"
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/full"
)

func TestDecoder(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	pieceCount := 128
	pieceLength := 8192
	codedPieceCount := pieceCount + 2
	pieces := generatePieces(uint(pieceCount), uint(pieceLength))
	enc := full.NewFullRLNCEncoder(pieces)

	coded := make([]*kodr.CodedPiece, 0, codedPieceCount)
	for i := 0; i < codedPieceCount; i++ {
		coded = append(coded, enc.CodedPiece())
	}

	dec := full.NewFullRLNCDecoder(uint(pieceCount))
	for i := 0; i < codedPieceCount; i++ {
		dec.AddPiece(coded[i])
	}

	d_pieces := dec.GetPieces()
	if d_pieces == nil {
		t.Fatal("decoding failed !")
	}

	if len(pieces) != len(d_pieces) {
		t.Fatal("didn't decode all !")
	}

	for i := 0; i < pieceCount; i++ {
		if !bytes.Equal(pieces[i], d_pieces[i]) {
			t.Fatal("decoded data doesn't match !")
		}
	}
}
