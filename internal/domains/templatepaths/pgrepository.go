package templatepaths

import (
	cm "orientation-training-api/internal/common"
	param "orientation-training-api/internal/interfaces/requestparams"
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
func (repo *PgTemplatePathRepository) GetTemplatePathByID(tempPathID int) (m.TemplatePath, error) {
	tempPath := m.TemplatePath{}

	err := repo.DB.Model(&tempPath).
		Where("id = ?", tempPathID).
		Where("deleted_at IS NULL").
		First()

	return tempPath, err
}

// GetTemplatePathList retrieves all active template paths
func (repo *PgTemplatePathRepository) GetTemplatePathList() ([]m.TemplatePath, error) {
	var tempPaths []m.TemplatePath

	err := repo.DB.Model(&tempPaths).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
		Select()

	if err != nil {
		repo.Logger.Errorf("Error retrieving template paths: %v", err)
		return nil, err
	}

	return tempPaths, nil
}

// CreateTemplatePath creates a new template path
func (repo *PgTemplatePathRepository) CreateTemplatePath(createTemplatePathParams *param.CreateTemplatePathParams) (m.TemplatePath, error) {
	templatePath := m.TemplatePath{
		Name:        createTemplatePathParams.Name,
		Description: createTemplatePathParams.Description,
		CourseIds:   createTemplatePathParams.CourseIds,
		Duration:    createTemplatePathParams.Duration,
	}

	err := repo.DB.Insert(&templatePath)
	return templatePath, err
}

// UpdateTemplatePath updates an existing template path
func (repo *PgTemplatePathRepository) UpdateTemplatePath(updateTemplatePathParams *param.UpdateTemplatePathParams) (m.TemplatePath, error) {
	tempPath := m.TemplatePath{
		ID:          updateTemplatePathParams.TempPathID,
		Name:        updateTemplatePathParams.Name,
		Description: updateTemplatePathParams.Description,
		CourseIds:   updateTemplatePathParams.CourseIds,
		Duration:    updateTemplatePathParams.Duration,
	}

	_, err := repo.DB.Model(&tempPath).
		Where("id = ?", tempPath.ID).
		Where("deleted_at IS NULL").
		Column("name", "description", "course_ids", "duration").
		Update()

	return tempPath, err
}

// DeleteTemplatePath soft deletes a template path
func (repo *PgTemplatePathRepository) DeleteTemplatePath(tempPathID int) error {
	_, err := repo.DB.Model((*m.TemplatePath)(nil)).
		Set("deleted_at = NOW()").
		Where("id = ?", tempPathID).
		Where("deleted_at IS NULL").
		Update()
	return err
}
