package modules

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"

	valid "github.com/asaskevich/govalidator"

	"github.com/go-pg/pg/v9"
	"github.com/labstack/echo/v4"
)

type ModuleController struct {
	cm.BaseController

	ModuleRepo     rp.ModuleRepository
	ModuleItemRepo rp.ModuleItemRepository
	CourseRepo     rp.CourseRepository
	// cloud            cld.StorageUtility
}

func NewModuleController(logger echo.Logger, moduleRepo rp.ModuleRepository, moduleItemRepo rp.ModuleItemRepository, courseRepo rp.CourseRepository) (ctr *ModuleController) {
	ctr = &ModuleController{cm.BaseController{}, moduleRepo, moduleItemRepo, courseRepo}
	ctr.Init(logger)
	return
}

// GetModuleList : get list of modules(by moduleName keyword)
// Params : echo.Context
// Returns : return error
func (ctr *ModuleController) GetModuleList(c echo.Context) error {
	// userProfile := c.Get("user_profile").(m.User)
	moduleListParams := new(param.ModuleListParams)

	if err := c.Bind(moduleListParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(moduleListParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	modules, totalRow, err := ctr.ModuleRepo.GetModules(moduleListParams)

	if err != nil {
		if err.Error() == pg.ErrNoRows.Error() {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Get module list failed",
			})
		}

		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System Error",
		})
	}

	if moduleListParams.RowPerPage == 0 {
		moduleListParams.CurrentPage = 1
		moduleListParams.RowPerPage = totalRow
	}

	pagination := map[string]interface{}{
		"current_page": moduleListParams.CurrentPage,
		"total_row":    totalRow,
		"row_per_page": moduleListParams.RowPerPage,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
		Data: map[string]interface{}{
			"pagination": pagination,
			"modules":    modules,
		},
	})
}

func (ctr *ModuleController) GetModuleDetails(c echo.Context) error {
	// userProfile := c.Get("user_profile").(m.User)
	moduleListParams := new(param.ModuleListParams)

	if err := c.Bind(moduleListParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(moduleListParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	modules, totalRow, err := ctr.ModuleRepo.GetModules(moduleListParams)

	if err != nil {
		if err.Error() == pg.ErrNoRows.Error() {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Get module list failed",
			})
		}

		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System Error",
		})
	}

	if moduleListParams.RowPerPage == 0 {
		moduleListParams.CurrentPage = 1
		moduleListParams.RowPerPage = totalRow
	}

	pagination := map[string]interface{}{
		"current_page": moduleListParams.CurrentPage,
		"total_row":    totalRow,
		"row_per_page": moduleListParams.RowPerPage,
	}

	listModuleResponse := []map[string]interface{}{}
	for _, module := range modules {
		moduleItemParams := &param.ModuleItemListParams{ModuleID: module.ID}
		moduleListItem, _, err := ctr.ModuleItemRepo.GetModuleItems(moduleItemParams)

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

		var moduleItemResponse []map[string]interface{}

		for _, moduleItem := range moduleListItem {
			moduleItemResponse = append(moduleItemResponse, map[string]interface{}{
				"id":        moduleItem.ID,
				"title":     moduleItem.Title,
				"item_type": moduleItem.ItemType,
				"resource":  moduleItem.Resource,
				"position":  moduleItem.Position,
			})
		}

		itemDataResponse := map[string]interface{}{
			"id":    module.ID,
			"title": module.Title,
			"position": module.Position,
			"module_items": moduleItemResponse,
		}

		listModuleResponse = append(listModuleResponse, itemDataResponse)
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
		Data: map[string]interface{}{
			"pagination": pagination,
			"modules":    listModuleResponse,
		},
	})
}

// AddModule : add new Module to database
// Params : echo.Context
// Returns : return error
func (ctr *ModuleController) AddModule(c echo.Context) error {
	createModuleParams := new(param.CreateModuleParams)

	if err := c.Bind(createModuleParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(createModuleParams); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	module, err := ctr.ModuleRepo.SaveModule(createModuleParams)
	if err != nil {
		ctr.Logger.Errorf("Error creating module: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Create Module Failed",
			Data:    err.Error(),
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Module Created Successfully",
		Data:    module,
	})
}

// DeleteModule : delete module by id
// Params : echo.Context
// Returns : object
func (ctr *ModuleController) DeleteModule(c echo.Context) error {
	moduleIDParam := new(param.ModuleIDParam)
	if err := c.Bind(moduleIDParam); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(moduleIDParam); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	_, er := ctr.ModuleRepo.GetModuleByID(moduleIDParam.ModuleID)

	if er != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Module not found",
			Data:    er,
		})
	}
	err := ctr.ModuleRepo.DeleteModule(moduleIDParam.ModuleID)

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
