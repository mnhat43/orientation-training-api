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
	UserRepo         rp.UserRepository
}

func NewUserProgressController(logger echo.Logger, userProgressRepo rp.UserProgressRepository, moduleRepo rp.ModuleRepository, moduleItemRepo rp.ModuleItemRepository, userRepo rp.UserRepository) (ctr *UserProgressController) {
	ctr = &UserProgressController{cm.BaseController{}, userProgressRepo, moduleRepo, moduleItemRepo, userRepo}
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

	userProgress := &m.UserProgress{
		UserID:   userProfile.ID,
		CourseID: updateUserProgressParams.CourseID,
	}

	if updateUserProgressParams.Completed {
		currentProgress, err := ctr.UserProgressRepo.GetUserProgress(userProfile.ID, updateUserProgressParams.CourseID)
		if err != nil {
			ctr.Logger.Errorf("Failed to fetch existing progress: %v", err)
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to fetch user progress",
			})
		}

		userProgress.ModulePosition = currentProgress.ModulePosition
		userProgress.ModuleItemPosition = currentProgress.ModuleItemPosition
		userProgress.Completed = true
	} else {
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

	existingProgress, err := ctr.UserProgressRepo.GetUserProgress(userProfile.ID, createUserProgressParams.CourseID)
	if err == nil && existingProgress.ID != 0 {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "User progress already exists for this course",
			Data:    existingProgress,
		})
	}

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

func (ctr *UserProgressController) GetListTraineeByCourseID(c echo.Context) error {
	getListProgressParams := new(param.GetListTraineeByCourseIDParams)

	if err := c.Bind(getListProgressParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(getListProgressParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	trainees, err := ctr.UserRepo.GetUsersByRoleID(cf.TraineeRoleID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch trainees: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch trainees",
		})
	}

	userProgressList, err := ctr.UserProgressRepo.GetUserProgressByCourseID(getListProgressParams.CourseID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch user progress list: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch user progress list",
		})
	}

	userProgressMap := make(map[int]m.UserProgress)
	for _, progress := range userProgressList {
		userProgressMap[progress.UserID] = progress
	}

	traineeInfoList := []map[string]interface{}{}

	for _, trainee := range trainees {
		status := cf.NotAssigned
		if progress, exists := userProgressMap[trainee.ID]; exists {
			if progress.Completed {
				status = cf.Completed
			} else {
				status = cf.InProgress
			}
		}

		traineeInfo := map[string]interface{}{
			"userID":     trainee.ID,
			"fullname":   trainee.UserProfile.FirstName + " " + trainee.UserProfile.LastName,
			"email":      trainee.UserProfile.PersonalEmail,
			"department": trainee.UserProfile.Department,
			"status":     status,
		}

		traineeInfoList = append(traineeInfoList, traineeInfo)
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Trainee progress list retrieved successfully",
		Data:    traineeInfoList,
	})
}

// AddListTraineeToCourse adds multiple trainees to a course
func (ctr *UserProgressController) AddListTraineeToCourse(c echo.Context) error {
	addListTraineeToCourseParams := new(param.AddListTraineeToCourseParams)
	if err := c.Bind(addListTraineeToCourseParams); err != nil {
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid request format",
		})
	}

	if _, err := valid.ValidateStruct(addListTraineeToCourseParams); err != nil {
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Validation failed: " + err.Error(),
		})
	}

	for _, traineeID := range addListTraineeToCourseParams.Trainees {
		progress := &m.UserProgress{
			UserID:             traineeID,
			CourseID:           addListTraineeToCourseParams.CourseID,
			ModulePosition:     1,
			ModuleItemPosition: 1,
			Completed:          false,
		}

		if err := ctr.UserProgressRepo.SaveUserProgress(progress); err != nil {
			ctr.Logger.Errorf("Failed to add trainee %d to course %d: %v", traineeID, addListTraineeToCourseParams.CourseID, err)
			continue
		}
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Trainees successfully added to course",
	})
}
