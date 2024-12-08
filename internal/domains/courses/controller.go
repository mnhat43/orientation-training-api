package courses

import (
	"io"
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	cld "orientation-training-api/internal/platform/cloud"
	"strconv"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-pg/pg"
	"github.com/labstack/echo/v4"
)

type CourseController struct {
	cm.BaseController

	CourseRepo     rp.CourseRepository
	UserCourseRepo rp.UserCourseRepository
	cloud          cld.StorageUtility
}

func NewCourseController(logger echo.Logger, courseRepo rp.CourseRepository, userCourse rp.UserCourseRepository, cloud cld.StorageUtility) (ctr *CourseController) {
	ctr = &CourseController{cm.BaseController{}, courseRepo, userCourse, cloud}
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
	for _, course := range courses {
		var secureURL interface{}
		if course.Thumbnail != "" {
			secureURL = ctr.cloud.GetFileByFileName(course.Thumbnail, cf.ThumbnailFolderCLD)
		}

		itemDataResponse := map[string]interface{}{
			"course_id":          course.ID,
			"course_title":       course.Title,
			"course_description": course.Description,
			"course_thumbnail":   secureURL,
			"created_by":         course.CreatedBy,
			"created_at":         course.CreatedAt.Format(cf.FormatDateDisplay),
			"updated_at":         course.UpdatedAt.Format(cf.FormatDateDisplay),
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
	title := c.FormValue("course_title")
	description := c.FormValue("course_description")
	createdBy, err := strconv.Atoi(c.FormValue("created_by"))

	if err != nil {
		return c.JSON(http.StatusUnauthorized, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid created_by",
		})
	}

	createCourseParams := &param.CreateCourseParams{
		Title:       title,
		Description: description,
		CreatedBy:   createdBy,
	}
	if _, err := valid.ValidateStruct(createCourseParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	createCourseDBParams := &param.CreateCourseDBParams{
		Title:       title,
		Description: description,
		CreatedBy:   createdBy,
	}

	thumbnailFile, err := c.FormFile("course_thumbnail")
	if err == nil {
		src, err := thumbnailFile.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Error opening thumbnail file",
			})
		}
		defer src.Close()

		src.Seek(0, io.SeekStart)

		millisecondTimeNow := int(time.Now().UnixNano() / int64(time.Millisecond))
		nameThumbnail := strconv.Itoa(createCourseParams.CreatedBy) + "_" + strconv.Itoa(millisecondTimeNow)

		err = ctr.cloud.UploadFileToCloud(src, nameThumbnail, cf.ThumbnailFolderCLD)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Upload thumbnail error",
			})
		}

		createCourseDBParams.Thumbnail = nameThumbnail
	}

	course, err := ctr.CourseRepo.SaveCourse(createCourseDBParams, ctr.UserCourseRepo)
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
		Message: "Course Created Successfully",
		Data:    courseResponse,
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

	course, er := ctr.CourseRepo.GetCourseByID(courseIDParam.CourseID)

	if er != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Course not found",
			Data:    er,
		})
	}

	if course.Thumbnail != "" {
		err := ctr.cloud.DeleteFileCloud(course.Thumbnail, cf.ThumbnailFolderCLD)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "System Error: Failed to delete thumbnail from Cloudinary",
			})
		}
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

// func (ctr *CourseController) GetCourseDetail(c echo.Context) error {
// 	courseIDParam := new(param.CourseIDParam)
// 	if err := c.Bind(courseIDParam); err != nil {
// 		return c.JSON(http.StatusOK, cf.JsonResponse{
// 			Status:  cf.FailResponseCode,
// 			Message: "Invalid Params",
// 			Data:    err,
// 		})
// 	}

// 	if _, err := valid.ValidateStruct(courseIDParam); err != nil {
// 		return c.JSON(http.StatusOK, cf.JsonResponse{
// 			Status:  cf.FailResponseCode,
// 			Message: err.Error(),
// 		})
// 	}

// 	course, er := ctr.CourseRepo.GetCourseByID(courseIDParam.CourseID)

// 	if er != nil {
// 		return c.JSON(http.StatusOK, cf.JsonResponse{
// 			Status:  cf.FailResponseCode,
// 			Message: "Course not found",
// 			Data:    er,
// 		})
// 	}

// 	if course.Thumbnail != "" {
// 		err := ctr.cloud.DeleteFileCloud(course.Thumbnail, cf.ThumbnailFolderCLD)
// 		if err != nil {
// 			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
// 				Status:  cf.FailResponseCode,
// 				Message: "System Error: Failed to delete thumbnail from Cloudinary",
// 			})
// 		}
// 	}

// 	err := ctr.UserCourseRepo.DeleteByCourseId(courseIDParam.CourseID)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
// 			Status:  cf.FailResponseCode,
// 			Message: "System Error",
// 		})
// 	}

// 	err = ctr.CourseRepo.DeleteCourse(courseIDParam.CourseID)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
// 			Status:  cf.FailResponseCode,
// 			Message: "System Error",
// 		})
// 	}

// 	return c.JSON(http.StatusOK, cf.JsonResponse{
// 		Status:  cf.SuccessResponseCode,
// 		Message: "Deleted",
// 	})
// }
