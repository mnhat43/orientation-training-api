package userprogress

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

type UserProgressController struct {
	cm.BaseController

	UserProgressRepo rp.UserProgressRepository
	ModuleRepo       rp.ModuleRepository
	ModuleItemRepo   rp.ModuleItemRepository
}

func NewUserProgressController(logger echo.Logger, userProgressRepo rp.UserProgressRepository, moduleRepo rp.ModuleRepository, moduleItemRepo rp.ModuleItemRepository) (ctr *UserProgressController) {
	ctr = &UserProgressController{cm.BaseController{}, userProgressRepo, moduleRepo, moduleItemRepo}
	ctr.Init(logger)
	return
}

// UpdateUserProgress updates a user's progress in a course
// Params: echo.Context
// Returns: error
func (ctr *UserProgressController) UpdateUserProgress(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	updateProgressParams := new(param.UpdateUserProgressParams)

	if err := c.Bind(updateProgressParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(updateProgressParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	// Override userID with the authenticated user's ID for security
	updateProgressParams.UserID = userProfile.ID

	// Get modules for this course
	moduleListParams := &param.ModuleListParams{
		CourseID: updateProgressParams.CourseID,
	}
	modules, _, err := ctr.ModuleRepo.GetModules(moduleListParams)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch modules: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch modules for the course",
		})
	}

	// Find the max module position
	maxModulePosition := 0
	var lastModuleID int
	for _, module := range modules {
		if module.Position > maxModulePosition {
			maxModulePosition = module.Position
			lastModuleID = module.ID
		}
	}

	// Get module items for the last module
	moduleItemListParams := &param.ModuleItemListParams{
		ModuleID: lastModuleID,
	}
	lastModuleItems, _, err := ctr.ModuleItemRepo.GetModuleItems(moduleItemListParams)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch module items: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch module items",
		})
	}

	// Find the max module item position in the last module
	maxModuleItemPosition := 0
	for _, item := range lastModuleItems {
		if item.Position > maxModuleItemPosition {
			maxModuleItemPosition = item.Position
		}
	}

	// Determine if the course is completed
	isCompleted := updateProgressParams.ModulePosition >= maxModulePosition && 
	               updateProgressParams.ModuleItemPosition >= maxModuleItemPosition

	// Create user progress entry with calculated completion status
	userProgress := &m.UserProgress{
		UserID:             updateProgressParams.UserID,
		CourseID:           updateProgressParams.CourseID,
		ModulePosition:     updateProgressParams.ModulePosition,
		ModuleItemPosition: updateProgressParams.ModuleItemPosition,
		Completed:          isCompleted,
	}

	err = ctr.UserProgressRepo.SaveUserProgress(userProgress)
	if err != nil {
		ctr.Logger.Errorf("Failed to save user progress: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to update progress",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Progress updated successfully",
		Data: map[string]interface{}{
			"module_position":      userProgress.ModulePosition,
			"module_item_position": userProgress.ModuleItemPosition,
			"completed":            userProgress.Completed,
		},
	})
}
