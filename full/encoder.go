package full

import (
	"github.com/itzmeanjan/kodr/kodr_internals"
)

type FullRLNCEncoder struct {
	pieces []kodr_internals.Piece
	extra  uint
}

// Total #-of pieces being coded together --- denoting
// these many linearly independent pieces are required
// successfully decoding back to original pieces
func (f *FullRLNCEncoder) PieceCount() uint {
	return uint(len(f.pieces))
}

// Pieces which are coded together are all of same size
//
// Total data being coded = pieceSize * pieceCount ( may include
// some padding bytes )
func (f *FullRLNCEncoder) PieceSize() uint {
	return uint(len(f.pieces[0]))
}

// How many bytes of data, constructed by concatenating
// coded pieces together, required at minimum for decoding
// back to original pieces ?
//
// As I'm coding N-many pieces together, I need at least N-many
// linearly independent pieces, which are concatenated together
// to form a byte slice & can be used for original data reconstruction.
//
// So it computes N * codedPieceLen
func (f *FullRLNCEncoder) DecodableLen() uint {
	return f.PieceCount() * f.CodedPieceLen()
}

// If N-many original pieces are coded together
// what could be length of one such coded piece
// obtained by invoking `CodedPiece` ?
//
// Here N = len(pieces), original pieces which are
// being coded together
func (f *FullRLNCEncoder) CodedPieceLen() uint {
	return f.PieceCount() + f.PieceSize()
}

// How many extra padding bytes added at end of
// original data slice so that splitted pieces are
// all of same size ?
func (f *FullRLNCEncoder) Padding() uint {
	return f.extra
}

// Returns a coded piece, which is constructed on-the-fly
// by randomly drawing elements from finite field i.e.
// coding coefficients & performing full-RLNC with
// all original pieces
func (f *FullRLNCEncoder) CodedPiece() *kodr_internals.CodedPiece {
	vector := kodr_internals.GenerateCodingVector(f.PieceCount())
	piece := make(kodr_internals.Piece, f.PieceSize())
	for i := range f.pieces {
		piece.Multiply(f.pieces[i], vector[i])
	}
	return &kodr_internals.CodedPiece{
		Vector: vector,
		Piece:  piece,
	}
}

// Provide with original pieces on which fullRLNC to be performed
// & get encoder, to be used for on-the-fly generation
// to N-many coded pieces
func NewFullRLNCEncoder(pieces []kodr_internals.Piece) *FullRLNCEncoder {
	return &FullRLNCEncoder{pieces: pieces}
}

// If you know #-of pieces you want to code together, invoking
// this function splits whole data chunk into N-pieces, with padding
// bytes appended at end of last piece, if required & prepares
// full RLNC encoder for obtaining coded pieces
func NewFullRLNCEncoderWithPieceCount(data []byte, pieceCount uint) (*FullRLNCEncoder, error) {
	pieces, padding, err := kodr_internals.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
	if err != nil {
		return nil, err
	}

	enc := NewFullRLNCEncoder(pieces)
	enc.extra = padding
	return enc, nil
}

// If you want to have N-bytes piece size for each, this
// function generates M-many pieces each of N-bytes size, which are ready
// to be coded together with full RLNC
func NewFullRLNCEncoderWithPieceSize(data []byte, pieceSize uint) (*FullRLNCEncoder, error) {
	pieces, padding, err := kodr_internals.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
	if err != nil {
		return nil, err
	}

	enc := NewFullRLNCEncoder(pieces)
	enc.extra = padding
	return enc, nil
}
