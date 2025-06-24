package systematic

import (
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/matrix"
)

type SystematicRLNCDecoder struct {
	expected, useful, received uint
	state                      *matrix.DecoderState
}

// Each piece of N-many bytes
//
// Note: If no pieces are yet added to decoder state, then
// returns 0, denoting **unknown**
func (s *SystematicRLNCDecoder) PieceLength() uint {
	if s.received > 0 {
		coded := s.state.CodedPieceMatrix()
		return coded.Cols()
	}

	return 0
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
func (s *SystematicRLNCDecoder) AddPiece(piece *kodr_internals.CodedPiece) error {
	if s.IsDecoded() {
		return kodr.ErrAllUsefulPiecesReceived
	}

	s.state.AddPiece(piece)
	s.received++
	if !(s.received > 1) {
		s.useful++
		return nil
	}

	s.state.Rref()
	s.useful = s.state.Rank()
	return nil
}

// GetPiece - Get a decoded piece by index, may ( not ) succeed !
//
// Note: It's not necessary that full decoding needs to happen
// for this method to return something useful
//
// If M-many pieces are received among N-many expected ( read M <= N )
// then pieces with index in [0..M] ( remember upper bound exclusive )
// can be attempted to be consumed, given algebric structure has revealed
// requested piece at index `i`
func (s *SystematicRLNCDecoder) GetPiece(i uint) (kodr_internals.Piece, error) {
	return s.state.GetPiece(i)
}

// All original pieces in order --- only when full decoding has happened
func (s *SystematicRLNCDecoder) GetPieces() ([]kodr_internals.Piece, error) {
	if !s.IsDecoded() {
		return nil, kodr.ErrMoreUsefulPiecesRequired
	}

	pieces := make([]kodr_internals.Piece, 0, s.useful)
	for i := 0; i < int(s.useful); i++ {
		piece, err := s.GetPiece(uint(i))
		if err != nil {
			return nil, err
		}
		pieces = append(pieces, piece)
	}
	return pieces, nil
}

// Pieces coded by systematic mean, along with randomly coded pieces,
// are decoded with this decoder
//
// @note Actually FullRLNCDecoder could have been used for same purpose
// making this one redundant
//
// I'll consider improving decoding by exploiting
// systematic coded pieces ( vectors )/ removing this
// in some future date
func NewSystematicRLNCDecoder(pieceCount uint) *SystematicRLNCDecoder {
	state := matrix.NewDecoderStateWithPieceCount(pieceCount)
	return &SystematicRLNCDecoder{expected: pieceCount, state: state}
}
