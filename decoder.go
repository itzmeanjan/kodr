package kodr

import (
	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr/matrix"
)

type Decoder struct {
	expected uint
	useful   uint
	received uint
	field    *galoisfield.GF
	pieces   []*CodedPiece
	rref     matrix.Matrix
}

// PieceLength - Returns piece length in bytes
func (d *Decoder) PieceLength() uint {
	return uint(len(d.pieces[0].piece))
}

// AddPiece - Adds a new received coded piece along with
// coding vector. After every new coded piece reception
// augmented matrix ( coding vector + coded piece )
// is rref-ed, to keep it ready for decoding
func (d *Decoder) AddPiece(piece *CodedPiece) {
	d.pieces = append(d.pieces, piece)
	d.received++
	if !(d.received > 1) {
		return
	}
	if d.useful >= d.expected {
		return
	}

	if d.rref == nil {
		rref := make(matrix.Matrix, d.received)
		for i := range rref {
			rref[i] = d.pieces[i].flatten()
		}

		d.rref = rref
	} else {
		d.rref = append(d.rref, piece.flatten())
	}

	d.rref = d.rref.Rref(d.field)
	d.useful = d.rref.Rank_()
}

// GetPiece - Get a decoded piece by index
func (d *Decoder) GetPiece(i uint) Piece {
	if !(d.useful >= d.expected) || i >= d.useful {
		return nil
	}

	return d.rref[i][uint(len(d.rref[i]))-d.PieceLength():]
}

// GetPieces - Get a list of all decoded pieces
func (d *Decoder) GetPieces() []Piece {
	if !(d.useful >= d.expected) {
		return nil
	}

	pieces := make([]Piece, 0, d.useful)
	for i := 0; i < int(d.useful); i++ {
		pieces = append(pieces, d.GetPiece(uint(i)))
	}
	return pieces
}

func NewDecoder(pieceCount uint) *Decoder {
	return &Decoder{expected: pieceCount, field: galoisfield.DefaultGF256, pieces: make([]*CodedPiece, 0)}
}
