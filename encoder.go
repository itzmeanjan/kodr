package kodr

import (
	"github.com/cloud9-tools/go-galoisfield"
)

type Encoder struct {
	field  *galoisfield.GF
	pieces []Piece
}

func (e *Encoder) CodedPiece() *CodedPiece {
	pieceCount := uint(len(e.pieces))
	vector := GenerateCodingVector(pieceCount)
	piece := make(Piece, len(e.pieces[0]))
	for i := range e.pieces {
		piece.Multiply(e.pieces[i], vector[i], e.field)
	}
	return &CodedPiece{
		vector: vector,
		piece:  piece,
	}
}

func NewEncoder(pieces []Piece) *Encoder {
	return &Encoder{pieces: pieces, field: galoisfield.DefaultGF256}
}
