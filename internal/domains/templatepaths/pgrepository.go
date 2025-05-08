package templatepaths

import (
	cm "orientation-training-api/internal/common"
	m "orientation-training-api/internal/models"

	"github.com/labstack/echo/v4"
)

type PgTemplatePathRepository struct {
	cm.AppRepository
}

func NewPgTemplatePathRepository(logger echo.Logger) (repo *PgTemplatePathRepository) {
	repo = &PgTemplatePathRepository{}
	repo.Init(logger)
	return
}

// GetTemplatePathByID retrieves a template path by its ID
func (repo *PgTemplatePathRepository) GetTemplatePathByID(pathID int) (m.TemplatePath, error) {
	path := m.TemplatePath{}

	err := repo.DB.Model(&path).
		Where("id = ?", pathID).
		Where("deleted_at IS NULL").
		First()

	return path, err
}

// GetTemplatePathList retrieves all active template paths
func (repo *PgTemplatePathRepository) GetTemplatePathList() ([]m.TemplatePath, error) {
	var paths []m.TemplatePath

	err := repo.DB.Model(&paths).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Select()

	return paths, err
}

// CreateTemplatePath creates a new template path
func (repo *PgTemplatePathRepository) CreateTemplatePath(path *m.TemplatePath) error {
	_, err := repo.DB.Model(path).Insert()
	return err
}

// UpdateTemplatePath updates an existing template path
func (repo *PgTemplatePathRepository) UpdateTemplatePath(path *m.TemplatePath) error {
	_, err := repo.DB.Model(path).
		Where("id = ?", path.ID).
		Where("deleted_at IS NULL").
		Update()
	return err
}

// DeleteTemplatePath soft deletes a template path
func (repo *PgTemplatePathRepository) DeleteTemplatePath(pathID int) error {
	_, err := repo.DB.Model((*m.TemplatePath)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", pathID).
		Where("deleted_at IS NULL").
		Update()
	return err
}
