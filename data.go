package kodr

import (
	"crypto/rand"
	"math"

	"github.com/cloud9-tools/go-galoisfield"
)

// A piece of data is nothing but a byte array
type Piece []byte

// Multiple pieces are coded together by performing
// symbol by symbol finite field arithmetic, where
// a single byte is a symbol
//
// `by` is coding coefficient
func (p *Piece) Multiply(piece Piece, by byte, field *galoisfield.GF) {
	for i := range piece {
		(*p)[i] = field.Add((*p)[i], field.Mul(piece[i], by))
	}
}

// One component of coded piece; holding
// information regarding how original pieces are
// combined together
type CodingVector []byte

// Coded piece along with randomly generated coding vector
// to be used by recoder/ decoder
type CodedPiece struct {
	Vector CodingVector
	Piece  Piece
}

// Flattens coded piece into single byte
// array ( vector ++ piece ), so that
// decoding steps can be performed -- rref
// on received data matrix
func (c *CodedPiece) Flatten() []byte {
	res := make([]byte, len(c.Vector)+len(c.Piece))
	copy(res[:len(c.Vector)], c.Vector)
	copy(res[len(c.Vector):], c.Piece)
	return res
}

// Generates random coding vector of specified length
//
// No specific randomization choice is made, default available
// source is used
func GenerateCodingVector(n uint) CodingVector {
	vector := make(CodingVector, n)
	// ignoring error, because it always succeeds
	rand.Read(vector)
	return vector
}

// Given whole chunk of data & desired size of each pieces ( in terms of bytes ),
// it'll split chunk into pieces, which are to be used by encoder for performing RLNC
//
// In case whole data chunk can't be properly divided into pieces of requested size,
// extra zero bytes may be appended at end, considered as padding bytes --- given that
// each piece must be of same size
func OriginalPiecesFromDataAndPieceSize(data []byte, pieceSize uint) ([]Piece, error) {
	if pieceSize == 0 {
		return nil, ErrZeroPieceSize
	}

	if int(pieceSize) >= len(data) {
		return nil, ErrBadPieceCount
	}

	pieceCount := int(math.Ceil(float64(len(data)) / float64(pieceSize)))
	data_ := make([]byte, pieceCount*int(pieceSize))
	if n := copy(data_, data); n != len(data) {
		return nil, ErrCopyFailedDuringPieceConstruction
	}

	pieces := make([]Piece, pieceCount)
	for i := 0; i < pieceCount; i++ {
		piece := data_[int(pieceSize)*i : int(pieceSize)*(i+1)]
		pieces[i] = piece
	}

	return pieces, nil
}

// When you want to split whole data chunk into N-many original pieces, this function
// will do it, while appending extra zero bytes ( read padding bytes ) at end of last piece
// if exact division is not feasible
func OriginalPiecesFromDataAndPieceCount(data []byte, pieceCount uint) ([]Piece, error) {
	if pieceCount < 2 {
		return nil, ErrBadPieceCount
	}

	if int(pieceCount) > len(data) {
		return nil, ErrPieceCountMoreThanTotalBytes
	}

	pieceSize := uint(math.Ceil(float64(len(data)) / float64(pieceCount)))
	data_ := make([]byte, pieceSize*pieceCount)
	if n := copy(data_, data); n != len(data) {
		return nil, ErrCopyFailedDuringPieceConstruction
	}
	return OriginalPiecesFromDataAndPieceSize(data_, pieceSize)
}
