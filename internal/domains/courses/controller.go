package courses

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	"strconv"

	"github.com/go-pg/pg"
	"github.com/labstack/echo/v4"
)

type CourseController struct {
	cm.BaseController

	CourseRepo rp.CourseRepository
}

func NewCourseController(
	logger echo.Logger,
	courseRepo rp.CourseRepository,
) (ctr *CourseController) {
	ctr = &CourseController{
		cm.BaseController{},
		courseRepo,
	}
	ctr.Init(logger)
	return
}

// GetCourse retrieves a course by its ID
func (ctr *CourseController) GetCourse(c echo.Context) error {
	// Get the course ID from the URL parameter
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid course ID",
		})
	}

	// Retrieve the course from the repository
	course, err := ctr.CourseRepo.GetCourseByID(int(id))
	if err != nil {
		if err.Error() == pg.ErrNoRows.Error() {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "User is not exists",
			})
		}

		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System error",
		})
	}

	// Return the course details
	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
		Data:    course,
	})
}
