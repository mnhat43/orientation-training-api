package courses

import (
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"

	"github.com/go-pg/pg/v9"
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

func (repo *PgCourseRepository) GetCourses(courseListParams *param.CourseListParams) ([]m.Course, int, error) {
	courses := []m.Course{}
	queryObj := repo.DB.Model(&courses)
	if courseListParams.Keyword != "" {
		queryObj.Where("LOWER(title) LIKE LOWER(?)", "%"+courseListParams.Keyword+"%")
	}
	queryObj.Offset((courseListParams.CurrentPage - 1) * courseListParams.RowPerPage)
	queryObj.Order("created_at DESC")
	queryObj.Limit(courseListParams.RowPerPage)
	totalRow, err := queryObj.SelectAndCount()
	return courses, totalRow, err
}

// GetAllCourses retrieves all courses without pagination
func (repo *PgCourseRepository) GetAllCourses() ([]m.Course, error) {
	courses := []m.Course{}
	query := repo.DB.Model(&courses).
		Where("deleted_at IS NULL").
		Order("created_at DESC")

	err := query.Select()
	if err != nil {
		repo.Logger.Errorf("Error fetching all courses: %v", err)
		return nil, err
	}

	return courses, nil
}

// SaveCourse : insert data to course and associated skill keywords
// Params : createCourseParams, courseSkillKeywordRepo
// Returns : return object of record that 've just been inserted
func (repo *PgCourseRepository) SaveCourse(createCourseParams *param.CreateCourseParams, courseSkillKeywordRepo rp.CourseSkillKeywordRepository) (m.Course, error) {
	course := m.Course{}
	err := repo.DB.RunInTransaction(func(tx *pg.Tx) error {
		var transErr error
		course, transErr = repo.InsertCourseWithTx(
			tx,
			createCourseParams.Title,
			createCourseParams.Description,
			createCourseParams.Thumbnail,
			createCourseParams.Category,
			createCourseParams.CreatedBy,
		)
		if transErr != nil {
			repo.Logger.Error(transErr)
			return transErr
		}

		if len(createCourseParams.SkillKeywordIDs) > 0 {
			for _, skillKeywordID := range createCourseParams.SkillKeywordIDs {
				transErr = courseSkillKeywordRepo.InsertCourseSkillKeywordWithTx(tx, course.ID, skillKeywordID)
				if transErr != nil {
					repo.Logger.Error(transErr)
					return transErr
				}
			}
		}

		return nil
	})

	return course, err
}

// InsertCourseWithTx : insert data to courses
// Params : pg.Tx, title, description, thumbnail
// Returns : return course object , error
func (repo *PgCourseRepository) InsertCourseWithTx(tx *pg.Tx, title string, description string, thumbnail string, category string, createdBy int) (m.Course, error) {
	course := m.Course{
		Title:       title,
		Description: description,
		Thumbnail:   thumbnail,
		Category:    category,
		CreatedBy:   createdBy,
	}
	err := tx.Insert(&course)
	return course, err
}

// UpdateCourse : update course
// Params : CourseID, Title,Thumbnail, Description
// Returns : error
func (repo *PgCourseRepository) UpdateCourse(courseParams *param.UpdateCourseParams, userCourseRepo rp.UserCourseRepository, courseSkillKeywordRepo rp.CourseSkillKeywordRepository) error {
	err := repo.DB.RunInTransaction(func(tx *pg.Tx) error {
		var transErr error

		updateQuery := tx.Model((*m.Course)(nil)).Where("id = ?", courseParams.ID)

		updateQuery = updateQuery.Set("updated_at = NOW()")

		if courseParams.Title != "" {
			updateQuery = updateQuery.Set("title = ?", courseParams.Title)
		}
		if courseParams.Description != "" {
			updateQuery = updateQuery.Set("description = ?", courseParams.Description)
		}
		if courseParams.Thumbnail != "" {
			updateQuery = updateQuery.Set("thumbnail = ?", courseParams.Thumbnail)
		}
		if courseParams.Category != "" {
			updateQuery = updateQuery.Set("category = ?", courseParams.Category)
		}

		if _, transErr = updateQuery.Update(); transErr != nil {
			repo.Logger.Error()
			return transErr
		}
		if courseParams.SkillKeywordIDs != nil {
			repo.Logger.Infof("Updating skill keywords for course %d. New keywords: %v", courseParams.ID, courseParams.SkillKeywordIDs)

			existingKeywords, checkErr := courseSkillKeywordRepo.GetSkillKeywordsByCourseID(courseParams.ID)
			if checkErr == nil {
				repo.Logger.Infof("Course %d currently has %d skill keywords", courseParams.ID, len(existingKeywords))
			}

			repo.Logger.Infof("Force deleting ALL existing skill keywords for course %d", courseParams.ID)
			if transErr = courseSkillKeywordRepo.DeleteByCourseIDWithTx(tx, courseParams.ID); transErr != nil {
				repo.Logger.Errorf("Failed to delete existing skill keywords for course %d: %v", courseParams.ID, transErr)
				return transErr
			}

			repo.Logger.Infof("Inserting %d new skill keywords for course %d", len(courseParams.SkillKeywordIDs), courseParams.ID)
			for i, skillKeywordID := range courseParams.SkillKeywordIDs {
				repo.Logger.Infof("Inserting skill keyword %d/%d: ID %d for course %d", i+1, len(courseParams.SkillKeywordIDs), skillKeywordID, courseParams.ID)
				if transErr = courseSkillKeywordRepo.InsertCourseSkillKeywordWithTx(tx, courseParams.ID, skillKeywordID); transErr != nil {
					repo.Logger.Errorf("Failed to insert skill keyword %d for course %d: %v", skillKeywordID, courseParams.ID, transErr)
					return transErr
				}
			}

			finalKeywords, finalErr := courseSkillKeywordRepo.GetSkillKeywordsByCourseID(courseParams.ID)
			if finalErr == nil {
				repo.Logger.Infof("FINAL RESULT: Course %d now has %d skill keywords in database", courseParams.ID, len(finalKeywords))
				if len(finalKeywords) != len(courseParams.SkillKeywordIDs) {
					repo.Logger.Errorf("MISMATCH: Expected %d skill keywords but found %d in database", len(courseParams.SkillKeywordIDs), len(finalKeywords))
				}
			}
		}

		return nil
	})

	return err
}

// DeleteCourse : delete course by ID
// Params : courseID
// Returns : error
func (repo *PgCourseRepository) DeleteCourse(courseID int) error {
	course := m.Course{}
	_, err := repo.DB.Model(&course).
		Where("id = ?", courseID).
		Delete()

	return err
}

// GetUserCourses retrieves all courses that a specific user is enrolled in,
// sorted by course_position from UserProgress table using relations
func (repo *PgCourseRepository) GetUserCourses(userID int) ([]m.Course, error) {
	var userProgresses []m.UserProgress

	query := repo.DB.Model(&userProgresses).
		Relation("Course").
		Where("user_progress.user_id = ?", userID).
		Where("user_progress.deleted_at IS NULL").
		Where("course.deleted_at IS NULL").
		Order("user_progress.course_position ASC")

	err := query.Select()
	if err != nil {
		repo.Logger.Errorf("Error fetching user courses: %v", err)
		return nil, err
	}

	var courses []m.Course
	for _, progress := range userProgresses {
		if progress.Course != nil {
			courses = append(courses, *progress.Course)
		}
	}

	return courses, nil
}
