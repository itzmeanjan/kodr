package full_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/full"
	"github.com/itzmeanjan/kodr/kodr_internals"
)

func TestNewFullRLNCDecoder(t *testing.T) {
	pieceCount := 128
	pieceLength := 8192
	codedPieceCount := pieceCount + 2
	pieces := generatePieces(uint(pieceCount), uint(pieceLength))
	enc := full.NewFullRLNCEncoder(pieces)

	coded := make([]*kodr_internals.CodedPiece, 0, codedPieceCount)
	for range codedPieceCount {
		coded = append(coded, enc.CodedPiece())
	}

	dec := full.NewFullRLNCDecoder(uint(pieceCount))
	neededPieceCount := uint(pieceCount)
	for i := range codedPieceCount {

		// test whether required piece count is monotonically decreasing or not
		switch i {
		case 0:
			if req_ := dec.Required(); req_ != neededPieceCount {
				t.Fatalf("expected still needed piece count to be %d, found it to be %d\n", neededPieceCount, req_)
			}
			// skip unnecessary assignment to `needPieceCount`

		default:
			if req_ := dec.Required(); !(req_ <= neededPieceCount) {
				t.Fatal("expected required piece count to monotonically decrease")
			} else {
				neededPieceCount = req_
			}
		}

		if err := dec.AddPiece(coded[i]); errors.Is(err, kodr.ErrAllUsefulPiecesReceived) {
			break
		}
	}

	if !dec.IsDecoded() {
		t.Fatal("expected to be fully decoded !")
	}

	for i := range codedPieceCount - pieceCount {
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

	for i := range pieceCount {
		if !bytes.Equal(pieces[i], d_pieces[i]) {
			t.Fatal("decoded data doesn't match !")
		}
	}
}
