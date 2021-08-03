package systematic_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/systematic"
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

func TestSystematicRLNCCoding(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var (
		pieceCount      uint                              = 128
		pieceLength     uint                              = 8192
		codedPieceCount uint                              = pieceCount * 2
		pieces          []kodr.Piece                      = generatePieces(pieceCount, pieceLength)
		enc             *systematic.SystematicRLNCEncoder = systematic.NewSystematicRLNCEncoder(pieces)
	)

	for i := 0; i < int(codedPieceCount); i++ {
		c_piece := enc.CodedPiece()
		if i < int(pieceCount) {
			if !c_piece.IsSystematic() {
				t.Fatal("expected piece to be systematic coded")
			}
		} else {
			if c_piece.IsSystematic() {
				t.Fatal("expected piece to be random coded")
			}
		}
	}
}
