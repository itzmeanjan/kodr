package systematic

import (
	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/matrix"
)

type SystematicRLNCDecoder struct {
	expected, useful, received uint
	field                      *galoisfield.GF
	pieces                     []*kodr.CodedPiece
	rref                       matrix.Matrix
}

// Each piece of N-many bytes
func (s *SystematicRLNCDecoder) PieceLength() uint {
	return uint(len(s.pieces[0].Piece))
}

// Already decoded back to original pieces, with collected pieces ?
//
// If yes, no more pieces need to be collected
func (s *SystematicRLNCDecoder) IsDecoded() bool {
	return s.useful >= s.expected
}

// How many more pieces are required to be collected so that
// whole data can be decoded successfully ?
//
// After collecting these many pieces, original data can be decoded
func (s *SystematicRLNCDecoder) Required() uint {
	return s.expected - s.useful
}

// Add one more collected coded piece, which will be used for decoding
// back to original pieces
//
// If all required pieces are already collected i.e. successful decoding
// has happened --- new pieces to be discarded, with an error denoting same
func (s *SystematicRLNCDecoder) AddPiece(piece *kodr.CodedPiece) error {
	s.pieces = append(s.pieces, piece)
	s.received++
	if !(s.received > 1) {
		s.useful++
		return nil
	}
	// no more piece collection is required, decoding
	// has been performed successfully
	//
	// good time to start reading decoded pieces
	if s.IsDecoded() {
		return kodr.ErrAllUsefulPiecesReceived
	}

	if s.rref == nil {
		rref := make(matrix.Matrix, s.received)
		for i := range rref {
			rref[i] = s.pieces[i].Flatten()
		}

		s.rref = rref
	} else {
		s.rref = append(s.rref, piece.Flatten())
	}

	s.rref.Rref(s.field)
	s.useful = s.rref.Rank_()
	return nil
}

// After full decoding has happened, this method can be used for finding
// i-th original piece, among N-pieces coded together --- so i < N, always
func (s *SystematicRLNCDecoder) GetPiece(i uint) (kodr.Piece, error) {
	if !s.IsDecoded() || i >= s.useful {
		return nil, kodr.ErrMoreUsefulPiecesRequired
	}

	return s.rref[i][uint(len(s.rref[i]))-s.PieceLength():], nil
}

// All original pieces in order --- only when full decoding has happened
func (s *SystematicRLNCDecoder) GetPieces() ([]kodr.Piece, error) {
	if !s.IsDecoded() {
		return nil, kodr.ErrMoreUsefulPiecesRequired
	}

	pieces := make([]kodr.Piece, 0, s.useful)
	for i := 0; i < int(s.useful); i++ {
		// safe to ignore error, because I've
		// already checked it above
		piece, _ := s.GetPiece(uint(i))
		pieces = append(pieces, piece)
	}
	return pieces, nil
}

// Pieces coded by systematic mean, alogn with randomly coded pieces,
// are decoded with this decoder
//
// @note Actually FullRLNCDecoder could have been used for same purpose
// making this one redundant
//
// I'll consider improving decoding by exploiting
// systematic coded pieces ( vectors )/ removing this
// in some future date
func NewSystematicRLNCDecoder(pieceCount uint) *SystematicRLNCDecoder {
	return &SystematicRLNCDecoder{expected: pieceCount, field: galoisfield.DefaultGF256, pieces: make([]*kodr.CodedPiece, 0)}
}
