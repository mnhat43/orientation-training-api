package repository

import (
	param "orientation-training-api/internal/interfaces/requestparams"

	"github.com/go-pg/pg/v9"
)

type UserCourseRepository interface {
	DeleteByCourseId(courseId int) error
	InsertUserCourseWithTx(tx *pg.Tx, userID int, courseID int) error
	SelectMembersInCourse(courseId int) ([]param.UserCourseInfoRecords, error)
}
