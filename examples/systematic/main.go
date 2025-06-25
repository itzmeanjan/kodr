package main

import (
	"bytes"
	"crypto/sha3"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/systematic"
)

type Data struct {
	FieldA uint    `json:"fieldA"`
	FieldB float64 `json:"fieldB"`
	FieldC bool    `json:"fieldC"`
	FieldD []byte  `json:"fieldD"`
}

// Generates random byte array of size N
func generateData(n uint) []byte {
	container := make([]byte, 0, n)
	for range n {
		container = append(container, byte(rand.Intn(255)))
	}

	return container
}

// Generates random `Data` i.e. values associated with
// respective fields are random
func randData() *Data {
	d := Data{
		FieldA: uint(rand.Uint64()),
		FieldB: rand.Float64(),
		FieldC: rand.Intn(2) == 0,
		FieldD: generateData(uint(1<<10 + rand.Intn(1<<10))),
	}
	return &d
}

func main() {
	data := randData()

	m_data, err := json.Marshal(&data)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	log.Printf("Original serialised data of %d bytes\n", len(m_data))

	hasher := sha3.New256()
	hasher.Write(m_data)
	originalSum := hasher.Sum(nil)
	log.Printf("SHA3-256(original): 0x%s\n", hex.EncodeToString(originalSum))

	var (
		pieceSize uint = 1 << 3 // in bytes
	)

	enc, err := systematic.NewSystematicRLNCEncoderWithPieceSize(m_data, pieceSize)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	log.Printf("%d pieces being coded together, each of %d bytes\n", enc.PieceCount(), enc.PieceSize())
	log.Printf("%d bytes of padding used\n\n", enc.Padding())
	log.Printf("%d bytes of coded data to be consumed for successful decoding\n", enc.DecodableLen())

	dec := systematic.NewSystematicRLNCDecoder(enc.PieceCount())
	for {
		c_piece := enc.CodedPiece()

		// simulating piece drop/ loss
		if rand.Intn(2) == 0 {
			continue
		}

		err := dec.AddPiece(c_piece)
		if err != nil && errors.Is(err, kodr.ErrAllUsefulPiecesReceived) {
			break
		}
	}

	d_pieces, err := dec.GetPieces()
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	d_flattened := make([]byte, 0, len(m_data)+int(enc.Padding()))
	for i := range d_pieces {
		d_flattened = append(d_flattened, d_pieces[i]...)
	}

	log.Printf("Recovered %d ( = %d + %d ) bytes flattened data\n", len(d_flattened), len(m_data), enc.Padding())
	d_flattened = d_flattened[:len(m_data)]

	hasher.Reset()
	hasher.Write(d_flattened)
	recoveredSum := hasher.Sum(nil)
	log.Printf("SHA3-256(recovered): 0x%s\n", hex.EncodeToString(recoveredSum))

	if !bytes.Equal(originalSum, recoveredSum) {
		log.Fatalln("SHA3-256 digest of original data and recovered data doesn't match !")
	}

	var rec_data Data
	if err := json.Unmarshal(d_flattened, &rec_data); err != nil {
		log.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}
