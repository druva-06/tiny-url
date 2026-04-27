package repository

import (
	"log"

	"github.com/druva-06/tiny-url/internal/db"
)

type URLRepository struct {
	cluster *db.DBCluster
}

func NewURLRepository(cluster *db.DBCluster) *URLRepository {
	return &URLRepository{cluster: cluster}
}

func (r *URLRepository) Create(longUrl string) (int64, error) {
	query := `INSERT INTO url_mapping (long_url) VALUES (?)`
	result, err := r.cluster.Exec(query, longUrl)
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
	_, err = r.cluster.Exec(query, shortCode, id)
	return
}

func (r *URLRepository) GetLongURL(shortCode string) (longURL string, err error) {
	query := `SELECT long_url FROM url_mapping WHERE short_code = ?`
	log.Printf("[URLRepository] QUERY code=%s", shortCode)
	rows, err := r.cluster.Query(query, shortCode)
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
	query := `UPDATE url_mapping SET long_url = ? WHERE short_code = ?`
	log.Printf("[URLRepository] QUERY code=%s long_url=%s query=%s", shortCode, longUrl, query)
	_, err = r.cluster.Exec(query, longUrl, shortCode)
	return
}
