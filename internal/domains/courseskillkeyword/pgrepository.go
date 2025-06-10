package courseskillkeyword

import (
	cm "orientation-training-api/internal/common"
	m "orientation-training-api/internal/models"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
)

type PgCourseSkillKeywordRepository struct {
	cm.AppRepository
}

func NewPgCourseSkillKeywordRepository(logger echo.Logger) (repo *PgCourseSkillKeywordRepository) {
	repo = &PgCourseSkillKeywordRepository{}
	repo.Init(logger)
	return
}

// InsertCourseSkillKeywordWithTx : Insert course skill keyword relation
func (repo *PgCourseSkillKeywordRepository) InsertCourseSkillKeywordWithTx(tx *pg.Tx, courseID int, skillKeywordID int) error {
	courseSkillKeyword := m.CourseSkillKeyword{
		CourseID:       courseID,
		SkillKeywordID: skillKeywordID,
	}

	err := tx.Insert(&courseSkillKeyword)
	if err != nil {
		repo.Logger.Error(err)
	}
	return err
}

func (repo *PgCourseSkillKeywordRepository) DeleteByCourseID(courseID int) error {
	_, err := repo.DB.Model(&m.CourseSkillKeyword{}).
		TableExpr("course_skill_keywords AS csk").
		Where("csk.course_id = ?", courseID).
		Delete()

	if err != nil {
		repo.Logger.Error(err)
	}
	return err
}
