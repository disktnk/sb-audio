package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gordonklaus/portaudio"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"time"
)

var (
	tickPath = data.MustCompilePath("tick")
)

// NewDeviceSource returns audio source from default device.
func NewDeviceSource(ctx *core.Context, ioParams *bql.IOParams, params data.Map) (
	core.Source, error) {
	tick := 3
	if ti, err := params.Get(tickPath); err == nil {
		itick, err := data.AsInt(ti)
		if err != nil {
			return nil, err
		}
		tick = int(itick)
	}

	return &device{
		tick: tick,
		stop: make(chan struct{}),
	}, nil
}

type device struct {
	tick int
	stop chan struct{}
}

func (d *device) GenerateStream(ctx *core.Context, w core.Writer) error {
	if err := portaudio.Initialize(); err != nil {
		return fmt.Errorf(
			"PortAudio library is failed to initialize: %v", err)
	}
	defer portaudio.Terminate()

	in := make([]int32, 64)
	stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
	if err != nil {
		return fmt.Errorf("cannot open default device: %v", err)
	}
	defer stream.Close()

	if err := stream.Start(); err != nil {
		return fmt.Errorf("cannot start audio stream: %v", err)
	}
	defer stream.Stop()

	b := bytes.NewBuffer([]byte{})
	ticker := time.NewTicker(time.Duration(d.tick) * time.Second)
	for {
		if err := stream.Read(); err != nil {
			return fmt.Errorf("cannot read audio data from device: %v",
				err)
		}
		if err := binary.Write(b, binary.BigEndian, in); err != nil {
			return err
		}

		select {
		case <-ticker.C:
			audio, err := soundStyle(b.Bytes())
			if err != nil {
				return err
			}

			da := data.Map{
				"type": data.String("AIFF"),
				"data": data.Blob(audio),
			}
			tu := core.NewTuple(da)
			if err := w.Write(ctx, tu); err != nil {
				return err
			}

			b = bytes.NewBuffer([]byte{})
		case <-d.stop:
			return nil
		default:
		}
	}
	return nil
}

func soundStyle(in []byte) ([]byte, error) {
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

func (d *device) Stop(ctx *core.Context) error {
	close(d.stop)
	return nil
}
