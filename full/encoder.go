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

// If you know #-of pieces you want to code together, invoking
// this function splits whole data chunk into N-pieces, with padding
// bytes appended at end of last piece, if required & prepares
// full RLNC encoder for obtaining coded pieces
func NewFullRLNCEncoderWithPieceCount(data []byte, pieceCount uint) (*FullRLNCEncoder, error) {
	pieces, err := kodr.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
	if err != nil {
		return nil, err
	}

	return NewFullRLNCEncoder(pieces), nil
}

// If you want to have N-bytes piece size for each, this
// function generates M-many pieces each of N-bytes size, which are ready
// to be coded together with full RLNC
func NewFullRLNCEncoderWithPieceSize(data []byte, pieceSize uint) (*FullRLNCEncoder, error) {
	pieces, err := kodr.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
	if err != nil {
		return nil, err
	}

	return NewFullRLNCEncoder(pieces), nil
}
