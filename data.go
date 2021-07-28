package kodr

import (
	"crypto/rand"

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
