package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/fevse/songlib/internal/storage"
)

// CreateSong godoc
// @Summary Добавление новой песни
// @Description Добавление новой песни с добавлением данных от стороннего API
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body storage.Song true "Song to add"
// @Success 201 {object} storage.Song
// @Failure 400 {string} string "Invalid JSON"
// @Failure 500 {string} string "Failed to create song"
// @Router /songs [post]
func (s *Server) CreateSong() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var song storage.Song
		if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := s.app.CreateSong(&song); err != nil {
			log.Printf("Error creating song: %v", err)
			http.Error(w, "Failed to create song", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(song)
	}
}

// UpdateSong godoc
// @Summary Обновление песни в библиотеке
// @Description Обновляет информацию о песне, данные передаются в теле запроса в формате JSON
// @Tags songs
// @Produce  json
// @Param id path int true "Song ID"
// @Param song body storage.Song true "Updated song details"
// @Success 200 {object} storage.Song
// @Failure 400 {string} string "Invalid ID or JSON"
// @Failure 500 {string} string "Failed to update song"
// @Router /songs/{id} [put]
func (s *Server) UpdateSong() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Printf("Error converting id to int: %v", err)
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		var song storage.Song
		if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		song.ID = id
		if err := s.app.UpdateSong(&song); err != nil {
			log.Printf("Error updating song: %v", err)
			http.Error(w, "Failed to update song", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(song)
	}
}

// DeleteSong godoc
// @Summary Удаление песни из библиотеки
// @Description Удаляет песню из библиотеки по ID
// @Tags songs
// @Param id path int true "Song ID"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 500 {string} string "Failed to delete song"
// @Router /songs/{id} [delete]
func (s *Server) DeleteSong() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			log.Printf("Error converting id to int: %v", err)
			http.Error(w, "Invalid ID", http.StatusBadRequest)
			return
		}

		if err := s.app.DeleteSong(id); err != nil {
			log.Printf("Error deleting song: %v", err)
			http.Error(w, "Failed to delete song", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// GetAllSongs godoc
// @Summary Получение песни или списка песен
// @Description Получение песни или списка песен: limit - количество выводимых данных, offset - с какого элемента, можно фильтровать по любым полям
// @Tags songs
// @Produce  json
// @Param group query string false "Filter by group"
// @Param song query string false "Filter by song name"
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} storage.Song
// @Failure 500 {string} string "Failed to get songs"
// @Router /songs [get]
func (s *Server) GetSongs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := make(map[string]string)
		for key, values := range r.URL.Query() {
			if key != "limit" && key != "offset" {
				filter[key] = values[0]
			}
		}

		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))

		songs, err := s.app.GetSongs(filter, limit, offset)
		if err != nil {
			log.Printf("Error getting songs: %v", err)
			http.Error(w, "Failed to get songs", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(songs)
	}
}
