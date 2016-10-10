package plugin

import (
	"github.com/disktnk/sb-audio"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
)

func init() {
	bql.MustRegisterGlobalSourceCreator("audio_device", bql.SourceCreatorFunc(
		audio.NewDeviceSource))
}
