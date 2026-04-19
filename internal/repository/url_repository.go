package repository

import (
	"database/sql"
	"log"
)

type URLRepository struct {
	db *sql.DB
}

func NewURLRepository(db *sql.DB) *URLRepository {
	return &URLRepository{db: db}
}

func (r *URLRepository) Create(longUrl string) (int64, error) {
	query := `INSERT INTO url_mapping (long_url) VALUES (?)`
	result, err := r.db.Exec(query, longUrl)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *URLRepository) UpdateShortCode(id int64, shortCode string) (err error) {
	query := `UPDATE url_mapping SET short_code = ? WHERE id = ?`
	_, err = r.db.Exec(query, shortCode, id)
	return
}

func (r *URLRepository) GetLongURL(shortCode string) (longURL string, err error) {
	query := `SELECT long_url FROM url_mapping WHERE short_code = ?`
	log.Printf("[URLRepository] QUERY code=%s", shortCode)
	err = r.db.QueryRow(query, shortCode).Scan(&longURL)
	return
}

func (r *URLRepository) UpdateLongUrl(shortCode string, longUrl string) (err error) {
	query := `UPDATE url_mapping SET long_url = ? WHERE short_code = ?`
	log.Printf("[URLRepository] QUERY code=%s long_url=%s query=%s", shortCode, longUrl, query)
	_, err = r.db.Exec(query, longUrl, shortCode)
	return
}
