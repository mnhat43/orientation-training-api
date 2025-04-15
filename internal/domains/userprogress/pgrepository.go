package userprogress

import (
	cm "orientation-training-api/internal/common"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"

	"github.com/labstack/echo/v4"
)

type PgUserProgressRepository struct {
	cm.AppRepository
}

func NewPgUserProgressRepository(logger echo.Logger)(repo *PgUserProgressRepository) {
	repo = &PgUserProgressRepository{}
	repo.Init(logger)
	return
}

func (repo *PgUserProgressRepository) GetUserProgress(getUserProgressParams *param.GetUserProgressParams) (m.UserProgress, error) {
	userProgress := m.UserProgress{}

	err := repo.DB.Model(&userProgress).Where("user_id = ?", getUserProgressParams.UserID).
		Where("course_id = ?", getUserProgressParams.CourseID).
		Where("deleted_at IS NULL").
		First()

	return userProgress, err
}

// SaveUserProgress creates a new progress record or updates an existing one
func (repo *PgUserProgressRepository) SaveUserProgress(userProgress *m.UserProgress) error {
	// Check if record exists
	exists, err := repo.DB.Model((*m.UserProgress)(nil)).
		Where("user_id = ?", userProgress.UserID).
		Where("course_id = ?", userProgress.CourseID).
		Where("deleted_at IS NULL").
		Exists()
	
	if err != nil {
		return err
	}
	
	if exists {
		// Update existing record
		_, err = repo.DB.Model(userProgress).
			Set("module_position = ?", userProgress.ModulePosition).
			Set("module_item_position = ?", userProgress.ModuleItemPosition).
			Set("completed = ?", userProgress.Completed).
			Set("updated_at = NOW()").
			Where("user_id = ?", userProgress.UserID).
			Where("course_id = ?", userProgress.CourseID).
			Where("deleted_at IS NULL").
			Update()
	} else {
		// Create new record
		_, err = repo.DB.Model(userProgress).Insert()
	}
	
	return err
}
