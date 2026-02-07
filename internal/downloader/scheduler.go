package downloader

import (
	"time"

	"github.com/GaviDroselj/muSync/internal/config"
)

var roundRobinIndex = 0

func Schedule(playlists []*Playlist, downloadStrategy config.DownloadStrategyType) {

	for {
		var playlist *Playlist
		switch downloadStrategy {
		case config.Priority:
			playlist = priority(playlists)
		case config.RoundRobin:
			playlist = roundRobin(playlists)
		}
		playlist.Update()
		time.Sleep(time.Minute)
	}
}

// Find closest possible playlist with download attempt
func findNextPlaylist(playlists []*Playlist, start int) *Playlist {
	for i := 0; i < len(playlists); i++ {
		playlist := playlists[(i+start)%len(playlists)]
		if playlist.NeedsUpdate() {
			return playlist
		}
	}
	return nil
}

func roundRobin(playlists []*Playlist) *Playlist {
	playlist := findNextPlaylist(playlists, roundRobinIndex)
	roundRobinIndex = (roundRobinIndex + 1) % len(playlists)

	return playlist
}

func priority(playlists []*Playlist) *Playlist {
	return findNextPlaylist(playlists, 0)
}
