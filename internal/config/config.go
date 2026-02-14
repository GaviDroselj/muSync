package config

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

const ConfigPath = "config.yml"

type DownloadStrategyType string

const (
	Priority   DownloadStrategyType = "PRIORITY"
	RoundRobin DownloadStrategyType = "ROUND_ROBIN"
)

type Config struct {
	LogLevel         slog.Level           `yaml:"log_level"`
	MusicFolder      string               `yaml:"music_folder"`
	DownloadStrategy DownloadStrategyType `yaml:"download_strategy"`
	Playlists        []PlaylistEntry      `yaml:"playlists"`
}

type PlaylistEntry struct {
	Name      string `yaml:"name"`
	URL       string `yaml:"url"`
	Subfolder string `yaml:"subfolder"`
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	conf := Config{}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
