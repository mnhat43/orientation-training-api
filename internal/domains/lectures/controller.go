package lectures

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
	cld "orientation-training-api/internal/platform/cloud"
	"orientation-training-api/internal/platform/youtube"
	"os"

	valid "github.com/asaskevich/govalidator"

	"github.com/labstack/echo/v4"
)

type LectureController struct {
	cm.BaseController

	ModuleRepo       rp.ModuleRepository
	ModuleItemRepo   rp.ModuleItemRepository
	CourseRepo       rp.CourseRepository
	UserProgressRepo rp.UserProgressRepository
	Cloud            cld.StorageUtility
}

func NewLectureController(logger echo.Logger, moduleRepo rp.ModuleRepository, moduleItemRepo rp.ModuleItemRepository, courseRepo rp.CourseRepository, upRepo rp.UserProgressRepository, cloud cld.StorageUtility) (ctr *LectureController) {
	ctr = &LectureController{cm.BaseController{}, moduleRepo, moduleItemRepo, courseRepo, upRepo, cloud}
	ctr.Init(logger)
	return
}

func (ctr *LectureController) GetLectureList(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	lectureListParams := new(param.LectureListParams)

	if err := c.Bind(lectureListParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(lectureListParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	// Get user progress with explicit userID from authenticated user
	userProgress, err := ctr.UserProgressRepo.GetSingleUserProgress(userProfile.ID, lectureListParams.CourseID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch user progress: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch user progress",
		})
	}

	currentModulePosition, currentModuleItemPosition := userProgress.ModulePosition, userProgress.ModuleItemPosition

	moduleListParams := &param.ModuleListParams{
		CourseID: lectureListParams.CourseID,
	}
	modules, _, err := ctr.ModuleRepo.GetModules(moduleListParams)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch modules: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch modules for the course",
		})
	}

	moduleIDs, err := ctr.ModuleRepo.GetModuleIDsByCourseID(lectureListParams.CourseID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch module IDs: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch modules for the course",
		})
	}

	moduleItems, err := ctr.ModuleItemRepo.GetModuleItemsByModuleIDs(moduleIDs)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch module items: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch module items",
		})
	}

	modulePositions := make(map[int]int)
	for _, module := range modules {
		modulePositions[module.ID] = module.Position
	}

	lectureList := map[string][]map[string]interface{}{}
	for _, module := range modules {
		lectureList[module.Title] = []map[string]interface{}{}
	}

	for _, item := range moduleItems {
		moduleTitle := ""
		modulePosition := 0

		for _, module := range modules {
			if module.ID == item.ModuleID {
				moduleTitle = module.Title
				modulePosition = module.Position
				break
			}
		}

		isUnlocked := false
		if modulePosition < currentModulePosition ||
			(modulePosition == currentModulePosition && item.Position <= currentModuleItemPosition) {
			isUnlocked = true
		}

		if item.ItemType == "video" {
			videoID := item.Resource

			ytService := youtube.NewYouTubeService()

			videoInfo, err := ytService.GetVideoDetails(videoID)
			if err != nil {
				ctr.Logger.Errorf("Failed to fetch video details for video ID %s: %v", videoID, err)
				continue
			}

			lectureData := map[string]interface{}{
				"module_item_id":       item.ID,
				"title":                item.Title,
				"item_type":            item.ItemType,
				"module_item_position": item.Position,
				"module_id":            item.ModuleID,
				"module_position":      modulePosition,
				"videoId":              videoID,
				"thumbnail":            videoInfo.ThumbnailURL,
				"duration":             videoInfo.Duration,
				"publishedAt":          videoInfo.PublishedAt,
				"required_time":        item.RequiredTime,
				"unlocked":             isUnlocked,
			}
			lectureList[moduleTitle] = append(lectureList[moduleTitle], lectureData)
		} else if item.ItemType == "file" {
			var filePath string
			if item.Resource != "" {
				filePath = "https://storage.cloud.google.com/" + os.Getenv("GOOGLE_STORAGE_BUCKET") + "/" + cf.FileFolderGCS + item.Resource
			}

			lectureData := map[string]interface{}{
				"module_item_id":       item.ID,
				"title":                item.Title,
				"item_type":            item.ItemType,
				"module_item_position": item.Position,
				"module_id":            item.ModuleID,
				"module_position":      modulePosition,
				"file_path":            filePath,
				"required_time":        item.RequiredTime,
				"duration":             item.RequiredTime,
				"unlocked":             isUnlocked,
			}
			lectureList[moduleTitle] = append(lectureList[moduleTitle], lectureData)
		}
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
		Data:    lectureList,
	})
}
