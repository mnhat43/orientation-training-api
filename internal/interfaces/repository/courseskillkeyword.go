package repository

import (
	m "orientation-training-api/internal/models"

	"github.com/go-pg/pg/v9"
)

type CourseSkillKeywordRepository interface {
	InsertCourseSkillKeywordWithTx(tx *pg.Tx, courseID int, skillKeywordID int) error
	DeleteByCourseID(courseID int) error
	DeleteByCourseIDWithTx(tx *pg.Tx, courseID int) error
	GetSkillKeywordsByCourseID(courseID int) ([]m.SkillKeyword, error)
}
