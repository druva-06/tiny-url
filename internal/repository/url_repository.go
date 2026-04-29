package repository

import (
	"log"

	"github.com/druva-06/tiny-url/internal/db"
)

type URLRepository struct {
	shard *db.ShardManager
}

func NewURLRepository(shard *db.ShardManager) *URLRepository {
	return &URLRepository{shard: shard}
}

func (r *URLRepository) Create(shortCode, longUrl string) (string, error) {
	db := r.shard.GetDB(shortCode)
	query := `INSERT INTO url_mapping (short_code, long_url) VALUES (?,?)`
	_, err := db.Exec(query, shortCode, longUrl)
	if err != nil {
		return "", err
	}
	return shortCode, nil
}

func (r *URLRepository) UpdateShortCode(id int64, shortCode string) (err error) {
	db := r.shard.GetDB(shortCode)
	query := `UPDATE url_mapping SET short_code = ? WHERE id = ?`
	_, err = db.Exec(query, shortCode, id)
	return
}

func (r *URLRepository) GetLongURL(shortCode string) (longURL string, err error) {
	db := r.shard.GetDB(shortCode)
	query := `SELECT long_url FROM url_mapping WHERE short_code = ?`
	log.Printf("[URLRepository] QUERY code=%s", shortCode)
	rows, err := db.Query(query, shortCode)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&longURL)
		if err != nil {
			return
		}
	}
	return
}

func (r *URLRepository) UpdateLongUrl(shortCode string, longUrl string) (err error) {
	db := r.shard.GetDB(shortCode)
	query := `UPDATE url_mapping SET long_url = ? WHERE short_code = ?`
	log.Printf("[URLRepository] QUERY code=%s long_url=%s query=%s", shortCode, longUrl, query)
	_, err = db.Exec(query, longUrl, shortCode)
	return
}
