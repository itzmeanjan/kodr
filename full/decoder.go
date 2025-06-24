package full

import (
	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/kodr_internals"
	"github.com/itzmeanjan/kodr/kodr_internals/matrix"
)

type FullRLNCDecoder struct {
	expected, useful, received uint
	state                      *matrix.DecoderState
}

// PieceLength - Returns piece length in bytes
//
// If no pieces are yet added to decoder state, then
// returns 0, denoting **unknown**
func (d *FullRLNCDecoder) PieceLength() uint {
	if d.received > 0 {
		coded := d.state.CodedPieceMatrix()
		return coded.Cols()
	}

	return 0
}

// IsDecoded - Use it for checking whether more piece
// collection is required or not
//
// If it returns false, denotes more linearly independent pieces
// need to be collected, only then decoding can be completed
func (d *FullRLNCDecoder) IsDecoded() bool {
	return d.useful >= d.expected
}

// Required - How many more linearly independent pieces
// are required for successfully decoding pieces ?
func (d *FullRLNCDecoder) Required() uint {
	return d.expected - d.useful
}

// AddPiece - Adds a new received coded piece along with
// coding vector. After every new coded piece reception
// augmented matrix ( coding vector + coded piece )
// is rref-ed, to keep it as ready as possible for consuming
// decoded pieces
//
// Note: As soon as all pieces are decoded, no more calls to
// this method does anything useful --- so better check for error & proceed !
func (d *FullRLNCDecoder) AddPiece(piece *kodr_internals.CodedPiece) error {
	// good time to start reading decoded pieces
	if d.IsDecoded() {
		return kodr.ErrAllUsefulPiecesReceived
	}

	d.state.AddPiece(piece)
	d.received++
	if !(d.received > 1) {
		d.useful++
		return nil
	}

	d.state.Rref()
	d.useful = d.state.Rank()
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
func (d *FullRLNCDecoder) GetPiece(i uint) (kodr_internals.Piece, error) {
	return d.state.GetPiece(i)
}

// GetPieces - Get a list of all decoded pieces, given full
// decoding has happened
func (d *FullRLNCDecoder) GetPieces() ([]kodr_internals.Piece, error) {
	if !d.IsDecoded() {
		return nil, kodr.ErrMoreUsefulPiecesRequired
	}

	pieces := make([]kodr_internals.Piece, 0, d.useful)
	for i := 0; i < int(d.useful); i++ {
		// error mustn't happen at this point, it should
		// have been returned fromvery first `if-block` in function
		piece, err := d.GetPiece(uint(i))
		if err != nil {
			return nil, err
		}
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
	gf := galoisfield.DefaultGF256
	state := matrix.NewDecoderStateWithPieceCount(gf, pieceCount)
	return &FullRLNCDecoder{expected: pieceCount, state: state}
}
