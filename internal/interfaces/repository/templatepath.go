package repository

import (
	m "orientation-training-api/internal/models"
)

// TemplatePathRepository defines methods for accessing template path data
type TemplatePathRepository interface {
	GetTemplatePathByID(pathID int) (m.TemplatePath, error)
	GetTemplatePathList() ([]m.TemplatePath, error)
	CreateTemplatePath(path *m.TemplatePath) error
	UpdateTemplatePath(path *m.TemplatePath) error
	DeleteTemplatePath(pathID int) error
}
