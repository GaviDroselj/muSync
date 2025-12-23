package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/lrstanley/go-ytdlp"
)

func main() {
	slog.Info("Downloading yt-dlp dependencies...")
	_, err := ytdlp.InstallAll(context.TODO())
	if err != nil {
		slog.Error("Failed to download yt-dlp dependencies, terminating", "err", err)
		os.Exit(1)
	}
	slog.Info("Successfully downloaded yt-dlp dependencies")
}
