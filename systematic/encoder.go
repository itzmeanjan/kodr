package systematic

import (
	"github.com/itzmeanjan/kodr/kodr_internals"
)

type SystematicRLNCEncoder struct {
	currentPieceId uint
	pieces         []kodr_internals.Piece
	extra          uint
}

// Total #-of pieces being coded together --- denoting
// these many linearly independent pieces are required
// successfully decoding back to original pieces
func (s *SystematicRLNCEncoder) PieceCount() uint {
	return uint(len(s.pieces))
}

// Pieces which are coded together are all of same size
//
// Total data being coded = pieceSize * pieceCount ( may include
// some padding bytes )
func (s *SystematicRLNCEncoder) PieceSize() uint {
	return uint(len(s.pieces[0]))
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
func (s *SystematicRLNCEncoder) DecodableLen() uint {
	return s.PieceCount() * s.CodedPieceLen()
}

// If N-many original pieces are coded together
// what could be length of one such coded piece
// obtained by invoking `CodedPiece` ?
//
// Here N = len(pieces), original pieces which are
// being coded together
func (s *SystematicRLNCEncoder) CodedPieceLen() uint {
	return s.PieceCount() + s.PieceSize()
}

// If any extra padding bytes added at end of original
// data slice for making all pieces of same size,
// returned value will be >0
func (s *SystematicRLNCEncoder) Padding() uint {
	return s.extra
}

// Generates a systematic coded piece's coding vector, which has
// only one non-zero element ( 1 )
func (s *SystematicRLNCEncoder) systematicCodingVector(idx uint) kodr_internals.CodingVector {
	if !(idx < s.PieceCount()) {
		return nil
	}

	vector := make(kodr_internals.CodingVector, s.PieceCount())
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
func (s *SystematicRLNCEncoder) CodedPiece() *kodr_internals.CodedPiece {
	if s.currentPieceId < s.PieceCount() {
		// `nil` coding vector can be returned, which is
		// not being checked at all, as in that case we'll
		// never get into `if` branch
		vector := s.systematicCodingVector(s.currentPieceId)
		piece := make(kodr_internals.Piece, s.PieceSize())
		copy(piece, s.pieces[s.currentPieceId])

		s.currentPieceId++
		return &kodr_internals.CodedPiece{
			Vector: vector,
			Piece:  piece,
		}
	}

	vector := kodr_internals.GenerateCodingVector(s.PieceCount())
	piece := make(kodr_internals.Piece, s.PieceSize())

	for i := range s.pieces {
		piece.Multiply(s.pieces[i], vector[i])
	}

	return &kodr_internals.CodedPiece{
		Vector: vector,
		Piece:  piece,
	}
}

// When you've already splitted original data chunk into pieces
// of same length ( in terms of bytes ), this function can be used
// for creating one systematic RLNC encoder, which delivers coded pieces
// on-the-fly
func NewSystematicRLNCEncoder(pieces []kodr_internals.Piece) *SystematicRLNCEncoder {
	return &SystematicRLNCEncoder{currentPieceId: 0, pieces: pieces}
}

// If you know #-of pieces you want to code together, invoking
// this function splits whole data chunk into N-pieces, with padding
// bytes appended at end of last piece, if required & prepares
// full RLNC encoder for obtaining coded pieces
func NewSystematicRLNCEncoderWithPieceCount(data []byte, pieceCount uint) (*SystematicRLNCEncoder, error) {
	pieces, padding, err := kodr_internals.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
	if err != nil {
		return nil, err
	}

	enc := NewSystematicRLNCEncoder(pieces)
	enc.extra = padding
	return enc, nil
}

// If you want to have N-bytes piece size for each, this
// function generates M-many pieces each of N-bytes size, which are ready
// to be coded together with full RLNC
func NewSystematicRLNCEncoderWithPieceSize(data []byte, pieceSize uint) (*SystematicRLNCEncoder, error) {
	pieces, padding, err := kodr_internals.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
	if err != nil {
		return nil, err
	}

	enc := NewSystematicRLNCEncoder(pieces)
	enc.extra = padding
	return enc, nil
}
