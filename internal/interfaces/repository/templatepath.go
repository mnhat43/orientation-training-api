package repository

import (
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
)

// TemplatePathRepository defines methods for accessing template path data
type TemplatePathRepository interface {
	GetTemplatePathByID(TempPathID int) (m.TemplatePath, error)
	GetTemplatePathList() ([]m.TemplatePath, error)
	CreateTemplatePath(createTemplatePathParams *param.CreateTemplatePathParams) (m.TemplatePath, error)
	UpdateTemplatePath(updateTemplatePathParams *param.UpdateTemplatePathParams) (m.TemplatePath, error)
	DeleteTemplatePath(TempPathID int) error
}
