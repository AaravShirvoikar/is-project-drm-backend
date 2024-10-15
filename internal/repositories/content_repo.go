package repositories

import (
	"database/sql"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
)

type ContentRepository interface {
	Create(content *models.Content) error
	List() ([]*models.Content, error)
	GetById(id string) (*models.Content, error)
}

type contentRepo struct {
	db *sql.DB
}

func NewContentRepository(db *sql.DB) ContentRepository {
	return &contentRepo{db: db}
}

func (r *contentRepo) Create(content *models.Content) error {
	query := `INSERT INTO content (id, title, description, creator_id, price, created_at, updated_at, file_id, file_size)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query, content.ContentID, content.Title, content.Description, content.CreatorID,
		content.Price, content.CreatedAt, content.UpdatedAt, content.FileID, content.FileSize)
	if err != nil {
		return err
	}

	return nil
}

func (r *contentRepo) List() ([]*models.Content, error) {
	query := "SELECT id, title, description, creator_id, price, created_at, updated_at, file_id, file_size FROM content"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contents []*models.Content
	for rows.Next() {
		var content models.Content
		err := rows.Scan(&content.ContentID, &content.Title, &content.Description, &content.CreatorID,
			&content.Price, &content.CreatedAt, &content.UpdatedAt, &content.FileID, &content.FileSize)
		if err != nil {
			return nil, err
		}
		contents = append(contents, &content)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contents, nil
}

func (r *contentRepo) GetById(id string) (*models.Content, error) {
	query := `SELECT id, title, description, creator_id, price, created_at, updated_at, file_id, file_size
              FROM content WHERE id = $1`

	var content models.Content
	err := r.db.QueryRow(query, id).Scan(&content.ContentID, &content.Title, &content.Description,
		&content.CreatorID, &content.Price, &content.CreatedAt, &content.UpdatedAt, &content.FileID,
		&content.FileSize)
	if err != nil {
		return nil, err
	}

	return &content, nil
}
