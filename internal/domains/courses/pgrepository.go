package courses

import (
	cm "orientation-training-api/internal/common"
	m "orientation-training-api/internal/models"

	"github.com/labstack/echo/v4"
)

type PgCourseRepository struct {
	cm.AppRepository
}

func NewPgCourseRepository(logger echo.Logger) (repo *PgCourseRepository) {
	repo = &PgCourseRepository{}
	repo.Init(logger)
	return
}

func (repo *PgCourseRepository) GetCourseByID(id int) (m.Course, error) {
	course := m.Course{}
	err := repo.DB.Model(&course).
		Where("id = ?", id).
		Where("deleted_at is null").
		First()

	return course, err
}
