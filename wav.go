package audio

import (
	"bytes"
	"encoding/binary"
)

// wavFormat convert "in" binary to WAVE format style.
// Skip eror in writing to buffer
func wavFormat(in []byte) ([]byte, error) {
	// RIFF header     4 "RIFF"
	// ckSize          4 =total size
	// formType        4 "WAVE"
	// chID            4 "fmt "
	// ckSize          4 16(PCM)
	// formatID        2 1(= WAVE_FORMAT_PCM)
	// numChannels     2 1 (default device is set as 1)
	// sampleRate      4 44100
	// numSampleFrames 4 88200(44.1kHz * 2 (16bit))
	// blockSize       2 2(16bit, monaural)
	// bitRate         2 16(bit)
	// extendSize      2 0 (not exist)
	// extendData      n not exist
	// dataCk          4 "data"
	// dataCkSize      4
	// soundData         =in
	b := bytes.NewBuffer([]byte{})
	nSamples := len(in)
	totalBytes := (4 + 4 + 4) + (4 + 4 + 16) + (4 + 4 + nSamples)

	b.WriteString("RIFF")
	binary.Write(b, binary.LittleEndian, int32(totalBytes-8))
	b.WriteString("WAVE")

	b.WriteString("fmt ")
	binary.Write(b, binary.LittleEndian, int32(16))
	binary.Write(b, binary.LittleEndian, int16(1))
	binary.Write(b, binary.LittleEndian, int16(1))
	binary.Write(b, binary.LittleEndian, int32(44100))
	binary.Write(b, binary.LittleEndian, int32(44100*4))
	binary.Write(b, binary.LittleEndian, int16(4))
	binary.Write(b, binary.LittleEndian, int16(32))

	b.WriteString("data")
	binary.Write(b, binary.LittleEndian, int32(nSamples))
	binary.Write(b, binary.LittleEndian, in)

	return b.Bytes(), nil
}
