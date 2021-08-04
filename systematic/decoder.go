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

func (s *SystematicRLNCDecoder) PieceLength() uint {
	return uint(len(s.pieces[0].Piece))
}

func (s *SystematicRLNCDecoder) IsDecoded() bool {
	return s.useful >= s.expected
}

func (s *SystematicRLNCDecoder) Required() uint {
	return s.expected - s.useful
}

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

func (s *SystematicRLNCDecoder) GetPiece(i uint) (kodr.Piece, error) {
	if !s.IsDecoded() || i >= s.useful {
		return nil, kodr.ErrMoreUsefulPiecesRequired
	}

	return s.rref[i][uint(len(s.rref[i]))-s.PieceLength():], nil
}

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

func NewSystematicRLNCDecoder(pieceCount uint) *SystematicRLNCDecoder {
	return &SystematicRLNCDecoder{expected: pieceCount, field: galoisfield.DefaultGF256, pieces: make([]*kodr.CodedPiece, 0)}
}
