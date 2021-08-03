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

func isSystematicCoded(piece *kodr.CodedPiece) bool {
	pos := -1
	for i := 0; i < len(piece.Vector); i++ {
		switch piece.Vector[i] {
		case 0:
			continue

		case 1:
			if pos != -1 {
				return false
			}
			pos = i

		default:
			return false

		}
	}
	return pos >= 0 && pos < len(piece.Vector)
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
			if !isSystematicCoded(c_piece) {
				t.Fatal("expected piece to be systematic coded")
			}
		} else {
			if isSystematicCoded(c_piece) {
				t.Fatal("expected piece to be random coded")
			}
		}
	}
}
