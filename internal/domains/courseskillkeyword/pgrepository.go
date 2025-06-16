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
		repo.Logger.Errorf("Error inserting course skill keyword (course: %d, skill: %d): %v", courseID, skillKeywordID, err)
	} else {
		repo.Logger.Infof("Successfully inserted course skill keyword (course: %d, skill: %d)", courseID, skillKeywordID)
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

func (repo *PgCourseSkillKeywordRepository) DeleteByCourseIDWithTx(tx *pg.Tx, courseID int) error {
	var count int
	_, err := tx.Query(&count, "SELECT COUNT(*) FROM course_skill_keywords WHERE course_id = ?", courseID)
	if err != nil {
		repo.Logger.Errorf("Error counting existing skill keywords for course %d: %v", courseID, err)
		return err
	}
	repo.Logger.Infof("Found %d existing skill keywords for course %d", count, courseID)

	result, err := tx.Exec("DELETE FROM course_skill_keywords WHERE course_id = ?", courseID)
	if err != nil {
		repo.Logger.Errorf("Error deleting course skill keywords for course %d: %v", courseID, err)
		return err
	}

	rowsAffected := result.RowsAffected()
	repo.Logger.Infof("Successfully deleted %d skill keyword relations for course %d", rowsAffected, courseID)

	var remainingCount int
	_, err = tx.Query(&remainingCount, "SELECT COUNT(*) FROM course_skill_keywords WHERE course_id = ?", courseID)
	if err != nil {
		repo.Logger.Errorf("Error verifying deletion for course %d: %v", courseID, err)
		return err
	}

	if remainingCount > 0 {
		repo.Logger.Errorf("WARNING: %d skill keywords still remain for course %d after deletion", remainingCount, courseID)
	} else {
		repo.Logger.Infof("Verified: All skill keywords successfully deleted for course %d", courseID)
	}

	return nil
}

// GetSkillKeywordsByCourseID : Get all skill keywords associated with a course
func (repo *PgCourseSkillKeywordRepository) GetSkillKeywordsByCourseID(courseID int) ([]m.SkillKeyword, error) {
	var skillKeywords []m.SkillKeyword

	_, err := repo.DB.Query(&skillKeywords, `
		SELECT sk.* 
		FROM skill_keywords AS sk
		JOIN course_skill_keywords AS csk ON sk.id = csk.skill_keyword_id
		WHERE csk.course_id = ? 
		AND sk.deleted_at IS NULL
		AND csk.deleted_at IS NULL
		ORDER BY sk.id ASC
	`, courseID)

	if err != nil {
		repo.Logger.Errorf("Failed to get skill keywords for course ID %d: %v", courseID, err)
	}

	return skillKeywords, err
}
