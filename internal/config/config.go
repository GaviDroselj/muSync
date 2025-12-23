package config

import (
	"encoding/json"
	"log/slog"
	"os"
)

const ConfigPath = "config.json"

type Config struct {
	LogLevel    slog.Level      `json:"log_level"`
	MusicFolder string          `json:"music_folder"`
	Playlists   []PlaylistEntry `json:"playlists"`
}

type PlaylistEntry struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	Subfolder string `json:"subfolder"`
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	conf := Config{}
	err = json.Unmarshal(data, &conf)
	if err != nil {
		return Config{}, err
	}

	return conf, nil
}
