package moduleitem

import (
	"net/http"
	"net/url"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	cld "orientation-training-api/internal/platform/cloud"
	"orientation-training-api/internal/platform/utils"
	"orientation-training-api/internal/platform/youtube"
	"strconv"
	"strings"
	"time"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
)

type ModuleItemController struct {
	cm.BaseController

	ModuleItemRepo rp.ModuleItemRepository
	cloud          cld.StorageUtility
}

func NewModuleItemController(logger echo.Logger, moduleItemRepo rp.ModuleItemRepository, cloud cld.StorageUtility) (ctr *ModuleItemController) {
	ctr = &ModuleItemController{cm.BaseController{}, moduleItemRepo, cloud}
	ctr.Init(logger)
	return
}

// GetModuleItemList : get list of moduleItems(by moduleName keyword)
// Params : echo.Context
// Returns : return error
func (ctr *ModuleItemController) GetModuleItemList(c echo.Context) error {
	// userProfile := c.Get("user_profile").(m.User)
	moduleItemListParams := new(param.ModuleItemListParams)

	if err := c.Bind(moduleItemListParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(moduleItemListParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	moduleItems, totalRow, err := ctr.ModuleItemRepo.GetModuleItems(moduleItemListParams)

	if err != nil {
		if err.Error() == pg.ErrNoRows.Error() {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Get module item list failed",
			})
		}

		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System Error",
		})
	}

	if moduleItemListParams.RowPerPage == 0 {
		moduleItemListParams.CurrentPage = 1
		moduleItemListParams.RowPerPage = totalRow
	}

	pagination := map[string]interface{}{
		"current_page": moduleItemListParams.CurrentPage,
		"total_row":    totalRow,
		"row_per_page": moduleItemListParams.RowPerPage,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
		Data: map[string]interface{}{
			"pagination":  pagination,
			"moduleItems": moduleItems,
		},
	})
}

// AddModuleItem : add new ModuleItem to database
// Params : echo.Context
// Returns : return error
func (ctr *ModuleItemController) AddModuleItem(c echo.Context) error {
	createModuleItemParams := new(param.CreateModuleItemParams)

	if err := c.Bind(createModuleItemParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(createModuleItemParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	if createModuleItemParams.ItemType == "" || (createModuleItemParams.ItemType != "video" && createModuleItemParams.ItemType != "file") {
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid item_type. Allowed values: video, file",
		})
	}

	if createModuleItemParams.ItemType == "video" {
		if !valid.IsURL(createModuleItemParams.Resource) {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid URL format for video",
			})
		}

		parsedURL, err := url.Parse(createModuleItemParams.Resource)
		if err != nil {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to parse video URL",
			})
		}

		queryParams := parsedURL.Query()
		videoId := queryParams.Get("v")
		if videoId == "" {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid YouTube video URL: missing 'v' parameter",
			})
		}

		ytService := youtube.NewYouTubeService()
		videoInfo, err := ytService.GetVideoDetails(videoId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to fetch video details",
			})
		}

		requiredTimeInSeconds := utils.CalculateRequiredTime(videoInfo.Duration)

		createModuleItemParams.RequiredTime = requiredTimeInSeconds
		createModuleItemParams.Resource = videoId
	} else if createModuleItemParams.ItemType == "file" {
		parts := strings.SplitN(createModuleItemParams.Resource, ",", 2)
		if len(parts) != 2 {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid File Format",
			})
		}

		mimeType := parts[0]
		base64Data := parts[1]

		formatFile := ""
		if strings.HasPrefix(mimeType, "data:application/") {
			formatFile = strings.TrimPrefix(mimeType, "data:application/")
			formatFile = strings.Split(formatFile, ";")[0]
		}

		if formatFile == "" {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid File Format",
			})
		}

		if _, check := utils.FindStringInArray(cf.AllowFormatFileList, formatFile); !check {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "File not allowed",
			})
		}

		millisecondTimeNow := int(time.Now().UnixNano() / int64(time.Millisecond))
		fileName := strconv.Itoa(createModuleItemParams.ModuleID) + "_" + strconv.Itoa(millisecondTimeNow)

		err := ctr.cloud.UploadFileToCloud(
			base64Data,
			fileName,
			cf.FileFolderGCS,
		)
		if err != nil {
			ctr.Logger.Error(err)
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to upload file to cloud",
			})
		}

		createModuleItemParams.Resource = fileName
	}

	savedItem, err := ctr.ModuleItemRepo.SaveModuleItem(createModuleItemParams)
	if err != nil {
		ctr.Logger.Error(err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to save Module Item to database",
		})
	}

	moduleItemResponse := map[string]interface{}{
		"id":        savedItem.ID,
		"type":      savedItem.ItemType,
		"title":     savedItem.Title,
		"resource":  savedItem.Resource,
		"position":  savedItem.Position,
		"required_time":  savedItem.RequiredTime,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Module Item Created Successfully",
		Data:    moduleItemResponse,
	})
}

// DeleteModuleItem : delete module item by id
// Params : echo.Context
// Returns : object
func (ctr *ModuleItemController) DeleteModuleItem(c echo.Context) error {
	moduleItemIDParam := new(param.ModuleItemIDParam)
	if err := c.Bind(moduleItemIDParam); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(moduleItemIDParam); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	moduleItem, er := ctr.ModuleItemRepo.GetModuleItemByID(moduleItemIDParam.ModuleItemID)

	if er != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Module item not found",
			Data:    er,
		})
	}

	if moduleItem.ItemType == "file" {
		err := ctr.cloud.DeleteFileCloud(moduleItem.Resource, cf.FileFolderGCS)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to delete file from cloud",
				Data:    err,
			})
		}
	}
	err := ctr.ModuleItemRepo.DeleteModuleItem(moduleItemIDParam.ModuleItemID)

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
