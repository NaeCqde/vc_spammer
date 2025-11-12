package main

import (
	"os"

	"github.com/goccy/go-yaml"
)

type Config struct {
	Prefix           string
	FFmpegExecutable string `yaml:"ffmpeg_executable"`
	Audio            map[string]Audio
	OwnerIDs         []string `yaml:"owner_ids"`
	Tokens           []string
}

var config Config

func LoadConfig() error {
	yml, err := os.ReadFile("config.yml")
	if err != nil {
		return err
	}

	return yaml.Unmarshal([]byte(yml), &config)
}
