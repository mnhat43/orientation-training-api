package courses

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-pg/pg"
	"github.com/labstack/echo/v4"
)

type CourseController struct {
	cm.BaseController

	CourseRepo     rp.CourseRepository
	UserCourseRepo rp.UserCourseRepository
}

func NewCourseController(logger echo.Logger, courseRepo rp.CourseRepository, userCourse rp.UserCourseRepository) (ctr *CourseController) {
	ctr = &CourseController{cm.BaseController{}, courseRepo, userCourse}
	ctr.Init(logger)
	return
}

// GetCourseList : get list of courses(by courseName keyword)
// Params : echo.Context
// Returns : return error
func (ctr *CourseController) GetCourseList(c echo.Context) error {
	// userProfile := c.Get("user_profile").(m.User)
	courseListParams := new(param.CourseListParams)

	if err := c.Bind(courseListParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(courseListParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	courses, totalRow, err := ctr.CourseRepo.GetCourses(courseListParams)

	if err != nil {
		if err.Error() == pg.ErrNoRows.Error() {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Get course list failed",
			})
		}

		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System Error",
		})
	}

	if courseListParams.RowPerPage == 0 {
		courseListParams.CurrentPage = 1
		courseListParams.RowPerPage = totalRow
	}

	pagination := map[string]interface{}{
		"current_page": courseListParams.CurrentPage,
		"total_row":    totalRow,
		"row_per_page": courseListParams.RowPerPage,
	}

	listCourseResponse := []map[string]interface{}{}
	for i := 0; i < len(courses); i++ {
		itemDataResponse := map[string]interface{}{
			"course_id":          courses[i].ID,
			"course_title":       courses[i].Title,
			"course_description": courses[i].Description,
			"course_thumbnail":   courses[i].Thumbnail,
			"created_by":         courses[i].CreatedBy,
			"created_at":         courses[i].CreatedAt.Format(cf.FormatDateDisplay),
			"updated_at":         courses[i].UpdatedAt.Format(cf.FormatDateDisplay),
		}

		listCourseResponse = append(listCourseResponse, itemDataResponse)
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
		Data: map[string]interface{}{
			"pagination": pagination,
			"courses":    listCourseResponse,
		},
	})
}

// AddCourse : add new Course to database
// Params : echo.Context
// Returns : return error
func (ctr *CourseController) AddCourse(c echo.Context) error {
	createCourseParams := new(param.CreateCourseParams)
	if err := c.Bind(createCourseParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(createCourseParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	course, err := ctr.CourseRepo.SaveCourse(createCourseParams, ctr.UserCourseRepo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Create Course Failed",
		})
	}

	courseResponse := map[string]interface{}{
		"course_id":          course.ID,
		"course_title":       course.Title,
		"course_thumbnail":   course.Thumbnail,
		"course_description": course.Description,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Course Created Successful",
		Data:    courseResponse,
	})
}

// UpdateCourse : update project
// Params : echo.Context
// Returns : return error
func (ctr *CourseController) UpdateCourse(c echo.Context) error {
	updateCourseParams := new(param.UpdateCourseParams)
	if err := c.Bind(updateCourseParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(updateCourseParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	_, er := ctr.CourseRepo.GetCourseByID(updateCourseParams.ID)

	if er != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Course not found",
			Data:    er,
		})
	}

	if err := ctr.CourseRepo.UpdateCourse(updateCourseParams, ctr.UserCourseRepo); err != nil {
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System Error",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Update project success",
	})
}

// DeleteCourse : delete course by id
// Params : echo.Context
// Returns : object
func (ctr *CourseController) DeleteCourse(c echo.Context) error {
	courseIDParam := new(param.CourseIDParam)
	if err := c.Bind(courseIDParam); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(courseIDParam); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	_, er := ctr.CourseRepo.GetCourseByID(courseIDParam.CourseID)

	if er != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Course not found",
			Data:    er,
		})
	}

	err := ctr.UserCourseRepo.DeleteByCourseId(courseIDParam.CourseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System Error",
		})
	}

	err = ctr.CourseRepo.DeleteCourse(courseIDParam.CourseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System Error",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Deleted",
	})
}
