package repository

import (
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
)

type ModuleItemRepository interface {
	GetModuleItems(moduleItemListParams *param.ModuleItemListParams) ([]m.ModuleItem, int, error)
	SaveModuleItem(moduleItemListParams *param.CreateModuleItemParams) (m.ModuleItem, error)
	GetModuleItemByID(id int) (m.ModuleItem, error)
	DeleteModuleItem(moduleItemID int) error
	GetModuleItemsByModuleIDs(moduleIDs []int, itemType string) ([]m.ModuleItem, error)
}
