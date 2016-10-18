package audio

import (
	"bytes"
	"encoding/binary"
)

// aiffFormat convert "in" binary to AIFF format style.
// Skip error in writing to buffer
func aiffFormat(in []byte) ([]byte, error) {
	// ckID            4  "FORM"
	// ckSize          4  =total size
	// formType        4  "AIFF"
	// chID            4  "COMM"
	// ckSize          4  18
	// numChannels     2  1 (default device is set as 1)
	// numSampleFrames 4  88200(=44.1kHz * 2(16bit) * 1(monaural))
	// sampleSize      2  16
	// sampleRate      10 44100
	// ckID            4  "SSND"
	// ckSize          4
	// offset          4  0
	// blockSize       4  0
	// soundData          =in
	b := bytes.NewBuffer([]byte{})
	nSamples := len(in)
	totalBytes := (4 + 4) + (4 + 18) + (4 + 12 + nSamples)

	b.WriteString("FORM")
	binary.Write(b, binary.BigEndian, int32(totalBytes-8)) // total bytes
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
	binary.Write(b, binary.BigEndian, int32(nSamples)) // size
	binary.Write(b, binary.BigEndian, int32(0))        // offset
	binary.Write(b, binary.BigEndian, int32(0))        // block

	binary.Write(b, binary.BigEndian, in)
	return b.Bytes(), nil
}
