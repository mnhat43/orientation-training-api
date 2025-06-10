package courses

import (
	"fmt"
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	resp "orientation-training-api/internal/interfaces/response"
	m "orientation-training-api/internal/models"
	gc "orientation-training-api/internal/platform/cloud"
	"orientation-training-api/internal/platform/utils"
	"orientation-training-api/internal/platform/youtube"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-pg/pg"
	"github.com/labstack/echo/v4"
)

type CourseController struct {
	cm.BaseController

	CourseRepo             rp.CourseRepository
	UserCourseRepo         rp.UserCourseRepository
	UserProgressRepo       rp.UserProgressRepository
	ModuleRepo             rp.ModuleRepository
	ModuleItemRepo         rp.ModuleItemRepository
	UserRepo               rp.UserRepository
	CourseSkillKeywordRepo rp.CourseSkillKeywordRepository
	cloud                  gc.StorageUtility
}

func NewCourseController(
	logger echo.Logger,
	courseRepo rp.CourseRepository,
	ucRepo rp.UserCourseRepository,
	upRepo rp.UserProgressRepository,
	moduleRepo rp.ModuleRepository,
	moduleItemRepo rp.ModuleItemRepository,
	userRepo rp.UserRepository,
	courseSkillKeywordRepo rp.CourseSkillKeywordRepository,
	cloud gc.StorageUtility,
) (ctr *CourseController) {
	ctr = &CourseController{
		cm.BaseController{},
		courseRepo,
		ucRepo,
		upRepo,
		moduleRepo,
		moduleItemRepo,
		userRepo,
		courseSkillKeywordRepo,
		cloud,
	}
	ctr.Init(logger)
	return
}

// GetCourseList : get list of all courses without pagination or filtering
// Params : echo.Context
// Returns : return error
func (ctr *CourseController) GetCourseList(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	var courses []m.Course
	var err error

	if userProfile.RoleID == cf.ManagerRoleID {
		courses, err = ctr.CourseRepo.GetAllCourses()
	} else if userProfile.RoleID == cf.EmployeeRoleID {
		courses, err = ctr.CourseRepo.GetUserCourses(userProfile.ID)
	} else {
		courses, err = ctr.CourseRepo.GetAllCourses()
	}

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

	listCourseResponse := []map[string]interface{}{}
	for _, course := range courses {
		thumbnailURL := ""
		if course.Thumbnail != "" {
			thumbnailURL = ctr.cloud.GetURL(course.Thumbnail, cf.ThumbnailFolderGCS)
		}

		itemDataResponse := map[string]interface{}{
			"course_id":   course.ID,
			"title":       course.Title,
			"description": course.Description,
			"thumbnail":   thumbnailURL,
			"category":    course.Category,
			"duration":    course.Duration,
			"created_by":  course.CreatedBy,
			"created_at":  course.CreatedAt.Format(cf.FormatDateDisplay),
			"updated_at":  course.UpdatedAt.Format(cf.FormatDateDisplay),
		}
		// Get skill keywords for the course
		skillKeywords, err := ctr.CourseSkillKeywordRepo.GetSkillKeywordsByCourseID(course.ID)
		if err == nil {
			skillKeywordNames := []string{}
			for _, sk := range skillKeywords {
				skillKeywordNames = append(skillKeywordNames, sk.Name)
			}
			itemDataResponse["skill_keyword"] = skillKeywordNames
		} else {
			ctr.Logger.Errorf("Failed to fetch skill keywords for course %d: %v", course.ID, err)
			itemDataResponse["skill_keyword"] = []string{}
		}

		if userProfile.RoleID == cf.EmployeeRoleID {
			userProgress, err := ctr.UserProgressRepo.GetSingleUserProgress(userProfile.ID, course.ID)
			if err == nil && userProgress.ID > 0 {
				itemDataResponse["course_position"] = userProgress.CoursePosition
				itemDataResponse["module_position"] = userProgress.ModulePosition
				itemDataResponse["module_item_position"] = userProgress.ModuleItemPosition
				itemDataResponse["completed"] = userProgress.Completed

				// Add assessment information for completed courses
				if userProgress.Completed && userProgress.ReviewedBy > 0 {
					assessment := resp.Assessment{
						PerformanceRating:  userProgress.PerformanceRating,
						PerformanceComment: userProgress.PerformanceComment,
					}

					// Get reviewer information
					if userProgress.ReviewedBy > 0 {
						reviewer, err := ctr.UserRepo.GetUserProfile(userProgress.ReviewedBy)
						if err == nil {
							assessment.ReviewerName = reviewer.UserProfile.FirstName + " " + reviewer.UserProfile.LastName
						}
					}

					itemDataResponse["assessment"] = assessment
				}
			}
		}

		listCourseResponse = append(listCourseResponse, itemDataResponse)
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
		Data: map[string]interface{}{
			"courses": listCourseResponse,
			"total":   len(courses),
		},
	})
}

// AddCourse : add new Course to database
// Params : echo.Context
// Returns : return error
func (ctr *CourseController) AddCourse(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
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
	createCourseParams.CreatedBy = userProfile.ID
	course, err := ctr.CourseRepo.SaveCourse(createCourseParams, ctr.CourseSkillKeywordRepo)
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
	// Delete course skill keywords first
	err := ctr.CourseSkillKeywordRepo.DeleteByCourseID(courseIDParam.CourseID)
	if err != nil {
		ctr.Logger.Errorf("Failed to delete course skill keywords: %v", err)
		// Continue even if there's an error, as the ON DELETE CASCADE in the database should handle this
	}

	// Then delete user course relations
	err = ctr.UserCourseRepo.DeleteByCourseId(courseIDParam.CourseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System Error",
		})
	}

	// Finally delete the course
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

// GetCourseDetail retrieves detailed information about a specific course
// Params: echo.Context
// Returns: error
func (ctr *CourseController) GetCourseDetail(c echo.Context) error {
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

	course, err := ctr.CourseRepo.GetCourseByID(courseIDParam.CourseID)
	if err != nil {
		if err.Error() == pg.ErrNoRows.Error() {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Course not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System Error",
		})
	}

	modules, err := ctr.ModuleRepo.GetModulesByCourseID(course.ID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch modules: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch course modules",
		})
	}

	moduleList := []map[string]interface{}{}
	ytService := youtube.NewYouTubeService()

	for _, module := range modules {
		moduleItems, err := ctr.ModuleItemRepo.GetModuleItemsByModuleID(module.ID)
		if err != nil {
			ctr.Logger.Errorf("Failed to fetch module items for module %d: %v", module.ID, err)
			continue
		}
		itemsList := []map[string]interface{}{}
		for _, item := range moduleItems {
			itemData := map[string]interface{}{
				"id":        item.ID,
				"title":     item.Title,
				"item_type": item.ItemType,
				"position":  item.Position,
			}

			if item.ItemType != "quiz" {
				duration := 0

				if item.ItemType == "file" {
					duration = item.RequiredTime
				} else if item.ItemType == "video" {
					videoID := item.Resource

					videoInfo, err := ytService.GetVideoDetails(videoID)
					if err == nil && videoInfo != nil {
						duration = utils.ParseDurationToSeconds(videoInfo.Duration)
					} else {
						ctr.Logger.Errorf("Failed to fetch video details for %s: %v", videoID, err)
						duration = item.RequiredTime
					}
				}

				itemData["required_time"] = item.RequiredTime
				itemData["duration"] = duration
			}

			itemsList = append(itemsList, itemData)
		}

		moduleData := map[string]interface{}{
			"id":           module.ID,
			"title":        module.Title,
			"position":     module.Position,
			"duration":     module.Duration,
			"module_items": itemsList,
		}

		moduleList = append(moduleList, moduleData)
	}

	courseDetail := map[string]interface{}{
		"course_id":   course.ID,
		"title":       course.Title,
		"description": course.Description,
		"category":    course.Category,
		"duration":    course.Duration,
		"modules":     moduleList,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Course details retrieved successfully",
		Data:    courseDetail,
	})
}
