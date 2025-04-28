package courses

import (
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
	"orientation-training-api/internal/platform/utils"

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

// SaveCourse : insert data to course
// Params : orgID, param.CreateCourseParams
// Returns : return object of record that 've just been inserted
func (repo *PgCourseRepository) SaveCourse(createCourseParams *param.CreateCourseParams, userCourseRepo rp.UserCourseRepository) (m.Course, error) {
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

		transErr = userCourseRepo.InsertUserCourseWithTx(tx, createCourseParams.CreatedBy, course.ID)
		if transErr != nil {
			repo.Logger.Error(transErr)
			return transErr
		}
		return transErr
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
func (repo *PgCourseRepository) UpdateCourse(courseParams *param.UpdateCourseParams, userCourseRepo rp.UserCourseRepository) error {
	currentCourse, err := repo.GetCourseByID(courseParams.ID)
	if err != nil {
		repo.Logger.Error()
		return err
	}

	course := &m.Course{
		Title:       courseParams.Title,
		Description: courseParams.Description,
		Thumbnail:   courseParams.Thumbnail,
		Category:    courseParams.Category,
		CreatedBy:   courseParams.CreatedBy,
	}

	users, transErr := userCourseRepo.SelectMembersInCourse(courseParams.ID)
	if transErr != nil && transErr.Error() != pg.ErrNoRows.Error() {
		repo.Logger.Error()
		return transErr
	}

	var usersId []int
	if len(users) > 0 {
		for _, user := range users {
			usersId = append(usersId, user.UserId)
		}
	}

	if courseParams.CreatedBy != currentCourse.CreatedBy && !utils.FindIntInSlice(usersId, courseParams.CreatedBy) {
		err = repo.DB.RunInTransaction(func(tx *pg.Tx) error {
			var transErr error
			if _, transErr = tx.Model(course).
				Column("title", "description", "thumbnail", "category", "created_by", "updated_at").
				Where("id = ?", courseParams.ID).
				Update(); transErr != nil {
				repo.Logger.Error()
				return transErr
			}

			if transErr = userCourseRepo.InsertUserCourseWithTx(tx, courseParams.CreatedBy, courseParams.ID); transErr != nil {
				repo.Logger.Error()
				return transErr
			}

			return transErr
		})
	} else {
		_, err = repo.DB.Model(course).
			Column("title", "description", "thumbnail", "category", "created_by", "updated_at").
			Where("id = ?", courseParams.ID).
			Update()
	}

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
