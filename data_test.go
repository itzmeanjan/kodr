package kodr_test

import (
	"bytes"
	"errors"
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

func TestSplitDataByCount(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	size := uint(2<<10 + rand.Intn(2<<10))
	count := uint(2<<1 + rand.Intn(int(size)))
	data := generateData(size)

	if _, err := kodr.OriginalPiecesFromDataAndPieceCount(data, 0); !(err != nil && errors.Is(err, kodr.ErrBadPieceCount)) {
		t.Fatalf("expected: %s\n", kodr.ErrBadPieceCount)
	}

	if _, err := kodr.OriginalPiecesFromDataAndPieceCount(data, size+1); !(err != nil && errors.Is(err, kodr.ErrPieceCountMoreThanTotalBytes)) {
		t.Fatalf("expected: %s\n", kodr.ErrPieceCountMoreThanTotalBytes)
	}

	pieces, err := kodr.OriginalPiecesFromDataAndPieceCount(data, count)
	if err != nil {
		t.Fatalf("didn't expect error: %s\n", err)
	}

	if len(pieces) != int(count) {
		t.Fatalf("expected %d pieces, found %d\n", count, len(pieces))
	}
}

func TestSplitDataBySize(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	size := uint(2<<10 + rand.Intn(2<<10))
	pieceSize := uint(2<<1 + rand.Intn(int(size/2)))
	data := generateData(size)

	if _, err := kodr.OriginalPiecesFromDataAndPieceSize(data, 0); !(err != nil && errors.Is(err, kodr.ErrZeroPieceSize)) {
		t.Fatalf("expected: %s\n", kodr.ErrZeroPieceSize)
	}

	if _, err := kodr.OriginalPiecesFromDataAndPieceSize(data, size); !(err != nil && errors.Is(err, kodr.ErrBadPieceCount)) {
		t.Fatalf("expected: %s\n", kodr.ErrBadPieceCount)
	}

	pieces, err := kodr.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
	if err != nil {
		t.Fatalf("didn't expect error: %s\n", err)
	}

	for i := 0; i < len(pieces); i++ {
		if len(pieces[i]) != int(pieceSize) {
			t.Fatalf("expected piece size of %d bytes; found of %d bytes", pieceSize, len(pieces[i]))
		}
	}
}

func TestCodedPieceFlattening(t *testing.T) {
	piece := &kodr.CodedPiece{Vector: generateData(2 << 5), Piece: generateData(2 << 10)}
	flat := piece.Flatten()
	if len(flat) != len(piece.Piece)+len(piece.Vector) {
		t.Fatal("coded piece flattening failed")
	}

	if !bytes.Equal(flat[:len(piece.Vector)], piece.Vector) || !bytes.Equal(flat[len(piece.Vector):], piece.Piece) {
		t.Fatal("flattened piece doesn't match << vector ++ piece >>")
	}
}

func TestCodedPiecesForRecoding(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	size := 6
	data := generateData(uint(size))
	pieceCount := 3
	enc, err := full.NewFullRLNCEncoderWithPieceCount(data, uint(pieceCount))
	if err != nil {
		t.Fatal(err.Error())
	}

	codedPieces := make([]*kodr.CodedPiece, 0, pieceCount)
	for i := 0; i < pieceCount; i++ {
		codedPieces = append(codedPieces, enc.CodedPiece())
	}

	flattenedCodedPieces := make([]byte, pieceCount*(pieceCount+size/pieceCount))
	for i := 0; i < pieceCount; i++ {
		flat := codedPieces[i].Flatten()
		if n := copy(flattenedCodedPieces[len(flat)*i:len(flat)*(i+1)], flat); n != len(flat) {
			t.Fatal("failed to copy everything !")
		}
	}

	codedPieces_, err := kodr.CodedPiecesForRecoding(flattenedCodedPieces, 3, 3)
	if err != nil {
		t.Fatal(err.Error())
	}
	for i := 0; i < len(codedPieces_); i++ {
		if !bytes.Equal(codedPieces_[i].Vector, codedPieces[i].Vector) {
			t.Fatal("coding vector mismatch !")
		}

		if !bytes.Equal(codedPieces_[i].Piece, codedPieces[i].Piece) {
			t.Fatal("coded piece mismatch !")
		}
	}
}
