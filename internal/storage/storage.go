package storage

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/pressly/goose"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) Migrate() error {
	if err := goose.SetDialect("postgres"); err != nil {
		log.Printf("Error migration settings: %v", err)
		return err
	}
	if err := goose.Up(s.db, "migrations"); err != nil {
		log.Printf("Error migrations: %v", err)
		return err
	}
	return nil
}

func (r *Storage) Create(song *Song) error {
	query := `
		INSERT INTO songs (band, song, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	err := r.db.QueryRow(
		query,
		song.Group, song.Song, song.ReleaseDate,
		song.Text, song.Link).Scan(&song.ID)
	if err != nil {
		log.Printf("Error creating song: %v", err)
		return err
	}
	return nil
}

func (r *Storage) GetByID(id int) (*Song, error) {
	query := `SELECT * FROM songs WHERE id = $1`
	row := r.db.QueryRow(query, id)

	var song Song
	err := row.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		log.Printf("Error getting song: %v", err)
		return nil, err
	}
	return &song, nil
}

func (r *Storage) Update(song *Song) error {
	query := `
		UPDATE songs
		SET band = $1, song = $2, release_date = $3, text = $4, link = $5
		WHERE id = $6`
	_, err := r.db.Exec(query, song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, song.ID)
	if err != nil {
		log.Printf("Error updating song: %v", err)
		return err
	}
	return nil
}

func (r *Storage) Delete(id int) error {
	query := `DELETE FROM songs WHERE id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting song: %v", err)
		return err
	}
	return nil
}

func (r *Storage) GetList(filter map[string]string, limit, offset int) ([]Song, error) {
	query := `SELECT * FROM songs WHERE 1=1`
	args := []any{}
	counter := 1

	for key, value := range filter {
		if key == "group" {
			key = "band"
		}
		query += " AND " + key + " = $" + strconv.Itoa(counter)
		args = append(args, value)
		counter++
	}

	query += " LIMIT $" + strconv.Itoa(counter) + " OFFSET $" + strconv.Itoa(counter+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Printf("Error getting songs: %v", err)
		return nil, err
	}
	defer rows.Close()

	var songs []Song
	for rows.Next() {
		var song Song
		err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
		if err != nil {
			log.Printf("Error scanning song: %v", err)
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}
