package kodr

import (
	"math/rand"

	"github.com/cloud9-tools/go-galoisfield"
)

type Piece []byte

func (p *Piece) multiply(piece Piece, by byte, field *galoisfield.GF) {
	for i := range piece {
		(*p)[i] = field.Add((*p)[i], field.Mul(piece[i], by))
	}
}

type CodingVector []byte

type CodedPiece struct {
	vector CodingVector
	piece  Piece
}

type Encoder struct {
	field  *galoisfield.GF
	pieces []Piece
}

func (e *Encoder) generateCodingVector(n uint) CodingVector {
	vector := make(CodingVector, n)
	// ignoring error, because it always succeeds
	rand.Read(vector)
	return vector
}

func (e *Encoder) CodedPiece() *CodedPiece {
	pieceCount := uint(len(e.pieces))
	vector := e.generateCodingVector(pieceCount)
	piece := make(Piece, len(e.pieces[0]))
	for i := range e.pieces {
		piece.multiply(e.pieces[i], vector[i], e.field)
	}
	return &CodedPiece{
		vector: vector,
		piece:  piece,
	}
}

func NewEncoder(pieces []Piece) *Encoder {
	return &Encoder{pieces: pieces, field: galoisfield.DefaultGF256}
}
