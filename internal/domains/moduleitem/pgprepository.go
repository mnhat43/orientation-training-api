package moduleitem

import (
	cm "orientation-training-api/internal/common"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"

	"github.com/go-pg/pg/v9"
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
	queryObj.Where("module_id = ?", moduleItemListParams.ModuleID)
	queryObj.Order("position ASC")
	totalRow, err := queryObj.SelectAndCount()
	return moduleItems, totalRow, err
}

// SaveModuleItem : insert data to item
// Params : param.CreateCourseParams
// Returns : return object of record that 've just been inserted
func (repo *PgModuleItemRepository) SaveModuleItem(createModuleItemParams *param.CreateModuleItemParams) (m.ModuleItem, error) {
	moduleItem := m.ModuleItem{
		Title:        createModuleItemParams.Title,
		ItemType:     createModuleItemParams.ItemType,
		Resource:     createModuleItemParams.Resource,
		Position:     createModuleItemParams.Position,
		RequiredTime: createModuleItemParams.RequiredTime,
		ModuleID:     createModuleItemParams.ModuleID,
	}

	err := repo.DB.Insert(&moduleItem)
	return moduleItem, err
}

func (repo *PgModuleItemRepository) GetModuleItemByID(id int) (m.ModuleItem, error) {
	moduleItem := m.ModuleItem{}
	err := repo.DB.Model(&moduleItem).
		Where("id = ?", id).
		Where("deleted_at is null").
		First()

	return moduleItem, err
}

func (repo *PgModuleItemRepository) DeleteModuleItem(moduleItemID int) error {
	moduleItem := m.ModuleItem{}
	_, err := repo.DB.Model(&moduleItem).
		Where("id = ?", moduleItemID).
		Delete()

	return err
}

func (repo *PgModuleItemRepository) GetModuleItemsByModuleIDs(moduleIDs []int) ([]m.ModuleItem, error) {
	moduleItems := []m.ModuleItem{}
	err := repo.DB.Model(&moduleItems).
		Where("module_id IN (?)", pg.In(moduleIDs)).
		Where("deleted_at is null").
		Order("position ASC").
		Select()
	if err != nil {
		return nil, err
	}
	return moduleItems, nil
}
