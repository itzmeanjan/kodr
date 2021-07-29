package kodr_test

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr"
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
