package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fevse/songlib/internal/storage"
)

type SongLibApp struct {
	storage *storage.Storage
	miURL   string
}

func NewSongLibApp(stor *storage.Storage, miURL string) *SongLibApp {
	return &SongLibApp{storage: stor, miURL: miURL}
}

func (s *SongLibApp) CreateSong(song *storage.Song) error {
	detail, err := s.getSongDetails(song.Group, song.Song)
	if err != nil {
		log.Printf("Error a song details: %v", err)
		// return err
	}

	song.ReleaseDate = detail.ReleaseDate
	song.Text = detail.Text
	song.Link = detail.Link

	if err := s.storage.Create(song); err != nil {
		log.Printf("Error creating song: %v", err)
		return err
	}

	return nil
}

func (s *SongLibApp) UpdateSong(song *storage.Song) error {
	return s.storage.Update(song)
}

func (s *SongLibApp) DeleteSong(id int) error {
	return s.storage.Delete(id)
}

func (s *SongLibApp) GetSongs(filter map[string]string, limit, offset int) ([]storage.Song, error) {
	return s.storage.GetList(filter, limit, offset)
}

func (s *SongLibApp) getSongDetails(group, song string) (*storage.SongDetail, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", s.miURL, group, song)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error getting song details: %v", err)
		return &storage.SongDetail{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status is not OK: %d", resp.StatusCode)
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var detail storage.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		log.Printf("Error decoding API response: %v", err)
		return nil, err
	}

	return &detail, nil
}
