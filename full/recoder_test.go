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

func recoderFlow(t *testing.T, rec *full.FullRLNCRecoder, codedPieceCount, pieceCount int, pieces []kodr.Piece) {
	recoded := make([]*kodr.CodedPiece, 0, codedPieceCount)
	for i := 0; i < codedPieceCount; i++ {
		rec_p, err := rec.CodedPiece()
		if err != nil {
			t.Fatal(err.Error())
		}
		recoded = append(recoded, rec_p)
	}

	dec := full.NewFullRLNCDecoder(uint(pieceCount))
	for i := 0; i < pieceCount; i++ {
		if _, err := dec.GetPieces(); !(err != nil && errors.Is(err, kodr.ErrMoreUsefulPiecesRequired)) {
			t.Fatal("expected error indicating more pieces are required for decoding")
		}

		if err := dec.AddPiece(recoded[i]); err != nil {
			t.Fatal(err.Error())
		}
	}

	for i := 0; i < codedPieceCount-pieceCount; i++ {
		if err := dec.AddPiece(recoded[pieceCount+i]); !(err != nil && errors.Is(err, kodr.ErrAllUsefulPiecesReceived)) {
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

func TestNewFullRLNCRecoder(t *testing.T) {
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

	rec := full.NewFullRLNCRecoder(coded)
	recoderFlow(t, rec, codedPieceCount, pieceCount, pieces)
}

func TestNewFullRLNCRecoderWithFlattenData(t *testing.T) {
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

	codedFlattened := make([]byte, 0)
	for i := 0; i < len(coded); i++ {
		codedFlattened = append(codedFlattened, coded[i].Flatten()...)
	}

	rec, err := full.NewFullRLNCRecoderWithFlattenData(codedFlattened, uint(codedPieceCount), uint(pieceCount))
	if err != nil {
		t.Fatal(err.Error())
	}

	recoderFlow(t, rec, codedPieceCount, pieceCount, pieces)
}
