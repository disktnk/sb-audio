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
	tickPath   = data.MustCompilePath("tick")
	foramtPath = data.MustCompilePath("format")
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
	format := "wav"
	if f, err := params.Get(foramtPath); err == nil {
		if format, err = data.AsString(f); err != nil {
			return nil, err
		}
	}
	var gen func([]byte) ([]byte, error)
	switch format {
	case "wav":
		gen = wavFormat
	case "aiff":
		gen = aiffFormat
	default:
		return nil, fmt.Errorf("'%v' is not supported", format)
	}

	return &device{
		tick:   tick,
		format: format,
		gen:    gen,
		stop:   make(chan struct{}),
	}, nil
}

type device struct {
	tick   int
	format string
	gen    func([]byte) ([]byte, error)
	stop   chan struct{}
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
			audio, err := d.gen(b.Bytes())
			if err != nil {
				return err
			}

			da := data.Map{
				"format": data.String(d.format),
				"data":   data.Blob(audio),
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

func (d *device) Stop(ctx *core.Context) error {
	close(d.stop)
	return nil
}
