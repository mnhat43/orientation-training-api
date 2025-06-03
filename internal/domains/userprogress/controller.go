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

	targetUserID := userProfile.ID

	// If user is a manager and a specific userID is provided, use that instead
	if userProfile.RoleID == cf.ManagerRoleID && updateUserProgressParams.UserID > 0 {
		targetUserID = updateUserProgressParams.UserID
	}

	userProgress := &m.UserProgress{
		UserID:         targetUserID,
		CourseID:       updateUserProgressParams.CourseID,
		CoursePosition: updateUserProgressParams.CoursePosition,
	}

	if updateUserProgressParams.Completed {
		currentProgress, err := ctr.UserProgressRepo.GetSingleUserProgress(targetUserID, updateUserProgressParams.CourseID)
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
		userProgress.CompletedDate = updateUserProgressParams.CompletedDate
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
			"course_position":      userProgress.CoursePosition,
			"module_position":      userProgress.ModulePosition,
			"module_item_position": userProgress.ModuleItemPosition,
			"completed":            userProgress.Completed,
		},
	})
}

// GetAllUserProgressByUserID is now modified to get progress for ALL courses of a user
// Params: echo.Context
// Returns: error
func (ctr *UserProgressController) GetAllUserProgressByUserID(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	getUserProgressParams := new(param.GetAllUserProgressParams)

	if err := c.Bind(getUserProgressParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	targetUserID := userProfile.ID

	if userProfile.RoleID == cf.ManagerRoleID && getUserProgressParams.UserID > 0 {
		targetUserID = getUserProgressParams.UserID
	}

	userProgressList, err := ctr.UserProgressRepo.GetAllUserProgressByUserID(targetUserID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch user progress list: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch user progress list",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "User progress for all courses retrieved successfully",
		Data:    userProgressList,
	})
}

// GetSingleCourseProgress retrieves a user's progress for a specific course
// Params: echo.Context
// Returns: error
func (ctr *UserProgressController) GetSingleCourseProgress(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	getSingleCourseProgressParams := new(param.GetSingleCourseProgressParams)

	if err := c.Bind(getSingleCourseProgressParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(getSingleCourseProgressParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	targetUserID := userProfile.ID

	if userProfile.RoleID == cf.ManagerRoleID && getSingleCourseProgressParams.UserID > 0 {
		targetUserID = getSingleCourseProgressParams.UserID
	}

	userProgress, err := ctr.UserProgressRepo.GetSingleUserProgress(targetUserID, getSingleCourseProgressParams.CourseID)
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

// AddUserProgress creates new user progress records for multiple courses
// Only managers can add progress
// Params: echo.Context
// Returns: error
func (ctr *UserProgressController) AddUserProgress(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)

	if userProfile.RoleID != cf.ManagerRoleID {
		return c.JSON(http.StatusForbidden, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Only managers can add progress for users",
		})
	}

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

	targetUserID := createUserProgressParams.UserID
	successCourses := []int{}

	for i, courseID := range createUserProgressParams.CourseIDs {
		coursePosition := i + 1

		existingProgress, err := ctr.UserProgressRepo.GetSingleUserProgress(targetUserID, courseID)
		if err == nil && existingProgress.ID != 0 {
			ctr.Logger.Infof("Progress already exists for user %d and course %d", targetUserID, courseID)
			continue
		}

		userProgress := &m.UserProgress{
			UserID:             targetUserID,
			CourseID:           courseID,
			CoursePosition:     coursePosition,
			ModulePosition:     1,
			ModuleItemPosition: 1,
			Completed:          false,
		}

		err = ctr.UserProgressRepo.SaveUserProgress(userProgress)
		if err != nil {
			ctr.Logger.Errorf("Failed to create user progress for course %d: %v", courseID, err)
			continue
		}

		successCourses = append(successCourses, courseID)
	}

	if len(successCourses) == 0 {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "No courses were added. Courses may already be assigned or an error occurred.",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Progress created successfully",
		Data:    nil,
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

	trainees, err := ctr.UserRepo.GetUsersByRoleID(cf.EmployeeRoleID)
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

// ReviewProgress allows managers to review the progress of trainees
func (ctr *UserProgressController) ReviewProgress(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)

	reviewProgressParams := new(param.ReviewProgressParams)

	if err := c.Bind(reviewProgressParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(reviewProgressParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}
	_, err := ctr.UserProgressRepo.GetSingleUserProgress(reviewProgressParams.UserID, reviewProgressParams.CourseID)
	if err != nil {
		ctr.Logger.Errorf("User progress not found: %v", err)
		return c.JSON(http.StatusNotFound, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "User progress not found",
		})
	}

	err = ctr.UserProgressRepo.ReviewUserProgress(
		reviewProgressParams.UserID,
		reviewProgressParams.CourseID,
		reviewProgressParams.PerformanceRating,
		reviewProgressParams.PerformanceComment,
		userProfile.ID,
	)

	if err != nil {
		ctr.Logger.Errorf("Failed to review user progress: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to review user progress",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "User progress reviewed successfully",
	})
}
