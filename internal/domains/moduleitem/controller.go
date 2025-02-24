package moduleitem

import (
	"net/http"
	"net/url"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	cld "orientation-training-api/internal/platform/cloud"
	"strconv"
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
	moduleId, err := strconv.Atoi(c.FormValue("module_id"))
	if err != nil || moduleId <= 0 {
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid module_id",
		})
	}

	itemType := c.FormValue("item_type")
	if itemType == "" || (itemType != "video" && itemType != "file") {
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid item_type. Allowed values: video, file",
		})
	}

	title := c.FormValue("title")
	if title == "" {
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Title is required",
		})
	}

	createModuleItemParams := &param.CreateModuleItemParams{
		Title:    title,
		ItemType: itemType,
		ModuleID: moduleId,
	}

	if itemType == "video" {
		videoURL := c.FormValue("url")
		if !valid.IsURL(videoURL) {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid URL format for video",
			})
		}

		parsedURL, err := url.Parse(videoURL)
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

		createModuleItemParams.Resource = videoId
	}

	if itemType == "file" {
		uploadedFile, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "File is required for item_type=file",
			})
		}

		if uploadedFile.Size > 10*1024*1024 {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "File size exceeds limit (10MB)",
			})
		}

		src, err := uploadedFile.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to open uploaded file",
			})
		}
		defer src.Close()

		millisecondTimeNow := int(time.Now().UnixNano() / int64(time.Millisecond))
		fileName := strconv.Itoa(moduleId) + "_" + strconv.Itoa(millisecondTimeNow)

		err = ctr.cloud.UploadFileToCloud(src, fileName, cf.FileFolderCLD)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to upload file to cloud",
			})
		}

		createModuleItemParams.Resource = fileName
	}

	savedItem, err := ctr.ModuleItemRepo.SaveModuleItem(createModuleItemParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to save Module Item to database",
		})
	}

	moduleItemResponse := map[string]interface{}{
		"id":        savedItem.ID,
		"module_id": savedItem.ModuleID,
		"type":      savedItem.ItemType,
		"title":     savedItem.Title,
		"resource":  savedItem.Resource,
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
		err := ctr.cloud.DeleteFileCloud(moduleItem.Resource, cf.FileFolderCLD)
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
