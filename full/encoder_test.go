package full_test

import (
	"bytes"
	"errors"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/full"
	"github.com/itzmeanjan/kodr/kodr_internals"
)

// Generates `N`-bytes of random data from default
// randomization source
func generateData(n uint) []byte {
	data := make([]byte, n)
	// can safely ignore error
	rand.Read(data)
	return data
}

// Generates N-many pieces each of M-bytes length, to be used
// for testing purposes
func generatePieces(pieceCount uint, pieceLength uint) []kodr_internals.Piece {
	pieces := make([]kodr_internals.Piece, 0, pieceCount)
	for i := 0; i < int(pieceCount); i++ {
		pieces = append(pieces, generateData(pieceLength))
	}
	return pieces
}

func encoderFlow(t *testing.T, enc *full.FullRLNCEncoder, pieceCount, codedPieceCount int, pieces []kodr_internals.Piece) {
	coded := make([]*kodr_internals.CodedPiece, 0, codedPieceCount)
	for i := 0; i < codedPieceCount; i++ {
		coded = append(coded, enc.CodedPiece())
	}

	dec := full.NewFullRLNCDecoder(uint(pieceCount))
	for i := 0; i < codedPieceCount; i++ {
		if i < pieceCount {
			if _, err := dec.GetPieces(); !(err != nil && errors.Is(err, kodr.ErrMoreUsefulPiecesRequired)) {
				t.Fatal("expected error indicating more pieces are required for decoding")
			}
		}

		if err := dec.AddPiece(coded[i]); errors.Is(err, kodr.ErrAllUsefulPiecesReceived) {
			break
		}
	}

	if !dec.IsDecoded() {
		t.Fatal("expected to be fully decoded !")
	}

	for i := 0; i < codedPieceCount-pieceCount; i++ {
		if err := dec.AddPiece(coded[pieceCount+i]); !(err != nil && errors.Is(err, kodr.ErrAllUsefulPiecesReceived)) {
			t.Fatal("expected error indication, received nothing !")
		}
	}

	d_pieces, err := dec.GetPieces()
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(pieces) != len(d_pieces) {
		t.Fatal("didn't decode all !")
	}

	for i := 0; i < pieceCount; i++ {
		if !bytes.Equal(pieces[i], d_pieces[i]) {
			t.Fatal("decoded data doesn't match !")
		}
	}
}

func TestNewFullRLNCEncoder(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	pieceCount := 128
	pieceLength := 8192
	codedPieceCount := pieceCount + 2
	pieces := generatePieces(uint(pieceCount), uint(pieceLength))
	enc := full.NewFullRLNCEncoder(pieces)

	encoderFlow(t, enc, pieceCount, codedPieceCount, pieces)
}

func TestNewFullRLNCEncoderWithPieceCount(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	size := uint(2<<10 + rand.Intn(2<<10))
	pieceCount := uint(2<<1 + rand.Intn(2<<8))
	codedPieceCount := pieceCount + 2
	data := generateData(size)
	t.Logf("\nTotal Data: %d bytes\nPiece Count: %d\nCoded Piece Count: %d\n", size, pieceCount, codedPieceCount)

	pieces, _, err := kodr_internals.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
	if err != nil {
		t.Fatal(err.Error())
	}

	enc, err := full.NewFullRLNCEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		t.Fatal(err.Error())
	}

	encoderFlow(t, enc, int(pieceCount), int(codedPieceCount), pieces)
}

func TestNewFullRLNCEncoderWithPieceSize(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	size := uint(2<<10 + rand.Intn(2<<10))
	pieceSize := uint(2<<5 + rand.Intn(2<<5))
	pieceCount := int(math.Ceil(float64(size) / float64(pieceSize)))
	codedPieceCount := pieceCount + 2
	data := generateData(size)
	t.Logf("\nTotal Data: %d bytes\nPiece Size: %d bytes\nPiece Count: %d\nCoded Piece Count: %d\n", size, pieceSize, pieceCount, codedPieceCount)

	pieces, _, err := kodr_internals.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
	if err != nil {
		t.Fatal(err.Error())
	}

	enc, err := full.NewFullRLNCEncoderWithPieceSize(data, pieceSize)
	if err != nil {
		t.Fatal(err.Error())
	}

	encoderFlow(t, enc, pieceCount, codedPieceCount, pieces)
}

func TestFullRLNCEncoderPadding(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	t.Run("WithPieceCount", func(t *testing.T) {
		for i := 0; i < 1<<5; i++ {
			size := uint(2<<10 + rand.Intn(2<<10))
			pieceCount := uint(2<<1 + rand.Intn(2<<8))
			data := generateData(size)

			enc, err := full.NewFullRLNCEncoderWithPieceCount(data, pieceCount)
			if err != nil {
				t.Fatalf("Error: %s\n", err.Error())
			}

			extra := enc.Padding()
			pieceSize := (size + extra) / pieceCount
			c_piece := enc.CodedPiece()
			if uint(len(c_piece.Piece)) != pieceSize {
				t.Fatalf("expected pieceSize to be %dB, found to be %dB\n", pieceSize, len(c_piece.Piece))
			}
		}
	})

	t.Run("WithPieceSize", func(t *testing.T) {
		for i := 0; i < 1<<5; i++ {
			size := uint(2<<10 + rand.Intn(2<<10))
			pieceSize := uint(2<<5 + rand.Intn(2<<5))
			pieceCount := uint(math.Ceil(float64(size) / float64(pieceSize)))
			data := generateData(size)

			enc, err := full.NewFullRLNCEncoderWithPieceSize(data, pieceSize)
			if err != nil {
				t.Fatalf("Error: %s\n", err.Error())
			}

			extra := enc.Padding()
			c_pieceSize := (size + extra) / pieceCount
			c_piece := enc.CodedPiece()
			if pieceSize != c_pieceSize || uint(len(c_piece.Piece)) != pieceSize {
				t.Fatalf("expected pieceSize to be %dB, found to be %dB\n", c_pieceSize, len(c_piece.Piece))
			}
		}
	})
}

func TestFullRLNCEncoder_CodedPieceLen(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	t.Run("WithPieceCount", func(t *testing.T) {
		size := uint(2<<10 + rand.Intn(2<<10))
		pieceCount := uint(2<<1 + rand.Intn(2<<8))
		data := generateData(size)

		enc, err := full.NewFullRLNCEncoderWithPieceCount(data, pieceCount)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		for i := 0; i <= int(pieceCount); i++ {
			c_piece := enc.CodedPiece()
			if c_piece.Len() != enc.CodedPieceLen() {
				t.Fatalf("expected coded piece to be of %dB, found to be of %dB\n", enc.CodedPieceLen(), c_piece.Len())
			}
		}
	})

	t.Run("WithPieceSize", func(t *testing.T) {
		size := uint(2<<10 + rand.Intn(2<<10))
		pieceSize := uint(2<<5 + rand.Intn(2<<5))
		pieceCount := uint(math.Ceil(float64(size) / float64(pieceSize)))
		data := generateData(size)

		enc, err := full.NewFullRLNCEncoderWithPieceSize(data, pieceSize)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		for i := 0; i <= int(pieceCount); i++ {
			c_piece := enc.CodedPiece()
			if c_piece.Len() != enc.CodedPieceLen() {
				t.Fatalf("expected coded piece to be of %dB, found to be of %dB\n", enc.CodedPieceLen(), c_piece.Len())
			}
		}
	})
}

func TestFullRLNCEncoder_DecodableLen(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	flow := func(enc *full.FullRLNCEncoder, dec *full.FullRLNCDecoder) {
		consumed_len := uint(0)
		for !dec.IsDecoded() {
			c_piece := enc.CodedPiece()
			// randomly drop piece
			if rand.Intn(2) == 0 {
				continue
			}
			if err := dec.AddPiece(c_piece); errors.Is(err, kodr.ErrAllUsefulPiecesReceived) {
				break
			}

			// as consumed this piece --- accounting
			consumed_len += c_piece.Len()
		}

		if consumed_len < enc.DecodableLen() {
			t.Fatalf("expected to consume >=%dB for decoding, but actually consumed %dB\n", enc.DecodableLen(), consumed_len)
		}
	}

	t.Run("WithPieceCount", func(t *testing.T) {
		size := uint(2<<10 + rand.Intn(2<<10))
		pieceCount := uint(2<<1 + rand.Intn(2<<8))
		data := generateData(size)

		enc, err := full.NewFullRLNCEncoderWithPieceCount(data, pieceCount)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		dec := full.NewFullRLNCDecoder(pieceCount)
		flow(enc, dec)
	})

	t.Run("WithPieceSize", func(t *testing.T) {
		size := uint(2<<10 + rand.Intn(2<<10))
		pieceSize := uint(2<<5 + rand.Intn(2<<5))
		pieceCount := uint(math.Ceil(float64(size) / float64(pieceSize)))
		data := generateData(size)

		enc, err := full.NewFullRLNCEncoderWithPieceSize(data, pieceSize)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		dec := full.NewFullRLNCDecoder(pieceCount)
		flow(enc, dec)
	})
}
