package kodr

import (
	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr/matrix"
)

type Recoder struct {
	field        *galoisfield.GF
	pieces       []*CodedPiece
	codingMatrix matrix.Matrix
}

func (r *Recoder) fill() {
	codingMatrix := make(matrix.Matrix, len(r.pieces))
	for i := 0; i < len(r.pieces); i++ {
		codingMatrix[i] = make([]byte, len(r.pieces[i].vector))
		copy(codingMatrix[i], r.pieces[i].vector)
	}
	r.codingMatrix = codingMatrix
}

func (r *Recoder) CodedPiece() *CodedPiece {
	pieceCount := uint(len(r.pieces))
	vector := generateCodingVector(pieceCount)
	piece := make(Piece, len(r.pieces[0].piece))
	for i := range r.pieces {
		piece.multiply(r.pieces[i].piece, vector[i], r.field)
	}

	vector_ := matrix.Matrix{vector}
	mult := vector_.Multiply(r.field, r.codingMatrix)
	if mult == nil {
		return nil
	}

	return &CodedPiece{
		vector: mult[0],
		piece:  piece,
	}
}

func NewRecoder(pieces []*CodedPiece) *Recoder {
	rec := &Recoder{field: galoisfield.DefaultGF256, pieces: pieces}
	rec.fill()
	return rec
}
