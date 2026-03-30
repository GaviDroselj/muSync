package main

import (
	"log/slog"
	"os"

	"github.com/gavidroselj/musync/dependencies"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	dependencies.DonwloadDependencies(logger)
}
