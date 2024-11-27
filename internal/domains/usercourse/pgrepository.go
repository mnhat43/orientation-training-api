package usercourse

import (
	cm "orientation-training-api/internal/common"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
)

type PgUserCourseRepository struct {
	cm.AppRepository
}

func NewPgUserCourseRepository(logger echo.Logger) (repo *PgUserCourseRepository) {
	repo = &PgUserCourseRepository{}
	repo.Init(logger)
	return
}

func (repo *PgUserCourseRepository) DeleteByCourseId(courseId int) error {
	_, err := repo.DB.Model(&m.UserCourse{}).
		TableExpr("user_courses AS uc").
		Where("uc.course_id = ?", courseId).
		Delete()

	if err != nil {
		repo.Logger.Error(err)
	}

	return err
}

// InsertUserCourse : Insert user course
func (repo *PgUserCourseRepository) InsertUserCourseWithTx(tx *pg.Tx, userID int, courseID int) error {
	userCourse := m.UserCourse{
		UserID:   userID,
		CourseID: courseID,
	}

	err := tx.Insert(&userCourse)
	if err != nil {
		repo.Logger.Error(err)
	}

	return err
}

func (repo *PgUserCourseRepository) SelectMembersInCourse(courseId int) ([]param.UserCourseInfoRecords, error) {
	var users []param.UserCourseInfoRecords
	err := repo.DB.Model(&m.UserCourse{}).
		Column("user_course.user_id").
		ColumnExpr("up.first_name || ' ' || up.last_name full_name").
		Join("JOIN courses AS p ON p.id = user_course.course_id").
		Join("JOIN user_profiles AS up ON up.user_id = user_course.user_id").
		Where("user_course.course_id = ?", courseId).
		Select(&users)

	if err != nil {
		repo.Logger.Error(err)
	}

	return users, err
}
