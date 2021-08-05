package full

import (
	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr"
)

type FullRLNCEncoder struct {
	field  *galoisfield.GF
	pieces []kodr.Piece
	extra  uint
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
	return uint(len(f.pieces)) * f.CodedPieceLen()
}

// If N-many original pieces are coded together
// what could be length of one such coded piece
// obtained by invoking `CodedPiece` ?
//
// Here N = len(pieces), original pieces which are
// being coded together
func (f *FullRLNCEncoder) CodedPieceLen() uint {
	return uint(len(f.pieces) + len(f.pieces[0]))
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
	pieces, padding, err := kodr.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
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
	pieces, padding, err := kodr.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
	if err != nil {
		return nil, err
	}

	enc := NewFullRLNCEncoder(pieces)
	enc.extra = padding
	return enc, nil
}
