package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/GaviDroselj/muSync/internal/config"
	"github.com/GaviDroselj/muSync/internal/downloader"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	var logLevel = new(slog.LevelVar)
	logLevel.Set(slog.LevelDebug)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	slog.SetDefault(logger)

	conf, err := config.LoadConfig(config.ConfigPath)
	if err != nil {
		logger.Error("Failed to load config", "path", config.ConfigPath, "err", err)
		os.Exit(1)
	}
	logLevel.Set(conf.LogLevel)
	logger.Debug("Successfully loaded config", "conf", conf)

	logger.Info("Downloading yt-dlp and ffmpeg...")
	_, err = ytdlp.InstallAll(context.TODO())
	if err != nil {
		logger.Error("Failed to download yt-dlp, terminating", "err", err)
		os.Exit(1)
	}

	logger.Info("muSync started", "version", ytdlp.Version, "logLevel", logLevel)

	var playlists []*downloader.Playlist

	for _, confEntry := range conf.Playlists {
		newPlaylist := downloader.NewPlaylist(confEntry, logger)

		playlists = append(playlists, newPlaylist)
	}

	for {
		for _, playlist := range playlists {
			if playlist.LastUpdate.Before(time.Now().Add(-time.Hour * 6)) {
				playlist.Update()
				break
			}

			downloaded := playlist.ProcessQueue()
			if downloaded {
				break
			}
		}
		time.Sleep(time.Minute)
	}
}
