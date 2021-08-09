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

// Total length of coded piece --- len(coding_vector) + len(piece)
func (c *CodedPiece) Len() uint {
	return uint(len(c.Vector) + len(c.Piece))
}

// Flattens coded piece into single byte
// slice ( vector ++ piece ), so that
// decoding steps can be performed -- rref
// on received data matrix
func (c *CodedPiece) Flatten() []byte {
	res := make([]byte, c.Len())
	copy(res[:len(c.Vector)], c.Vector)
	copy(res[len(c.Vector):], c.Piece)
	return res
}

// Returns true if finds this piece is coded
// systematically i.e. piece is actually
// uncoded, just being augmented that it's coded
// which is why coding vector has only one
// non-zero element ( 1 )
func (c *CodedPiece) IsSystematic() bool {
	pos := -1
	for i := 0; i < len(c.Vector); i++ {
		switch c.Vector[i] {
		case 0:
			continue

		case 1:
			if pos != -1 {
				return false
			}
			pos = i

		default:
			return false

		}
	}
	return pos >= 0 && pos < len(c.Vector)
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
func OriginalPiecesFromDataAndPieceSize(data []byte, pieceSize uint) ([]Piece, uint, error) {
	if pieceSize == 0 {
		return nil, 0, ErrZeroPieceSize
	}

	if int(pieceSize) >= len(data) {
		return nil, 0, ErrBadPieceCount
	}

	pieceCount := int(math.Ceil(float64(len(data)) / float64(pieceSize)))
	padding := uint(pieceCount*int(pieceSize) - len(data))

	var data_ []byte
	if padding > 0 {
		data_ = make([]byte, pieceCount*int(pieceSize))
		if n := copy(data_, data); n != len(data) {
			return nil, 0, ErrCopyFailedDuringPieceConstruction
		}
	} else {
		data_ = data
	}

	pieces := make([]Piece, pieceCount)
	for i := 0; i < pieceCount; i++ {
		piece := data_[int(pieceSize)*i : int(pieceSize)*(i+1)]
		pieces[i] = piece
	}

	return pieces, padding, nil
}

// When you want to split whole data chunk into N-many original pieces, this function
// will do it, while appending extra zero bytes ( read padding bytes ) at end of last piece
// if exact division is not feasible
func OriginalPiecesFromDataAndPieceCount(data []byte, pieceCount uint) ([]Piece, uint, error) {
	if pieceCount < 2 {
		return nil, 0, ErrBadPieceCount
	}

	if int(pieceCount) > len(data) {
		return nil, 0, ErrPieceCountMoreThanTotalBytes
	}

	pieceSize := uint(math.Ceil(float64(len(data)) / float64(pieceCount)))
	padding := pieceCount*pieceSize - uint(len(data))

	var data_ []byte
	if padding > 0 {
		data_ = make([]byte, pieceSize*pieceCount)
		if n := copy(data_, data); n != len(data) {
			return nil, 0, ErrCopyFailedDuringPieceConstruction
		}
	} else {
		data_ = data
	}

	// padding field being ignored, because I've already computed it
	// in line 134
	//
	// Here ignored field will always be 0, because it's already extended ( if required ) to be
	// properly divisible by `pieceSize`, which is checked in function invoked below
	splitted, _, err := OriginalPiecesFromDataAndPieceSize(data_, pieceSize)
	return splitted, padding, err
}

// Before recoding can be performed, coded pieces byte array i.e. []<< coding vector ++ coded piece >>
// where each coded piece is << coding vector ++ coded piece >> ( flattened ) is splitted into
// structured data i.e. into components {coding vector, coded piece}, where how many coded pieces are
// present in byte array ( read `data` ) & how many pieces are coded together ( read coding vector length )
// are provided
func CodedPiecesForRecoding(data []byte, pieceCount uint, piecesCodedTogether uint) ([]*CodedPiece, error) {
	codedPieceLength := len(data) / int(pieceCount)
	if codedPieceLength*int(pieceCount) != len(data) {
		return nil, ErrCodedDataLengthMismatch
	}

	if !(piecesCodedTogether < uint(codedPieceLength)) {
		return nil, ErrCodingVectorLengthMismatch
	}

	codedPieces := make([]*CodedPiece, pieceCount)
	for i := 0; i < int(pieceCount); i++ {
		codedPiece := data[codedPieceLength*i : codedPieceLength*(i+1)]
		codedPieces[i] = &CodedPiece{
			Vector: codedPiece[:piecesCodedTogether],
			Piece:  codedPiece[piecesCodedTogether:],
		}
	}

	return codedPieces, nil
}
