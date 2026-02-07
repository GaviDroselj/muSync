package downloader

import (
	"time"

	"github.com/GaviDroselj/muSync/internal/config"
)

var index = 0

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
func possiblePlaylist(playlists []*Playlist, start int) *Playlist {
	for i := start; i < len(playlists); i++ {
		if playlists[i].NeedsUpdate() {
			return playlists[i]
		}
	}
	return nil
}

func roundRobin(playlists []*Playlist) *Playlist {
	playlist := possiblePlaylist(playlists, index)
	index++
	if index%len(playlists) == 0 {
		index = 0
	}
	return playlist
}

func priority(playlists []*Playlist) *Playlist {
	return possiblePlaylist(playlists, 0)
}
