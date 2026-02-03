package downloader

import (
	"time"

	"github.com/GaviDroselj/muSync/internal/config"
)

func Schedul(pl []*Playlist, ds config.DownloadStrategyType) {

	index := 0
	length := len(pl)

	for {
		switch ds {
		case config.Priority:
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
		case config.RoundRobin:
			for i := index; i < length; i++ {
				if pl[i].LastUpdate.Before(time.Now().Add(-time.Hour * 6)) {
					pl[i].Update()
					break
				}

				downloaded := pl[i].ProcessQueue()
				if downloaded {
					break
				}
			}
			index++
			if index%length == 0 {
				index = 0
			}
		}
		time.Sleep(time.Minute)
	}
}
