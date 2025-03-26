package lectures

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	cld "orientation-training-api/internal/platform/cloud"
	"orientation-training-api/internal/platform/youtube"
	"os"

	valid "github.com/asaskevich/govalidator"

	"github.com/labstack/echo/v4"
)

type LectureController struct {
	cm.BaseController

	ModuleRepo     rp.ModuleRepository
	ModuleItemRepo rp.ModuleItemRepository
	CourseRepo     rp.CourseRepository
	Cloud          cld.StorageUtility
}

func NewLectureController(logger echo.Logger, moduleRepo rp.ModuleRepository, moduleItemRepo rp.ModuleItemRepository, courseRepo rp.CourseRepository, cloud cld.StorageUtility) (ctr *LectureController) {
	ctr = &LectureController{cm.BaseController{}, moduleRepo, moduleItemRepo, courseRepo, cloud}
	ctr.Init(logger)
	return
}

func (ctr *LectureController) GetLectureList(c echo.Context) error {
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

	lectureList := map[string][]map[string]interface{}{}
	for _, module := range modules {
		lectureList[module.Title] = []map[string]interface{}{}
	}

	for _, item := range moduleItems {
		moduleTitle := ""
		for _, module := range modules {
			if module.ID == item.ModuleID {
				moduleTitle = module.Title
				break
			}
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
				"id":          item.ID,
				"itemType":    item.ItemType,
				"title":       item.Title,
				"videoId":     videoID,
				"thumbnail":   videoInfo.ThumbnailURL,
				"duration":    videoInfo.Duration,
				"publishedAt": videoInfo.PublishedAt,
			}
			lectureList[moduleTitle] = append(lectureList[moduleTitle], lectureData)
		} else if item.ItemType == "file" {
			var filePath string
			if item.Resource != "" {
				filePath = "https://storage.cloud.google.com/" + os.Getenv("GOOGLE_STORAGE_BUCKET") + "/" + cf.FileFolderGCS + item.Resource
			}

			lectureData := map[string]interface{}{
				"id":        item.ID,
				"itemType":  item.ItemType,
				"title":     item.Title,
				"file_path": filePath,
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
