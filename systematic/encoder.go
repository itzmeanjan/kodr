package systematic

import (
	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr"
)

type SystematicRLNCEncoder struct {
	currentPieceId uint
	field          *galoisfield.GF
	pieces         []kodr.Piece
}

// Generates a systematic coded piece's coding vector, which has
// only one non-zero element ( 1 )
func (s *SystematicRLNCEncoder) systematicCodingVector(idx uint, pieceCount uint) kodr.CodingVector {
	if !(idx < pieceCount) {
		return nil
	}

	vector := make(kodr.CodingVector, pieceCount)
	vector[idx] = 1
	return vector
}

// For systematic coding, first N-piece are returned in uncoded form
// i.e. coding vectors are having only single non-zero element ( 1 )
// in respective index of piece
//
// Piece index `i` ( returned from this method ), where i < N
// is going to have coding vector = [N]byte, where only i'th index
// of this vector will have 1, all other fields will have 0.
//
// Here N = #-of pieces being coded together
//
// Later pieces are coded as they're done in Full RLNC scheme
// `i` keeps incrementing by +1, until it reaches N
func (s *SystematicRLNCEncoder) CodedPiece() *kodr.CodedPiece {
	pieceCount := uint(len(s.pieces))
	if s.currentPieceId < pieceCount {
		// `nil` coding vector can be returned, which is
		// not being checked at all, as in that case we'll
		// never get into `if` branch
		vector := s.systematicCodingVector(s.currentPieceId, pieceCount)
		piece := make(kodr.Piece, len(s.pieces[s.currentPieceId]))
		copy(piece, s.pieces[s.currentPieceId])

		s.currentPieceId++
		return &kodr.CodedPiece{
			Vector: vector,
			Piece:  piece,
		}
	}

	vector := kodr.GenerateCodingVector(pieceCount)
	piece := make(kodr.Piece, len(s.pieces[0]))
	for i := range s.pieces {
		piece.Multiply(s.pieces[i], vector[i], s.field)
	}
	return &kodr.CodedPiece{
		Vector: vector,
		Piece:  piece,
	}
}
