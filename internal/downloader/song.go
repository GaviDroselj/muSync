package downloader

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/lrstanley/go-ytdlp"
)

const deletionGracePeriod = 1

type Song struct {
	ID                string `json:"id"`
	URL               string `json:"url"`
	Title             string `json:"title"`
	DeletedIterations int    `json:"deleted_iterations"`

	logger *slog.Logger `json:"-"`
}

func NewSong(id, url, title string, logger *slog.Logger) Song {
	return Song{
		ID:     id,
		URL:    url,
		Title:  title,
		logger: logger.With("songID", id),
	}
}

func NewSongFromInfo(info *ytdlp.ExtractedInfo, logger *slog.Logger) Song {
	url := SafeDeref(info.URL)
	if url == "" {
		url = info.ExtractedFormat.URL
	}
	if url == "" {
		url = SafeDeref(info.WebpageURL)
	}

	return NewSong(info.ID, url, SafeDeref(info.Title), logger)
}

func (s *Song) ResetDeleteCounter() {
	s.DeletedIterations = 0
}

// Delete song if it has been absent for deletionGracePeriod playlist refreshes
// This ensures a single bad response from server will not wipe out the entire library
func (s *Song) DeleteAttempt(folderPath string) bool {
	s.DeletedIterations++
	if s.DeletedIterations < deletionGracePeriod {
		return false
	}

	err := s.Delete(folderPath)
	if err != nil {
		slog.Error("Failed to delete song", "song", *s, "folderPath", folderPath, "err", err)
		return false
	}
	return true
}

// Deletes song by searching for file containing s.ID in folderPath
func (s *Song) Delete(folderPath string) error {
	if folderPath == "" {
		return nil
	}

	dir, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, entry := range dir {
		if strings.Contains(entry.Name(), s.ID) {
			slog.Debug("Deleting song file", "path", folderPath, "file", entry.Name())
			err := os.Remove(fmt.Sprintf("%s/%s", folderPath, entry.Name()))

			if err != nil {
				return err
			}
			return nil
		}
	}

	slog.Debug("Failed to find song scheduled for deletion on disk", "song", *s, "dirEntries", dir)
	return nil
}
