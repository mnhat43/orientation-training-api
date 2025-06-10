package repository

import (
	"github.com/go-pg/pg/v9"
)

type CourseSkillKeywordRepository interface {
	InsertCourseSkillKeywordWithTx(tx *pg.Tx, courseID int, skillKeywordID int) error
	DeleteByCourseID(courseID int) error
}
