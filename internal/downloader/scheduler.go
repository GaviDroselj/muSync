package downloader

import (
	"time"

	"github.com/GaviDroselj/muSync/internal/config"
)

func Schedule(pl []*Playlist, ds config.DownloadStrategyType) {

	index := 0
	length := len(pl)

	for {
		switch ds {
		case config.Priority:
			Priority(pl)
		case config.RoundRobin:
			index = RoundRobin(pl, index, length)
		}
		//time.Sleep(time.Minute)
	}
}

func Priority(pl []*Playlist) {
	for _, playlist := range pl {
		if playlist.LastUpdate.Before(time.Now().Add(-time.Hour * 6)) {
			playlist.Update()
			break
		}

		downloaded := playlist.ProcessQueue()
		if downloaded {
			break
		}
	}
}

func RoundRobin(pl []*Playlist, i int, l int) int {
	if pl[i].LastUpdate.Before(time.Now().Add(-time.Hour * 6)) {
		pl[i].Update()
	}

	downloaded := pl[i].ProcessQueue()
	if downloaded {
	}

	i++
	if i%l == 0 {
		i = 0
	}
	return i
}
