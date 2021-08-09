package systematic_test

import (
	"bytes"
	"errors"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/systematic"
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
func generatePieces(pieceCount uint, pieceLength uint) []kodr.Piece {
	pieces := make([]kodr.Piece, 0, pieceCount)
	for i := 0; i < int(pieceCount); i++ {
		pieces = append(pieces, generateData(pieceLength))
	}
	return pieces
}

func TestSystematicRLNCCoding(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var (
		pieceCount      uint                              = uint(2<<1 + rand.Intn(2<<8))
		pieceLength     uint                              = 8192
		codedPieceCount uint                              = pieceCount * 2
		pieces          []kodr.Piece                      = generatePieces(pieceCount, pieceLength)
		enc             *systematic.SystematicRLNCEncoder = systematic.NewSystematicRLNCEncoder(pieces)
	)

	for i := 0; i < int(codedPieceCount); i++ {
		c_piece := enc.CodedPiece()
		if i < int(pieceCount) {
			if !c_piece.IsSystematic() {
				t.Fatal("expected piece to be systematic coded")
			}
		} else {
			if c_piece.IsSystematic() {
				t.Fatal("expected piece to be random coded")
			}
		}
	}
}

func TestNewSystematicRLNC(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	t.Run("Encoder", func(t *testing.T) {
		var (
			pieceCount  uint = 1 << 8
			pieceLength uint = 8192
		)

		pieces := generatePieces(pieceCount, pieceLength)
		enc := systematic.NewSystematicRLNCEncoder(pieces)
		dec := systematic.NewSystematicRLNCDecoder(pieceCount)

		encoderFlow(t, enc, dec, pieceCount, pieces)
	})

	t.Run("EncoderWithPieceCount", func(t *testing.T) {
		size := uint(2<<10 + rand.Intn(2<<10))
		pieceCount := uint(2<<1 + rand.Intn(2<<8))
		data := generateData(size)

		enc, err := systematic.NewSystematicRLNCEncoderWithPieceCount(data, pieceCount)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		pieces, _, err := kodr.OriginalPiecesFromDataAndPieceCount(data, pieceCount)
		if err != nil {
			t.Fatal(err.Error())
		}

		dec := systematic.NewSystematicRLNCDecoder(pieceCount)
		encoderFlow(t, enc, dec, pieceCount, pieces)
	})

	t.Run("EncoderWithPieceSize", func(t *testing.T) {
		size := uint(2<<10 + rand.Intn(2<<10))
		pieceSize := uint(2<<5 + rand.Intn(2<<5))
		pieceCount := uint(math.Ceil(float64(size) / float64(pieceSize)))
		data := generateData(size)

		enc, err := systematic.NewSystematicRLNCEncoderWithPieceSize(data, pieceSize)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		pieces, _, err := kodr.OriginalPiecesFromDataAndPieceSize(data, pieceSize)
		if err != nil {
			t.Fatal(err.Error())
		}

		dec := systematic.NewSystematicRLNCDecoder(pieceCount)
		encoderFlow(t, enc, dec, pieceCount, pieces)
	})
}

func encoderFlow(t *testing.T, enc *systematic.SystematicRLNCEncoder, dec *systematic.SystematicRLNCDecoder, pieceCount uint, pieces []kodr.Piece) {
	for {
		c_piece := enc.CodedPiece()

		if rand.Intn(2) == 0 {
			continue
		}

		if err := dec.AddPiece(c_piece); err != nil && errors.Is(err, kodr.ErrAllUsefulPiecesReceived) {
			break
		}
	}

	d_pieces, err := dec.GetPieces()
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(pieces) != len(d_pieces) {
		t.Fatal("didn't decode all !")
	}

	for i := 0; i < int(pieceCount); i++ {
		if !bytes.Equal(pieces[i], d_pieces[i]) {
			t.Fatal("decoded data doesn't match !")
		}
	}
}

func TestSystematicRLNCEncoder_Padding(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	t.Run("WithPieceCount", func(t *testing.T) {
		for i := 0; i < 1<<5; i++ {
			size := uint(2<<10 + rand.Intn(2<<10))
			pieceCount := uint(2<<1 + rand.Intn(2<<8))
			data := generateData(size)

			enc, err := systematic.NewSystematicRLNCEncoderWithPieceCount(data, pieceCount)
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

			enc, err := systematic.NewSystematicRLNCEncoderWithPieceSize(data, pieceSize)
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

func TestSystematicRLNCEncoder_CodedPieceLen(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	t.Run("WithPieceCount", func(t *testing.T) {
		size := uint(2<<10 + rand.Intn(2<<10))
		pieceCount := uint(2<<1 + rand.Intn(2<<8))
		data := generateData(size)

		enc, err := systematic.NewSystematicRLNCEncoderWithPieceCount(data, pieceCount)
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

		enc, err := systematic.NewSystematicRLNCEncoderWithPieceSize(data, pieceSize)
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

func TestSystematicRLNCEncoder_DecodableLen(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	flow := func(enc *systematic.SystematicRLNCEncoder, dec *systematic.SystematicRLNCDecoder) {
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

		enc, err := systematic.NewSystematicRLNCEncoderWithPieceCount(data, pieceCount)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		dec := systematic.NewSystematicRLNCDecoder(pieceCount)
		flow(enc, dec)
	})

	t.Run("WithPieceSize", func(t *testing.T) {
		size := uint(2<<10 + rand.Intn(2<<10))
		pieceSize := uint(2<<5 + rand.Intn(2<<5))
		pieceCount := uint(math.Ceil(float64(size) / float64(pieceSize)))
		data := generateData(size)

		enc, err := systematic.NewSystematicRLNCEncoderWithPieceSize(data, pieceSize)
		if err != nil {
			t.Fatalf("Error: %s\n", err.Error())
		}

		dec := systematic.NewSystematicRLNCDecoder(pieceCount)
		flow(enc, dec)
	})
}
