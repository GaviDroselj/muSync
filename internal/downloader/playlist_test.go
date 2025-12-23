package downloader

import (
	"log/slog"
	"testing"

	"github.com/GaviDroselj/muSync/internal/config"
	"github.com/GaviDroselj/muSync/internal/xtesting"
)

func getTestingPlaylist() *Playlist {
	playlist := NewPlaylist(config.Config{
		MusicFolder: "music",
	}, config.PlaylistEntry{
		Name:      "test",
		URL:       "http://test.com/test",
		Subfolder: "test",
	}, slog.Default())

	playlist.Songs = map[string]*Song{
		"1": {
			ID:                "1",
			URL:               "http://test.com/1",
			Title:             "Song1",
			DeletedIterations: 0,
		},
		"2": {
			ID:                "2",
			URL:               "http://test.com/2",
			Title:             "Song2",
			DeletedIterations: 4,
		},
		"3": {
			ID:                "3",
			URL:               "http://test.com/3",
			Title:             "Song3",
			DeletedIterations: 1,
		},
	}
	return playlist
}

func TestUpdateTrackedSongs(t *testing.T) {
	t.Run("no-change", func(t *testing.T) {
		g := xtesting.NewGoldie(t)
		playlist := getTestingPlaylist()

		playlist.updateTrackedSongs([]Song{
			{
				ID:    "1",
				URL:   "http://test.com/1",
				Title: "Song1",
			},
			{
				ID:    "2",
				URL:   "http://test.com/2",
				Title: "Song2",
			},
			{
				ID:    "3",
				URL:   "http://test.com/3",
				Title: "Song3",
			},
		})

		g.AssertJson(t, "playlist", playlist)
	})

	t.Run("add-only", func(t *testing.T) {
		g := xtesting.NewGoldie(t)
		playlist := getTestingPlaylist()

		playlist.updateTrackedSongs([]Song{
			{
				ID:    "1",
				URL:   "http://test.com/1",
				Title: "Song1",
			},
			{
				ID:    "2",
				URL:   "http://test.com/2",
				Title: "Song2",
			},
			{
				ID:    "3",
				URL:   "http://test.com/3",
				Title: "Song3",
			},
			{
				ID:  "4",
				URL: "http://test.com/4",
			},
		})

		g.AssertJson(t, "playlist", playlist)
	})

	t.Run("delete-once", func(t *testing.T) {
		g := xtesting.NewGoldie(t)
		playlist := getTestingPlaylist()

		playlist.updateTrackedSongs([]Song{
			{
				ID:    "2",
				URL:   "http://test.com/2",
				Title: "Song2",
			},
		})

		g.AssertJson(t, "playlist", playlist)
	})

	t.Run("delete-success", func(t *testing.T) {
		g := xtesting.NewGoldie(t)
		playlist := getTestingPlaylist()

		playlist.updateTrackedSongs([]Song{
			{
				ID:    "1",
				URL:   "http://test.com/1",
				Title: "Song1",
			},
		})

		g.AssertJson(t, "playlist", playlist)
	})
}
