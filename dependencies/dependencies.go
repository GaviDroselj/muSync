package dependencies

import (
	"context"
	"log/slog"
	"os"

	"github.com/GaviDroselj/go-ytdlp"
)

func DonwloadDependencies(logger *slog.Logger) {
	logger.Info("Downloading yt-dlp and dependencies...")

	logger.Debug("Downloading yt-dlp...")
	_, err := ytdlp.Install(context.TODO(), &ytdlp.InstallOptions{})
	if err != nil {
		logger.Error("Failed to download yt-dlp , terminating", "err", err)
		os.Exit(1)
	}
	logger.Debug("Successfully downloaded yt-dlp")

	logger.Debug("Downloading ffmpeg...")
	_, err = ytdlp.InstallFFmpeg(context.TODO(), &ytdlp.InstallFFmpegOptions{})
	if err != nil {
		logger.Error("Failed to download ffmpeg , terminating", "err", err)
		os.Exit(1)
	}
	logger.Debug("Successfully downloaded ffmpeg")

	logger.Debug("Downloading ffprobe...")
	_, err = ytdlp.InstallFFprobe(context.TODO(), &ytdlp.InstallFFmpegOptions{})
	if err != nil {
		logger.Error("Failed to download ffprobe , terminating", "err", err)
		os.Exit(1)
	}
	logger.Debug("Successfully downloaded ffprobe")
}
