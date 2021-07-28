package kodr

import "errors"

var (
	ErrMatrixDimensionMismatch  = errors.New("can't perform matrix multiplication")
	ErrAllUsefulPiecesReceived  = errors.New("no more pieces required for decoding")
	ErrMoreUsefulPiecesRequired = errors.New("not enough pieces received yet to decode")
)
