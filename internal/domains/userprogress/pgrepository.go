package userprogress

import (
	cm "orientation-training-api/internal/common"
	m "orientation-training-api/internal/models"

	"github.com/labstack/echo/v4"
)

type PgUserProgressRepository struct {
	cm.AppRepository
}

func NewPgUserProgressRepository(logger echo.Logger) (repo *PgUserProgressRepository) {
	repo = &PgUserProgressRepository{}
	repo.Init(logger)
	return
}

func (repo *PgUserProgressRepository) GetSingleUserProgress(userID int, courseID int) (m.UserProgress, error) {
	userProgress := m.UserProgress{}

	query := repo.DB.Model(&userProgress).
		Where("user_id = ?", userID).
		Where("course_id = ?", courseID).
		Where("deleted_at IS NULL")

	err := query.First()

	return userProgress, err
}

// SaveUserProgress creates a new progress record or updates an existing one
func (repo *PgUserProgressRepository) SaveUserProgress(userProgress *m.UserProgress) error {
	exists, err := repo.DB.Model((*m.UserProgress)(nil)).
		Where("user_id = ?", userProgress.UserID).
		Where("course_id = ?", userProgress.CourseID).
		Where("deleted_at IS NULL").
		Exists()

	if err != nil {
		repo.Logger.Errorf("Error checking if user progress exists: %v", err)
		return err
	}

	if exists {
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
		_, err = repo.DB.Model(userProgress).Insert()
	}

	return err
}

// GetUserProgressByCourseID retrieves all user progress records for a specific course
func (repo *PgUserProgressRepository) GetUserProgressByCourseID(courseID int) ([]m.UserProgress, error) {
	var userProgressList []m.UserProgress

	err := repo.DB.Model(&userProgressList).
		Where("course_id = ?", courseID).
		Where("deleted_at IS NULL").
		Select()

	if err != nil {
		return nil, err
	}

	return userProgressList, nil
}

// GetAllUserProgressByUserID retrieves all user progress records for a specific user
func (repo *PgUserProgressRepository) GetAllUserProgressByUserID(userID int) ([]m.UserProgress, error) {
	var userProgressList []m.UserProgress

	err := repo.DB.Model(&userProgressList).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Order("course_position ASC").
		Select()

	if err != nil {
		return nil, err
	}

	return userProgressList, nil
}
