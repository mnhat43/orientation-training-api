package modules

import (
	cm "orientation-training-api/internal/common"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"

	"github.com/labstack/echo/v4"
)

type PgModuleRepository struct {
	cm.AppRepository
}

func NewPgModuleRepository(logger echo.Logger) (repo *PgModuleRepository) {
	repo = &PgModuleRepository{}
	repo.Init(logger)
	return
}

func (repo *PgModuleRepository) GetModules(moduleListParams *param.ModuleListParams) ([]m.Module, int, error) {
	modules := []m.Module{}
	queryObj := repo.DB.Model(&modules)
	if moduleListParams.Keyword != "" {
		queryObj.Where("LOWER(title) LIKE LOWER(?)", "%"+moduleListParams.Keyword+"%")
	}
	queryObj.Where("course_id = ?", moduleListParams.CourseID)
	queryObj.Offset((moduleListParams.CurrentPage - 1) * moduleListParams.RowPerPage)
	queryObj.Limit(moduleListParams.RowPerPage)
	queryObj.Order("position ASC") // Added order by position
	totalRow, err := queryObj.SelectAndCount()
	return modules, totalRow, err
}

// SaveModule : insert data into the module table
// Params : createModuleParams contains module creation details (title, course ID, etc.)
// Returns : return the inserted module record or an error
func (repo *PgModuleRepository) SaveModule(createModuleParams *param.CreateModuleParams) (m.Module, error) {
	module := m.Module{
		Title:    createModuleParams.Title,
		CourseID: createModuleParams.CourseID,
		Position: createModuleParams.Position,
	}
	err := repo.DB.Insert(&module)
	return module, err
}

func (repo *PgModuleRepository) GetModuleByID(id int) (m.Module, error) {
	module := m.Module{}
	err := repo.DB.Model(&module).
		Where("id = ?", id).
		Where("deleted_at is null").
		First()

	return module, err
}

// DeleteModule : delete module by ID
// Params : moduleID
// Returns : error
func (repo *PgModuleRepository) DeleteModule(moduleID int) error {
	module := m.Module{}
	_, err := repo.DB.Model(&module).
		Where("id = ?", moduleID).
		Delete()

	return err
}

func (repo *PgModuleRepository) GetModuleIDsByCourseID(courseID int) ([]int, error) {
	moduleIDs := []int{}
	err := repo.DB.Model((*m.Module)(nil)).
		Column("id").
		Where("course_id = ?", courseID).
		Order("position ASC").
		Select(&moduleIDs)
	if err != nil {
		return nil, err
	}
	return moduleIDs, nil
}

// GetModulesByCourseID : retrieve all modules for a specific course
// Params : courseID
// Returns : slice of modules and error
func (repo *PgModuleRepository) GetModulesByCourseID(courseID int) ([]m.Module, error) {
	modules := []m.Module{}
	err := repo.DB.Model(&modules).
		Where("course_id = ?", courseID).
		Where("deleted_at is null").
		Order("position ASC").
		Select()
	if err != nil {
		return nil, err
	}
	return modules, nil
}

// GetModuleByPositionAndCourse retrieves a module by its course ID and position
// Params: courseID, position
// GetModuleByPositionAndCourse returns a module with a specific position in a course
// Returns: module and error
func (repo *PgModuleRepository) GetModuleByPositionAndCourse(courseID int, position int) (m.Module, error) {
	module := m.Module{}
	// Using a relation would require defining the relation in the Module model first
	err := repo.DB.Model(&module).
		Where("course_id = ?", courseID).
		Where("position = ?", position).
		Where("deleted_at is null").
		First()

	return module, err
}
