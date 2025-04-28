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
// Can both track progress and mark a course as completed
// Params: echo.Context
// Returns: error
func (ctr *UserProgressController) UpdateUserProgress(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	updateUserProgressParams := new(param.UpdateUserProgressParams)

	if err := c.Bind(updateUserProgressParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(updateUserProgressParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	// Create user progress entry
	userProgress := &m.UserProgress{
		UserID:   userProfile.ID,
		CourseID: updateUserProgressParams.CourseID,
	}

	// If the completed flag is set, we need to get the current progress first
	if updateUserProgressParams.Completed {
		currentProgress, err := ctr.UserProgressRepo.GetUserProgress(userProfile.ID, updateUserProgressParams.CourseID)
		if err != nil {
			ctr.Logger.Errorf("Failed to fetch existing progress: %v", err)
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to fetch user progress",
			})
		}

		// For marking complete, keep existing positions but set completed=true
		userProgress.ModulePosition = currentProgress.ModulePosition
		userProgress.ModuleItemPosition = currentProgress.ModuleItemPosition
		userProgress.Completed = true
	} else {
		// For regular progress updates, use the provided positions and set completed=false
		userProgress.ModulePosition = updateUserProgressParams.ModulePosition
		userProgress.ModuleItemPosition = updateUserProgressParams.ModuleItemPosition
		userProgress.Completed = false
	}

	err := ctr.UserProgressRepo.SaveUserProgress(userProgress)
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

// GetUserProgress retrieves a user's progress for a specific course
// Params: echo.Context
// Returns: error
func (ctr *UserProgressController) GetUserProgress(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	getUserProgressParams := new(param.GetUserProgressParams)

	if err := c.Bind(getUserProgressParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(getUserProgressParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	userProgress, err := ctr.UserProgressRepo.GetUserProgress(userProfile.ID, getUserProgressParams.CourseID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch user progress: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch user progress",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "User progress retrieved successfully",
		Data:    userProgress,
	})
}

// AddUserProgress creates a new user progress record
// Params: echo.Context
// Returns: error
func (ctr *UserProgressController) AddUserProgress(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	createUserProgressParams := new(param.CreateUserProgressParams)

	if err := c.Bind(createUserProgressParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(createUserProgressParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	// Check if user progress already exists
	existingProgress, err := ctr.UserProgressRepo.GetUserProgress(userProfile.ID, createUserProgressParams.CourseID)
	if err == nil && existingProgress.ID != 0 {
		// Progress already exists, return appropriate response
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "User progress already exists for this course",
			Data:    existingProgress,
		})
	}

	// Create new user progress entry
	userProgress := &m.UserProgress{
		UserID:             userProfile.ID,
		CourseID:           createUserProgressParams.CourseID,
		ModulePosition:     createUserProgressParams.ModulePosition,
		ModuleItemPosition: createUserProgressParams.ModuleItemPosition,
		Completed:          false,
	}

	err = ctr.UserProgressRepo.SaveUserProgress(userProgress)
	if err != nil {
		ctr.Logger.Errorf("Failed to create user progress: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to create progress",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Progress created successfully",
		Data:    userProgress,
	})
}
