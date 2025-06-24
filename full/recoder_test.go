package full_test

import (
	"bytes"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/full"
	"github.com/itzmeanjan/kodr/kodr_internals"
)

func recoderFlow(t *testing.T, rec *full.FullRLNCRecoder, pieceCount int, pieces []kodr_internals.Piece) {
	dec := full.NewFullRLNCDecoder(uint(pieceCount))
	for {
		r_piece, err := rec.CodedPiece()
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}
		if err := dec.AddPiece(r_piece); errors.Is(err, kodr.ErrAllUsefulPiecesReceived) {
			break
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

	coded := make([]*kodr_internals.CodedPiece, 0, codedPieceCount)
	for i := 0; i < codedPieceCount; i++ {
		coded = append(coded, enc.CodedPiece())
	}

	rec := full.NewFullRLNCRecoder(coded)
	recoderFlow(t, rec, pieceCount, pieces)
}

func TestNewFullRLNCRecoderWithFlattenData(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	pieceCount := 128
	pieceLength := 8192
	codedPieceCount := pieceCount + 2
	pieces := generatePieces(uint(pieceCount), uint(pieceLength))
	enc := full.NewFullRLNCEncoder(pieces)

	coded := make([]*kodr_internals.CodedPiece, 0, codedPieceCount)
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

	recoderFlow(t, rec, pieceCount, pieces)
}
