package moduleitem

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

type ModuleItemController struct {
	cm.BaseController

	ModuleItemRepo rp.ModuleItemRepository
	ModuleRepo     rp.ModuleRepository
}

func NewModuleItemController(logger echo.Logger, moduleItemRepo rp.ModuleItemRepository, moduleRepo rp.ModuleRepository) (ctr *ModuleItemController) {
	ctr = &ModuleItemController{cm.BaseController{}, moduleItemRepo, moduleRepo}
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
