package repository

import (
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
)

type ModuleItemRepository interface {
	GetModuleItems(moduleItemListParams *param.ModuleItemListParams) ([]m.ModuleItem, int, error)
	// GetModuleItemDetail(moduleListParams *param.ModuleItemListParams) ([]m.ModuleItemDetail, int, error)
	// SaveModuleItem(createModuleItemDBParams *param.CreateModuleItemDBParams, userModuleItemRepo UserModuleItemRepository) (m.ModuleItem, error)
	// DeleteModuleItem(moduleID int) error
}
