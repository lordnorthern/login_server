package helpers

import (
	"bytes"
	"encoding/binary"
)

func HandleChunk(chunk []byte, inSize int) ([]ParsedChunk, error) {
	chunk = chunk[:inSize]
	var start = 0
	var chunks []ParsedChunk
	for {
		if len(chunk[start:]) < 17 {
			break
		}
		rChunk, err := divideChunk(chunk[start:])
		if err != nil {
			return nil, err
		}
		start += int(rChunk.Length) + 16
		chunks = append(chunks, rChunk)
	}

	return chunks, nil
}

func divideChunk(chunk []byte) (rChunk ParsedChunk, err error) {
	bufStep := bytes.NewReader(chunk[0:4])
	if len(chunk) < 17 {
		return
	}
	err = binary.Read(bufStep, binary.LittleEndian, &rChunk.Step)
	if err != nil {
		LogError(err, "divudeChunk", "Decode step")
		return
	}
	bufSerial := bytes.NewReader(chunk[4:8])
	err = binary.Read(bufSerial, binary.LittleEndian, &rChunk.Serial)
	if err != nil {
		LogError(err, "divudeChunk", "Decode serial")
		return
	}

	bufLength := bytes.NewReader(chunk[8:12])
	err = binary.Read(bufLength, binary.LittleEndian, &rChunk.Length)
	if err != nil {
		LogError(err, "divudeChunk", "Decode length")
		return
	}

	bufTotalChunks := bytes.NewReader(chunk[12:16])
	err = binary.Read(bufTotalChunks, binary.LittleEndian, &rChunk.TotalChunks)
	if err != nil {
		LogError(err, "divudeChunk", "Decode total chunks")
		return
	}
	rChunk.ChunkBytes = chunk[16 : rChunk.Length+16]
	return
}

type ParsedChunk struct {
	ChunkBytes  []byte
	Length      int32
	Step        int32
	Serial      int32
	TotalChunks int32
}
