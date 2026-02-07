package downloader

import (
	"time"

	"github.com/GaviDroselj/muSync/internal/config"
)

func Schedule(playlists []*Playlist, downloadStrategy config.DownloadStrategyType) {

	index := -1
	length := len(playlists)

	for {
		var start int
		switch downloadStrategy {
		case config.Priority:
			start = 0
		case config.RoundRobin:
			index = RoundRobin(playlists, index, length)
			start = index
		}

		for i := start; i < len(playlists); i++ {
			if playlists[i].LastUpdate.Before(time.Now().Add(-time.Hour * 6)) {
				playlists[i].Update()
				break
			}

			queueHasItems := playlists[i].ProcessQueue()
			if queueHasItems {
				break
			}
		}

		//time.Sleep(time.Minute)
	}
}

func RoundRobin(pl []*Playlist, i int, l int) int {
	i++
	if i%l == 0 {
		i = 0
	}
	return i
}
