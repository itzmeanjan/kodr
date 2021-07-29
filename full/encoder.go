package full

import (
	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr"
)

type FullRLNCEncoder struct {
	field  *galoisfield.GF
	pieces []kodr.Piece
}

// Returns a coded piece, which is constructed on-the-fly
// by randomly drawing elements from finite field i.e.
// coding coefficients & performing full-RLNC with
// all original pieces
func (e *FullRLNCEncoder) CodedPiece() *kodr.CodedPiece {
	pieceCount := uint(len(e.pieces))
	vector := kodr.GenerateCodingVector(pieceCount)
	piece := make(kodr.Piece, len(e.pieces[0]))
	for i := range e.pieces {
		piece.Multiply(e.pieces[i], vector[i], e.field)
	}
	return &kodr.CodedPiece{
		Vector: vector,
		Piece:  piece,
	}
}

// Provide with original pieces on which fullRLNC to be performed
// & get encoder, to be used for on-the-fly generation
// to N-many coded pieces
func NewFullRLNCEncoder(pieces []kodr.Piece) *FullRLNCEncoder {
	return &FullRLNCEncoder{pieces: pieces, field: galoisfield.DefaultGF256}
}

func NewFullRLNCEncoderWithPieceCount(data []byte, pieceCount uint) *FullRLNCEncoder {
	return NewFullRLNCEncoder(nil)
}

func NewFullRLNCEncoderWithPieceSize(data []byte, pieceSize uint) *FullRLNCEncoder {
	return NewFullRLNCEncoder(nil)
}
