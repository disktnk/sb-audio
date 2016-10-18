package audio

import (
	"bytes"
	"encoding/binary"
)

func aiffFormat(in []byte) ([]byte, error) {
	// below setup, skip errors when write to buffer
	b := bytes.NewBuffer([]byte{})
	nSamples := len(in)
	totalBytes := 4 + 8 + 18 + 8 + 8 + 4*nSamples

	b.WriteString("FORM")
	binary.Write(b, binary.BigEndian, int32(totalBytes)) // total bytes
	b.WriteString("AIFF")

	// common chunk
	b.WriteString("COMM")
	binary.Write(b, binary.BigEndian, int32(18))              // size
	binary.Write(b, binary.BigEndian, int16(1))               // channels
	binary.Write(b, binary.BigEndian, int32(nSamples))        // number of samples
	binary.Write(b, binary.BigEndian, int16(32))              // bits per simple
	b.Write([]byte{0x40, 0x0e, 0xac, 0x44, 0, 0, 0, 0, 0, 0}) // 80-bit sample rate 4410

	// sound chunk
	b.WriteString("SSND")
	binary.Write(b, binary.BigEndian, int32(4*nSamples+8)) // size
	binary.Write(b, binary.BigEndian, int32(0))            // offset
	binary.Write(b, binary.BigEndian, int32(0))            // block

	binary.Write(b, binary.BigEndian, in)
	return b.Bytes(), nil
}
