package lectures

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	"orientation-training-api/internal/platform/youtube"

	valid "github.com/asaskevich/govalidator"

	"github.com/labstack/echo/v4"
)

type LectureController struct {
	cm.BaseController

	ModuleRepo     rp.ModuleRepository
	ModuleItemRepo rp.ModuleItemRepository
	CourseRepo     rp.CourseRepository
}

func NewLectureController(logger echo.Logger, moduleRepo rp.ModuleRepository, moduleItemRepo rp.ModuleItemRepository, courseRepo rp.CourseRepository) (ctr *LectureController) {
	ctr = &LectureController{cm.BaseController{}, moduleRepo, moduleItemRepo, courseRepo}
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

	moduleIDs, err := ctr.ModuleRepo.GetModuleIDsByCourseID(lectureListParams.CourseID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch module IDs: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch modules for the course",
		})
	}

	moduleItems, err := ctr.ModuleItemRepo.GetModuleItemsByModuleIDs(moduleIDs, "video")
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch module items: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch module items",
		})
	}

	lectureList := []map[string]interface{}{}
	for _, item := range moduleItems {
		videoID := item.Resource

		ytService := youtube.NewYouTubeService()

		videoInfo, err := ytService.GetVideoDetails(videoID)
		if err != nil {
			ctr.Logger.Errorf("Failed to fetch video details for video ID %s: %v", videoID, err)
			continue
		}

		lectureData := map[string]interface{}{
			"id":          item.ID,
			"title":       item.Title,
			"videoId":     videoID,
			"thumbnail":   videoInfo.ThumbnailURL,
			"duration":    videoInfo.Duration,
			"publishedAt": videoInfo.PublishedAt,
		}
		lectureList = append(lectureList, lectureData)
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
		Data:    lectureList,
	})
}
