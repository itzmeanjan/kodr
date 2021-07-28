package full

import (
	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/matrix"
)

type FullRLNCRecoder struct {
	field        *galoisfield.GF
	pieces       []*kodr.CodedPiece
	codingMatrix matrix.Matrix
}

func (r *FullRLNCRecoder) fill() {
	codingMatrix := make(matrix.Matrix, len(r.pieces))
	for i := 0; i < len(r.pieces); i++ {
		codingMatrix[i] = make([]byte, len(r.pieces[i].Vector))
		copy(codingMatrix[i], r.pieces[i].Vector)
	}
	r.codingMatrix = codingMatrix
}

// Returns recoded piece, which is constructed on-the-fly
// by randomly drawing some coding coefficients from
// finite field & performing full RLNC with all coded pieces
func (r *FullRLNCRecoder) CodedPiece() (*kodr.CodedPiece, error) {
	pieceCount := uint(len(r.pieces))
	vector := kodr.GenerateCodingVector(pieceCount)
	piece := make(kodr.Piece, len(r.pieces[0].Piece))
	for i := range r.pieces {
		piece.Multiply(r.pieces[i].Piece, vector[i], r.field)
	}

	vector_ := matrix.Matrix{vector}
	mult, err := vector_.Multiply(r.field, r.codingMatrix)
	if err != nil {
		return nil, err
	}

	return &kodr.CodedPiece{
		Vector: mult[0],
		Piece:  piece,
	}, nil
}

// Provide with all coded pieces, which are to be used
// for performing fullRLNC ( read recoding of coded data )
// & get back recoder which is used for on-the-fly construction
// of N-many recoded pieces
func NewFullRLNCRecoder(pieces []*kodr.CodedPiece) *FullRLNCRecoder {
	rec := &FullRLNCRecoder{field: galoisfield.DefaultGF256, pieces: pieces}
	rec.fill()
	return rec
}
