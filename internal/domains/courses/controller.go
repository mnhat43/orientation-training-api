package courses

import (
	"fmt"
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	gc "orientation-training-api/internal/platform/cloud"
	"orientation-training-api/internal/platform/utils"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-pg/pg"
	"github.com/labstack/echo/v4"
)

type CourseController struct {
	cm.BaseController

	CourseRepo     rp.CourseRepository
	UserCourseRepo rp.UserCourseRepository
	cloud          gc.StorageUtility
}

func NewCourseController(logger echo.Logger, courseRepo rp.CourseRepository, userCourse rp.UserCourseRepository, cloud gc.StorageUtility) (ctr *CourseController) {
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
		var base64Img []byte = nil
		if course.Thumbnail != "" {
			base64Img, err = ctr.cloud.GetFileByFileName(course.Thumbnail, cf.ThumbnailFolderGCS)
			if err != nil {
				ctr.Logger.Error(err)
				base64Img = nil
			}
		}

		itemDataResponse := map[string]interface{}{
			"course_id":          course.ID,
			"course_title":       course.Title,
			"course_description": course.Description,
			"course_thumbnail":   base64Img,
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

	if createCourseParams.Thumbnail != "" {
		parts := strings.SplitN(createCourseParams.Thumbnail, ",", 2)
		if len(parts) != 2 {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid Thumbnail Format",
			})
		}

		mimeType := parts[0]
		base64Data := parts[1]

		formatImageThumbnail := ""
		if strings.HasPrefix(mimeType, "data:image/") {
			formatImageThumbnail = strings.TrimPrefix(mimeType, "data:image/")
			formatImageThumbnail = strings.Split(formatImageThumbnail, ";")[0]
		}

		if formatImageThumbnail == "" {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid Image Format",
			})
		}

		if _, check := utils.FindStringInArray(cf.AllowFormatImageList, formatImageThumbnail); !check {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "The Thumbnail field must be an image",
			})
		}

		millisecondTimeNow := int(time.Now().UnixNano() / int64(time.Millisecond))
		nameThumbnail := fmt.Sprintf("%d_%d.%s", createCourseParams.CreatedBy, millisecondTimeNow, formatImageThumbnail)

		err := ctr.cloud.UploadFileToCloud(base64Data, nameThumbnail, cf.ThumbnailFolderGCS)
		if err != nil {
			ctr.Logger.Error(err)
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Upload thumbnail error",
			})
		}

		createCourseParams.Thumbnail = nameThumbnail
	}

	course, err := ctr.CourseRepo.SaveCourse(createCourseParams, ctr.UserCourseRepo)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Create Course Failed",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Course Created Successfully",
		Data:    course,
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
		err := ctr.cloud.DeleteFileCloud(course.Thumbnail, cf.ThumbnailFolderGCS)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "System Error: Failed to delete thumbnail from cloud",
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
