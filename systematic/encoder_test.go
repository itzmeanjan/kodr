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
