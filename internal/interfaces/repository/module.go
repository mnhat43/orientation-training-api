package repository

import (
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
)

type ModuleRepository interface {
	GetModules(moduleListParams *param.ModuleListParams) ([]m.Module, int, error)
	SaveModule(createModuleParams *param.CreateModuleParams) (m.Module, error)
	GetModuleByID(id int) (m.Module, error)
	DeleteModule(moduleID int) error

	// GetModuleDetail(moduleListParams *param.ModuleListParams) ([]m.ModuleDetail, int, error)
	// DeleteModule(moduleID int) error
}
