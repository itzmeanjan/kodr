package matrix

import (
	"github.com/cloud9-tools/go-galoisfield"
	"github.com/itzmeanjan/kodr"
)

type DecoderState struct {
	field      *galoisfield.GF
	pieceCount uint
	coeffs     Matrix
	coded      Matrix
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func (d *DecoderState) clean_forward() {
	var (
		rows     int = int(d.coeffs.Rows())
		cols     int = int(d.coeffs.Cols())
		boundary int = min(rows, cols)
	)

	for i := 0; i < boundary; i++ {
		if d.coeffs[i][i] == 0 {
			non_zero_col := false
			pivot := i + 1
			for ; pivot < rows; pivot++ {
				if d.coeffs[pivot][i] != 0 {
					non_zero_col = true
					break
				}
			}

			if !non_zero_col {
				continue
			}

			// row switching in coefficient matrix
			{
				tmp := d.coeffs[i]
				d.coeffs[i] = d.coeffs[pivot]
				d.coeffs[pivot] = tmp
			}
			// row switching in coded piece matrix
			{
				tmp := d.coded[i]
				d.coded[i] = d.coded[pivot]
				d.coded[pivot] = tmp
			}
		}

		for j := i + 1; j < rows; j++ {
			if d.coeffs[j][i] == 0 {
				continue
			}

			quotient := d.field.Div(d.coeffs[j][i], d.coeffs[i][i])
			for k := i; k < cols; k++ {
				d.coeffs[j][k] = d.field.Add(d.coeffs[j][k], d.field.Mul(d.coeffs[i][k], quotient))
			}

			for k := 0; k < len(d.coded[0]); k++ {
				d.coded[j][k] = d.field.Add(d.coded[j][k], d.field.Mul(d.coded[i][k], quotient))
			}
		}
	}
}

func (d *DecoderState) clean_backward() {
	var (
		rows     int = int(d.coeffs.Rows())
		cols     int = int(d.coeffs.Cols())
		boundary int = min(rows, cols)
	)

	for i := boundary - 1; i >= 0; i-- {
		if d.coeffs[i][i] == 0 {
			continue
		}

		for j := 0; j < i; j++ {
			if d.coeffs[j][i] == 0 {
				continue
			}

			quotient := d.field.Div(d.coeffs[j][i], d.coeffs[i][i])
			for k := i; k < cols; k++ {
				d.coeffs[j][k] = d.field.Add(d.coeffs[j][k], d.field.Mul(d.coeffs[i][k], quotient))
			}

			for k := 0; k < len(d.coded[0]); k++ {
				d.coded[j][k] = d.field.Add(d.coded[j][k], d.field.Mul(d.coded[i][k], quotient))
			}

		}

		if d.coeffs[i][i] == 1 {
			continue
		}

		inv := d.field.Div(1, d.coeffs[i][i])
		d.coeffs[i][i] = 1
		for j := i + 1; j < cols; j++ {
			if d.coeffs[i][j] == 0 {
				continue
			}

			d.coeffs[i][j] = d.field.Mul(d.coeffs[i][j], inv)
		}

		for j := 0; j < len(d.coded[0]); j++ {
			d.coded[i][j] = d.field.Mul(d.coded[i][j], inv)
		}
	}
}

func (d *DecoderState) remove_zero_rows() {
	var (
		cols = len(d.coeffs[0])
	)

	for i := 0; i < len(d.coeffs); i++ {
		yes := true
		for j := 0; j < cols; j++ {
			if d.coeffs[i][j] != 0 {
				yes = false
				break
			}
		}
		if !yes {
			continue
		}

		// resize `coeffs` matrix
		d.coeffs[i] = nil
		copy((d.coeffs)[i:], (d.coeffs)[i+1:])
		d.coeffs = (d.coeffs)[:len(d.coeffs)-1]

		// resize `coded` matrix
		d.coded[i] = nil
		copy((d.coded)[i:], (d.coded)[i+1:])
		d.coded = (d.coded)[:len(d.coded)-1]

		i = i - 1
	}
}

// Calculates Reduced Row Echelon Form of coefficient
// matrix, while also modifying coded piece matrix
// First it forward, backward cleans up matrix
// i.e. cells other than pivots are zeroed,
// later it checks if some rows of coefficient matrix
// are linearly dependent or not, if yes it removes those,
// while respective rows of coded piece matrix is also
// removed --- considered to be `not useful piece`
//
// Note: All operations are in-place, no more memory
// allocations are performed
func (d *DecoderState) Rref() {
	d.clean_forward()
	d.clean_backward()
	d.remove_zero_rows()
}

// Expected to be invoked after RREF-ed, in other words
// it won't rref matrix first to calculate rank,
// rather that needs to first invoked
func (d *DecoderState) Rank() uint {
	return d.coeffs.Rows()
}

// Current state of coding coefficient matrix
func (d *DecoderState) CoefficientMatrix() Matrix {
	return d.coeffs
}

// Current state of coded piece matrix, which is updated
// along side coding coefficient matrix ( during rref )
func (d *DecoderState) CodedPieceMatrix() Matrix {
	return d.coded
}

// Adds a new coded piece to decoder state, which will hopefully
// help in decoding pieces, if linearly independent with other rows
// i.e. read pieces
func (d *DecoderState) AddPiece(codedPiece *kodr.CodedPiece) {
	d.coeffs = append(d.coeffs, codedPiece.Vector)
	d.coded = append(d.coded, codedPiece.Piece)
}

// Request decoded piece by index ( 0 based, definitely )
//
// If piece not yet decoded/ requested index is >= #-of
// pieces coded together, returns error message indicating so
//
// Otherwise piece is returned, without any error
//
// Note: This method will copy decoded piece into newly allocated memory
// when whole decoding hasn't yet happened, to prevent any chance
// that user mistakenly modifies slice returned ( read piece )
// & that affects next round of decoding ( when new piece is received )
func (d *DecoderState) GetPiece(idx uint) (kodr.Piece, error) {
	if idx >= d.pieceCount {
		return nil, kodr.ErrPieceOutOfBound
	}
	if idx >= d.coeffs.Rows() {
		return nil, kodr.ErrPieceNotDecodedYet
	}

	if d.Rank() >= d.pieceCount {
		return d.coded[idx], nil
	}

	cols := int(d.coeffs.Cols())
	decoded := true

OUT:
	for i := 0; i < cols; i++ {
		switch i {
		case int(idx):
			if d.coeffs[idx][i] != 1 {
				decoded = false
				break OUT
			}

		default:
			if d.coeffs[idx][i] == 0 {
				decoded = false
				break OUT
			}

		}
	}

	if !decoded {
		return nil, kodr.ErrPieceNotDecodedYet
	}

	buf := make([]byte, d.coded.Cols())
	copy(buf, d.coded[idx])
	return buf, nil
}

func NewDecoderStateWithPieceCount(gf *galoisfield.GF, pieceCount uint) *DecoderState {
	coeffs := make([][]byte, 0, pieceCount)
	coded := make([][]byte, 0, pieceCount)
	return &DecoderState{field: gf, pieceCount: pieceCount, coeffs: coeffs, coded: coded}
}

func NewDecoderState(gf *galoisfield.GF, coeffs, coded Matrix) *DecoderState {
	return &DecoderState{field: gf, pieceCount: uint(len(coeffs)), coeffs: coeffs, coded: coded}
}
