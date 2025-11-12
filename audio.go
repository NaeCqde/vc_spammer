package main

import (
	"strings"

	"github.com/NaeCqde/dca"
)

type Audio struct {
	Description string
	AudioFilter string `yaml:"audio_filter"`
	File        string
}

func GetAudio(audio Audio) (*dca.EncodeSession, error) {
	options := *dca.StdEncodeOptions
	options.Channels = 1
	options.Bitrate = 96
	options.RawOutput = true
	options.Application = dca.AudioApplicationLowDelay

	options.FFmpegExecutable = config.FFmpegExecutable
	options.AudioFilter = audio.AudioFilter
	path := audio.File

	if strings.HasPrefix(path, "https://") || strings.HasPrefix(path, "http://") {
		options.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/140.0.0.0 Safari/537.36 Edg/140.0.0.0"
	}

	return dca.EncodeFile(path, &options)
}
