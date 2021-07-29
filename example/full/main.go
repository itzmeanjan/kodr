package main

import (
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"log"
	"math/rand"
	"os"
	"path"

	"github.com/itzmeanjan/kodr"
	"github.com/itzmeanjan/kodr/full"
)

func main() {
	img := path.Join("..", "..", "img", "logo.png")
	log.Printf("Reading from %s\n", img)
	data, err := os.ReadFile(img)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	log.Printf("Read %d bytes\n", len(data))

	hasher := sha512.New512_224()
	hasher.Write(data)
	sum := hasher.Sum(nil)
	log.Printf("SHA512: 0x%s\n\n", hex.EncodeToString(sum))

	var (
		pieceCount        uint = 64
		codedPieceCount   uint = 2 * pieceCount
		droppedPieceCount uint = pieceCount / 2
	)

	log.Printf("Coding %d pieces together\n", pieceCount)
	enc, err := full.NewFullRLNCEncoderWithPieceCount(data, pieceCount)
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	codedPieces := make([]*kodr.CodedPiece, 0, codedPieceCount)
	for i := 0; i < int(codedPieceCount); i++ {
		codedPieces = append(codedPieces, enc.CodedPiece())
	}

	log.Printf("Coded into %d pieces\n", codedPieceCount)

	for i := 0; i < int(droppedPieceCount); i++ {
		idx := rand.Intn(len(codedPieces))
		codedPieces[idx] = nil
		copy(codedPieces[idx:], codedPieces[idx+1:])
		codedPieces[len(codedPieces)-1] = nil
		codedPieces = codedPieces[:len(codedPieces)-1]
	}

	log.Printf("Dropped %d pieces, remaining %d pieces\n\n", droppedPieceCount, len(codedPieces))

	var (
		recodedPieceCount uint = uint(len(codedPieces)) * 2
	)

	rec := full.NewFullRLNCRecoder(codedPieces)
	recodedPieces := make([]*kodr.CodedPiece, 0, recodedPieceCount)
	for i := 0; i < int(recodedPieceCount); i++ {
		rec_p, err := rec.CodedPiece()
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
		recodedPieces = append(recodedPieces, rec_p)
	}

	log.Printf("Recoded %d coded pieces into %d pieces\n", len(codedPieces), recodedPieceCount)
	rand.Shuffle(int(recodedPieceCount), func(i, j int) {
		recodedPieces[i], recodedPieces[j] = recodedPieces[j], recodedPieces[i]
	})
	log.Printf("Shuffled %d coded pieces\n", recodedPieceCount)

	for i := 0; i < int(recodedPieceCount)/2; i++ {
		idx := rand.Intn(len(recodedPieces))
		recodedPieces[idx] = nil
		copy(recodedPieces[idx:], recodedPieces[idx+1:])
		recodedPieces[len(recodedPieces)-1] = nil
		recodedPieces = recodedPieces[:len(recodedPieces)-1]
	}

	log.Printf("Dropped %d pieces, remaining %d pieces\n\n", recodedPieceCount/2, len(recodedPieces))

	log.Printf("Decoding with %d pieces\n", pieceCount)
	dec := full.NewFullRLNCDecoder(pieceCount)
	for i := 0; i < int(pieceCount); i++ {
		if err := dec.AddPiece(recodedPieces[i]); err != nil {
			log.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
	}

	if err := dec.AddPiece(codedPieces[pieceCount]); !(err != nil && errors.Is(err, kodr.ErrAllUsefulPiecesReceived)) {
		log.Printf("Error `%s` was expected to be thrown\n", kodr.ErrAllUsefulPiecesReceived)
	}

	dec_p, err := dec.GetPieces()
	if err != nil {
		log.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	decoded_data := make([]byte, 0)
	for i := 0; i < len(dec_p); i++ {
		decoded_data = append(decoded_data, dec_p[i]...)
	}

	log.Printf("Decoded into %d bytes\n", len(decoded_data))

	if !bytes.Equal(data, decoded_data[:len(data)]) {
		log.Println("Decoded data not matching !")
		os.Exit(1)
	}

	log.Printf("First %d bytes of decoded data matches original %d bytes\n", len(data), len(data))
	log.Printf("3 bytes of padding: %v\n\n", decoded_data[len(data):])

	hasher.Reset()
	hasher.Write(decoded_data[:len(data)])
	sum = hasher.Sum(nil)
	log.Printf("SHA512: 0x%s\n", hex.EncodeToString(sum))

	if err := os.WriteFile("recovered.png", decoded_data[:len(data)], 0o644); err != nil {
		log.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	log.Printf("Wrote %d bytes into `./recovered.png`", len(data))
}
