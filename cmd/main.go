package main

import (
	"log/slog"
	"os"

	"github.com/gavidroselj/musync/dependencies"
	"github.com/gavidroselj/musync/internal/config"
	"github.com/gavidroselj/musync/internal/downloader"
	"github.com/lrstanley/go-ytdlp"
)

func main() {
	logLevel := new(slog.LevelVar)
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

	dependencies.DonwloadDependencies(logger)

	logger.Info("muSync started", "version", ytdlp.Version, "logLevel", logLevel)

	var playlists []*downloader.Playlist

	for _, confEntry := range conf.Playlists {
		newPlaylist := downloader.NewPlaylist(conf, confEntry, logger)
		newPlaylist.SyncFromDisk()

		playlists = append(playlists, newPlaylist)
	}

	downloader.Schedule(playlists, conf.DownloadStrategy)
}
