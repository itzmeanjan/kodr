package full_test

import (
	"bytes"
	"errors"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/full"
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

func encoderFlow(t *testing.T, enc *full.FullRLNCEncoder, pieceCount, codedPieceCount int, pieces []kodr.Piece) {
	coded := make([]*kodr.CodedPiece, 0, codedPieceCount)
	for i := 0; i < codedPieceCount; i++ {
		coded = append(coded, enc.CodedPiece())
	}

	dec := full.NewFullRLNCDecoder(uint(pieceCount))
	for i := 0; i < pieceCount; i++ {
		if _, err := dec.GetPieces(); !(err != nil && errors.Is(err, kodr.ErrMoreUsefulPiecesRequired)) {
			t.Fatal("expected error indicating more pieces are required for decoding")
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

func TestNewFullRLNCEncoder(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	pieceCount := 128
	pieceLength := 8192
	codedPieceCount := pieceCount + 2
	pieces := generatePieces(uint(pieceCount), uint(pieceLength))
	enc := full.NewFullRLNCEncoder(pieces)

	encoderFlow(t, enc, pieceCount, codedPieceCount, pieces)
}

func TestNewFullRLNCEncoderWithPieceCount(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	size := uint(2<<10 + rand.Intn(2<<10))
	pieceCount := uint(2<<1 + rand.Intn(2<<8))
	codedPieceCount := pieceCount + 2
	data := generateData(size)
	t.Logf("\nTotal Data: %d bytes\nPiece Count: %d\nCoded Piece Count: %d\n", size, pieceCount, codedPieceCount)

	pieces, err := kodr.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
	if err != nil {
		t.Fatal(err.Error())
	}

	enc, err := full.NewFullRLNCEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		t.Fatal(err.Error())
	}

	encoderFlow(t, enc, int(pieceCount), int(codedPieceCount), pieces)
}

func TestNewFullRLNCEncoderWithPieceSize(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	size := uint(2<<10 + rand.Intn(2<<10))
	pieceSize := uint(2<<5 + rand.Intn(2<<5))
	pieceCount := int(math.Ceil(float64(size) / float64(pieceSize)))
	codedPieceCount := pieceCount + 2
	data := generateData(size)
	t.Logf("\nTotal Data: %d bytes\nPiece Size: %d bytes\nPiece Count: %d\nCoded Piece Count: %d\n", size, pieceSize, pieceCount, codedPieceCount)

	pieces, err := kodr.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
	if err != nil {
		t.Fatal(err.Error())
	}

	enc, err := full.NewFullRLNCEncoderWithPieceSize(data, pieceSize)
	if err != nil {
		t.Fatal(err.Error())
	}

	encoderFlow(t, enc, pieceCount, codedPieceCount, pieces)
}
