package full_test

import (
	"bytes"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/full"
)

func TestNewFullRLNCDecoder(t *testing.T) {
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
	neededPieceCount := uint(pieceCount)
	for i := 0; i < pieceCount; i++ {

		// test whether required piece count is monotonically decreasing or not
		switch i {
		case 0:
			if req_ := dec.Required(); req_ != neededPieceCount {
				t.Fatalf("expected still needed piece count to be %d, found it to be %d\n", neededPieceCount, req_)
			}
			// skip unnecessary assignment to `needPieceCount`

		default:
			if req_ := dec.Required(); !(req_ == neededPieceCount || req_ == neededPieceCount-1) {
				t.Fatal("expected required piece count monotonically decrease by 1")
			} else {
				neededPieceCount = req_
			}

		}

		// check is piece is already decoded or not --- which it should never be
		// because we're iterating over `pieceCount` #-of pieces & those many
		// we must need for decoding
		//
		// next piece to be added in follow code block
		if dec.IsDecoded() {
			t.Fatal("didn't expect it to be decoded already")
		}

		if err := dec.AddPiece(coded[i]); err != nil {
			t.Fatal(err.Error())
		}
	}

	for i := 0; i < codedPieceCount-pieceCount; i++ {
		if err := dec.AddPiece(coded[pieceCount+i]); !(err != nil && errors.Is(err, kodr.ErrAllUsefulPiecesReceived)) {
			t.Fatal("expected error indication, received nothing !")
		}
	}

	d_pieces, err := dec.GetPieces()
	if err != nil {
		t.Fatal(err.Error())
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
