package kodr

import "errors"

var (
	ErrMatrixDimensionMismatch           = errors.New("can't perform matrix multiplication")
	ErrAllUsefulPiecesReceived           = errors.New("no more pieces required for decoding")
	ErrMoreUsefulPiecesRequired          = errors.New("not enough pieces received yet to decode")
	ErrCopyFailedDuringPieceConstruction = errors.New("failed to copy whole data before splitting into pieces")
	ErrPieceCountMoreThanTotalBytes      = errors.New("requested piece count > total bytes of original data")
	ErrZeroPieceSize                     = errors.New("pieces can't be sized as zero byte")
	ErrBadPieceCount                     = errors.New("minimum 2 pieces required for RLNC")
	ErrCodedDataLengthMismatch           = errors.New("coded data length != coded piece count x coded piece length")
	ErrCodingVectorLengthMismatch        = errors.New("coding vector length > coded piece length ( in total )")
)
