package main

import (
	"log/slog"
	"os"

	"github.com/GaviDroselj/muSync/internal/downloader"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	downloader.DonwloadDependencies(logger)
}
