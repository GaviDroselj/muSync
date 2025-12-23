package downloader

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/GaviDroselj/muSync/internal/config"
	"github.com/lrstanley/go-ytdlp"
)

var fileIDRegex = regexp.MustCompile(`.*\[(.*)\].mp3`)

type Playlist struct {
	Name          string           `json:"name"`
	URL           string           `json:"url"`
	Songs         map[string]*Song `json:"songs"`
	Folder        string           `json:"folder"`
	DownloadQueue []Song           `json:"download_queue"`
	LastUpdate    time.Time        `json:"last_update"`

	logger *slog.Logger `json:"-"`
}

func NewPlaylist(c config.Config, pe config.PlaylistEntry, logger *slog.Logger) *Playlist {
	playlist := Playlist{
		Name:          pe.Name,
		URL:           pe.URL,
		Folder:        path.Join(c.MusicFolder, pe.Subfolder),
		Songs:         map[string]*Song{},
		logger:        logger.With("playlist", pe.Name),
		DownloadQueue: []Song{},
	}
	return &playlist
}

// Create dummy songs from files that exist on disk but are untracked
func (p *Playlist) SyncFromDisk() {
	err := os.MkdirAll(p.Folder, os.ModePerm)
	if err != nil {
		p.logger.Error("Failed to create playlist directory", "err", err)
	}
	dir, err := os.ReadDir(p.Folder)
	if err != nil {
		p.logger.Error("Failed to sync from disk", "err", err)
		return
	}

	for _, entry := range dir {
		idCandidates := fileIDRegex.FindStringSubmatch(entry.Name())
		id := idCandidates[len(idCandidates)-1]

		if _, exists := p.Songs[id]; exists {
			continue
		}

		p.logger.Debug("Found untracked file on disk, adding to list", "name", entry.Name(), "id", id)
		newSong := NewSong(id, "", entry.Name(), p.logger)
		p.Songs[id] = &newSong
	}

	p.logger.Info("Synced playlist from disk", "downloadedSongs", len(p.Songs))
}

func (p *Playlist) Update() {
	p.logger.Debug("Updating playlist from remote...")
	newSongs, err := p.fetchPlaylistSongs()
	if err != nil {
		p.logger.Error("Failed to update playlist", "err", err)
		return
	}
	if len(newSongs) == 0 {
		p.logger.Debug("Playlist returned length 0, skipping")
		return
	}

	p.updateTrackedSongs(newSongs)
	p.LastUpdate = time.Now()
}

// Get song list from URL
func (p *Playlist) fetchPlaylistSongs() ([]Song, error) {
	dl := ytdlp.New().
		PrintJSON().
		Simulate()

	result, err := dl.Run(context.TODO(), p.URL)
	if err != nil {
		p.logger.Error("Failed to query playlist", "err", err)
		return nil, err
	}

	info, err := result.GetExtractedInfo()
	if err != nil {
		p.logger.Error("Failed to extract playlist info", "err", err)
		return nil, err
	}

	songs := []Song{}

	for _, song := range info {
		songs = append(songs, NewSongFromInfo(song, p.logger))
	}

	return songs, nil
}

func (p *Playlist) updateTrackedSongs(newSongs []Song) {
	newSongCount, existingSongCount, deleteCountdownSongCount, deletedSongCount := 0, 0, 0, 0

	for _, newSong := range newSongs {
		song, exists := p.Songs[newSong.ID]
		if exists {
			// song is already tracked
			song.ResetDeleteCounter()
			p.Songs[newSong.ID] = &newSong

			existingSongCount++
		} else {
			// song is untracked (new)
			p.DownloadQueue = append(p.DownloadQueue, newSong)
			newSongCount++
		}
	}

	// handle deleted songs
	for _, song := range p.Songs {
		if !slices.ContainsFunc(newSongs, func(a Song) bool {
			return song.ID == a.ID
		}) {
			deleted := song.DeleteAttempt(p.Folder)
			deleteCountdownSongCount++
			if deleted {
				delete(p.Songs, song.ID)
				deletedSongCount++
			}
		}
	}

	p.logger.Debug("Playlist updated", "newSongs", newSongCount, "existingSongs", existingSongCount, "deleteCountdownSongs", deleteCountdownSongCount, "deletedSongs", deletedSongCount)
}

// Download one queue entry if available, return true if download occured
func (p *Playlist) ProcessQueue() bool {
	p.logger.Debug("Processing download queue...")
	defer p.logger.Debug("Finished processing download queue")
	if len(p.DownloadQueue) == 0 {
		p.logger.Debug("Queue empty, skipping")
		return false
	}

	err := p.pruneQueue()
	if err != nil {
		p.logger.Error("Failed to prune queue", "err", err)
		return false
	}
	if len(p.DownloadQueue) == 0 {
		p.logger.Debug("Queue empty, skipping")
		return false
	}

	targetSong := p.DownloadQueue[0]

	p.logger.Info("Downloading song", "id", targetSong.ID, "song", targetSong.Title, "url", targetSong.URL)
	err = p.downloadSong(targetSong)
	if err != nil {
		p.logger.Error("Failed to download song", "err", err)
		return false
	}

	p.logger.Debug("Song downloaded successfuly")

	p.DownloadQueue = p.DownloadQueue[1:]
	return true
}

// Remove already downloaded songs from download queue
func (p *Playlist) pruneQueue() error {
	files, err := os.ReadDir(p.Folder)
	if err != nil {
		return err
	}
	remainingQueue := []Song{}
	for _, queueItem := range p.DownloadQueue {
		if !slices.ContainsFunc(files, func(file os.DirEntry) bool {
			return strings.Contains(file.Name(), queueItem.ID)
		}) {
			remainingQueue = append(remainingQueue, queueItem)
		}
	}
	p.logger.Debug("Pruned queue", "removedEntried", len(p.DownloadQueue)-len(remainingQueue), "remainingEntries", len(remainingQueue))
	p.DownloadQueue = remainingQueue
	return nil
}

func (p *Playlist) downloadSong(song Song) error {
	dl := ytdlp.New().
		AudioQuality("0").
		ExtractAudio().
		EmbedThumbnail().
		AudioFormat("mp3").
		EmbedMetadata().
		AbortOnError().
		NoPart().
		ProgressFunc(time.Millisecond*200, func(update ytdlp.ProgressUpdate) {
			p.logger.Debug("Downloading...", "progress", fmt.Sprintf("%.1f%%", update.Percent()), "ETA", update.ETA())
		}).
		Output(fmt.Sprintf("%s/%s [%s]", p.Folder, strings.ReplaceAll(song.Title, "/", "|"), song.ID))

	_, err := dl.Run(context.TODO(), song.URL)
	if err != nil {
		return err
	}

	return nil
}
