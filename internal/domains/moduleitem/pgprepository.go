package moduleitem

import (
	cm "orientation-training-api/internal/common"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"

	"github.com/labstack/echo/v4"
)

type PgModuleItemRepository struct {
	cm.AppRepository
}

func NewPgModuleItemRepository(logger echo.Logger) (repo *PgModuleItemRepository) {
	repo = &PgModuleItemRepository{}
	repo.Init(logger)
	return
}

func (repo *PgModuleItemRepository) GetModuleItems(moduleItemListParams *param.ModuleItemListParams) ([]m.ModuleItem, int, error) {
	moduleItems := []m.ModuleItem{}
	queryObj := repo.DB.Model(&moduleItems)
	if moduleItemListParams.Keyword != "" {
		queryObj.Where("LOWER(title) LIKE LOWER(?)", "%"+moduleItemListParams.Keyword+"%")
	}
	queryObj.Where("module_id = ?", moduleItemListParams.ModuleID)
	queryObj.Offset((moduleItemListParams.CurrentPage - 1) * moduleItemListParams.RowPerPage)
	queryObj.Order("created_at DESC")
	queryObj.Limit(moduleItemListParams.RowPerPage)
	totalRow, err := queryObj.SelectAndCount()
	return moduleItems, totalRow, err
}
