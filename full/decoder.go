package full

import (
	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/matrix"
)

type FullRLNCDecoder struct {
	expected uint
	useful   uint
	received uint
	field    *galoisfield.GF
	pieces   []*kodr.CodedPiece
	rref     matrix.Matrix
}

// PieceLength - Returns piece length in bytes
func (d *FullRLNCDecoder) PieceLength() uint {
	return uint(len(d.pieces[0].Piece))
}

// AddPiece - Adds a new received coded piece along with
// coding vector. After every new coded piece reception
// augmented matrix ( coding vector + coded piece )
// is rref-ed, to keep it as ready as possible for consuming
// decoded pieces
func (d *FullRLNCDecoder) AddPiece(piece *kodr.CodedPiece) error {
	d.pieces = append(d.pieces, piece)
	d.received++
	if !(d.received > 1) {
		return nil
	}
	// no more piece collection is required, decoding
	// has been performed successfully
	//
	// good time to start reading decoded pieces
	if d.useful >= d.expected {
		return kodr.ErrAllUsefulPiecesReceived
	}

	if d.rref == nil {
		rref := make(matrix.Matrix, d.received)
		for i := range rref {
			rref[i] = d.pieces[i].Flatten()
		}

		d.rref = rref
	} else {
		d.rref = append(d.rref, piece.Flatten())
	}

	d.rref = d.rref.Rref(d.field)
	d.useful = d.rref.Rank_()
	return nil
}

// GetPiece - Get a decoded piece by index, given full
// decoding has happened
func (d *FullRLNCDecoder) GetPiece(i uint) (kodr.Piece, error) {
	if !(d.useful >= d.expected) || i >= d.useful {
		return nil, kodr.ErrMoreUsefulPiecesRequired
	}

	return d.rref[i][uint(len(d.rref[i]))-d.PieceLength():], nil
}

// GetPieces - Get a list of all decoded pieces, given full
// decoding has happened
func (d *FullRLNCDecoder) GetPieces() ([]kodr.Piece, error) {
	if !(d.useful >= d.expected) {
		return nil, kodr.ErrMoreUsefulPiecesRequired
	}

	pieces := make([]kodr.Piece, 0, d.useful)
	for i := 0; i < int(d.useful); i++ {
		// safe to ignore error, because I've
		// already checked it above
		piece, _ := d.GetPiece(uint(i))
		pieces = append(pieces, piece)
	}
	return pieces, nil
}

// If minimum #-of linearly independent coded pieces required
// for decoding coded pieces --- is provided with,
// it returns a decoder, which keeps applying
// full RLNC decoding step on received coded pieces
//
// As soon as minimum #-of linearly independent pieces are obtained
// which is generally equal to original #-of pieces, decoded pieces
// can be read back
func NewFullRLNCDecoder(pieceCount uint) *FullRLNCDecoder {
	return &FullRLNCDecoder{expected: pieceCount, field: galoisfield.DefaultGF256, pieces: make([]*kodr.CodedPiece, 0)}
}
