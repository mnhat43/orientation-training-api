package users

import (
	"fmt"
	"net/http"
	"time"

	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	resp "orientation-training-api/internal/interfaces/response"
	m "orientation-training-api/internal/models"
	"orientation-training-api/internal/platform/utils"

	valid "github.com/asaskevich/govalidator"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	cm.BaseController

	UserRepo         rp.UserRepository
	UserProgressRepo rp.UserProgressRepository
	CourseRepo       rp.CourseRepository
	ModuleRepo       rp.ModuleRepository
	ModuleItemRepo   rp.ModuleItemRepository
	QuizRepo         rp.QuizRepository
}

func NewUserController(
	logger echo.Logger,
	userRepo rp.UserRepository,
	userProgressRepo rp.UserProgressRepository,
	courseRepo rp.CourseRepository,
	moduleRepo rp.ModuleRepository,
	moduleItemRepo rp.ModuleItemRepository,
	quizRepo rp.QuizRepository,
) (ctr *UserController) {
	ctr = &UserController{
		cm.BaseController{},
		userRepo,
		userProgressRepo,
		courseRepo,
		moduleRepo,
		moduleItemRepo,
		quizRepo,
	}
	ctr.Init(logger)
	return
}

// GetLoginUser : get information user login
// Params  : echo.Context
// Returns : JSON
func (ctr *UserController) GetLoginUser(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	user, err := ctr.UserRepo.GetUserProfile(int(userID))
	if err != nil {
		if err.Error() == pg.ErrNoRows.Error() {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "User is not exists",
			})
		}

		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System error",
		})
	}

	dataResponse := map[string]interface{}{
		"id":           user.ID,
		"email":        user.Email,
		"phone_number": user.UserProfile.PhoneNumber,
		"first_name":   user.UserProfile.FirstName,
		"last_name":    user.UserProfile.LastName,
		"fullname":     user.UserProfile.FirstName + " " + user.UserProfile.LastName,
		"avatar":       user.UserProfile.Avatar,
		"birthday":     user.UserProfile.Birthday,
		"department":   user.UserProfile.Department,
		"role_id":      user.RoleID,
		"role_name":    user.Role.Name,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
		Data:    dataResponse,
	})
}

// RegisterUser : register a new user
// Params  : echo.Context
// Returns : JSON
func (ctr *UserController) Register(c echo.Context) error {
	registerParams := new(param.RegisterParams)

	if err := c.Bind(registerParams); err != nil {
		ctr.Logger.Errorf("Error binding request params: %v", err)
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid request parameters",
		})
	}

	if _, err := valid.ValidateStruct(registerParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	exists, err := ctr.UserRepo.CheckEmailExists(registerParams.Email)
	if err != nil {
		ctr.Logger.Errorf("Error checking email existence: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System error",
		})
	}

	if exists {
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Email already exists",
		})
	}

	hashedPassword := utils.GetSHA256Hash(registerParams.Password)

	var birthday time.Time
	if registerParams.Birthday != "" {
		birthday, err = time.Parse(cf.FormatDateDatabase, registerParams.Birthday)
		if err != nil {
			return c.JSON(http.StatusBadRequest, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid birthday format. Use YYYY-MM-DD",
			})
		}
	}

	newUser := m.User{
		Email:    registerParams.Email,
		Password: hashedPassword,
		RoleID:   registerParams.RoleID,
		UserProfile: m.UserProfile{
			FirstName:     registerParams.FirstName,
			LastName:      registerParams.LastName,
			PhoneNumber:   registerParams.PhoneNumber,
			PersonalEmail: registerParams.PersonnalEmail,
			Department:    registerParams.Department,
			Avatar:        registerParams.Avatar,
			Gender:        registerParams.Gender,
			Birthday:      birthday,
		},
	}

	// Create user in the database
	userID, err := ctr.UserRepo.CreateUser(newUser)
	if err != nil {
		ctr.Logger.Errorf("Error creating user: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to create user",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "User registered successfully",
		Data: map[string]interface{}{
			"user_id": userID,
		},
	})
}

// GetListTrainee retrieves all users with trainee role
// Params: echo.Context
// Returns: error
func (ctr *UserController) GetListTrainee(c echo.Context) error {
	trainees, err := ctr.UserRepo.GetUsersByRoleID(cf.EmployeeRoleID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch trainees: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch trainees",
		})
	}

	traineeList := []map[string]interface{}{}
	for _, trainee := range trainees {
		traineeInfo := map[string]interface{}{
			"userID":      trainee.ID,
			"email":       trainee.Email,
			"fullname":    trainee.UserProfile.FirstName + " " + trainee.UserProfile.LastName,
			"phoneNumber": trainee.UserProfile.PhoneNumber,
			"avatar":      trainee.UserProfile.Avatar,
			"birthday":    nil,
			"department":  trainee.UserProfile.Department,
			"gender":      cf.Gender[trainee.UserProfile.Gender],
			"joinedDate":  nil,
		}
		if !trainee.UserProfile.Birthday.IsZero() {
			traineeInfo["birthday"] = trainee.UserProfile.Birthday.Format(cf.FormatDateDatabase)
		}
		if !trainee.UserProfile.CompanyJoinedDate.IsZero() {
			traineeInfo["joinedDate"] = trainee.UserProfile.CompanyJoinedDate.Format(cf.FormatDateDatabase)
		}
		traineeList = append(traineeList, traineeInfo)
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Trainee list retrieved successfully",
		Data:    traineeList,
	})
}

// GetEmployeeOverview retrieves an overview of employees
// Params: echo.Context
// Returns: error
func (ctr *UserController) GetEmployeeOverview(c echo.Context) error {
	employees, err := ctr.UserRepo.GetUsersByRoleID(cf.EmployeeRoleID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch employees: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch employees",
		})
	}

	employeeList := []resp.EmployeeOverview{}
	for _, employee := range employees {
		var status string
		// Get user progress
		userProgresses, err := ctr.UserRepo.GetUserProgressByUserID(employee.ID)
		if err != nil && err.Error() != pg.ErrNoRows.Error() {
			ctr.Logger.Errorf("Failed to fetch user progress for user ID %d: %v", employee.ID, err)
			return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Failed to fetch user progress",
			})
		} else {
			allCompleted := true
			for _, progress := range userProgresses {
				if !progress.Completed {
					allCompleted = false
					break
				}
			}
			if len(userProgresses) == 0 {
				status = resp.StatusNotAssigned
			} else if allCompleted {
				status = resp.StatusCompleted
			} else {
				status = resp.StatusInProgress
			}
		}

		employeeInfo := resp.EmployeeOverview{
			UserID:      employee.ID,
			Fullname:    employee.UserProfile.FirstName + " " + employee.UserProfile.LastName,
			Email:       employee.Email,
			PhoneNumber: employee.UserProfile.PhoneNumber,
			Avatar:      employee.UserProfile.Avatar,
			Department:  employee.UserProfile.Department,
			Status:      status,
		}
		employeeList = append(employeeList, employeeInfo)
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Employee overview retrieved successfully",
		Data:    employeeList,
	})
}

// EmployeeDetail retrieves detailed information of an employee
// Params: echo.Context
// Returns: error
func (ctr *UserController) EmployeeDetail(c echo.Context) error {
	// Get employee detail params
	employeeDetailParams := new(param.EmployeeDetailParams)
	if err := c.Bind(employeeDetailParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid params",
			Data:    err,
		})
	}

	if _, err := valid.ValidateStruct(employeeDetailParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	// Get user info
	employee, err := ctr.UserRepo.GetUserProfile(employeeDetailParams.UserID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch employee details: %v", err)
		return c.JSON(http.StatusNotFound, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Employee not found",
		})
	}

	// Get user progress for all courses
	userProgresses, err := ctr.UserProgressRepo.GetAllUserProgressByUserID(employeeDetailParams.UserID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch user progress: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch user progress",
		})
	}
	// Prepare response
	response := buildEmployeeDetailResponse(employee, userProgresses, ctr.CourseRepo, ctr.ModuleRepo, ctr.ModuleItemRepo, ctr.QuizRepo, ctr.Logger)

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Employee details retrieved successfully",
		Data:    response,
	})
}

// buildEmployeeDetailResponse creates a structured response for employee details
func buildEmployeeDetailResponse(
	employee m.User,
	userProgresses []m.UserProgress,
	courseRepo rp.CourseRepository,
	moduleRepo rp.ModuleRepository,
	moduleItemRepo rp.ModuleItemRepository,
	quizRepo rp.QuizRepository,
	logger echo.Logger,
) resp.EmployeeDetail {
	// Prepare user info
	userInfo := resp.UserInfo{
		ID:          employee.ID,
		Fullname:    employee.UserProfile.FirstName + " " + employee.UserProfile.LastName,
		Email:       employee.Email,
		PhoneNumber: employee.UserProfile.PhoneNumber,
		Department:  employee.UserProfile.Department,
		Avatar:      employee.UserProfile.Avatar,
		JoinedDate:  employee.UserProfile.CompanyJoinedDate.Format(cf.FormatDateDatabase),
	}

	// We only need to consider courses in the user's training path (from UserProgress)
	totalCourses := len(userProgresses)

	// Initialize process info
	processInfo := []resp.CourseInfo{}
	completedCourses := 0
	totalUserScore := float64(0)
	totalMaxScore := float64(0)
	latestCompletionDate := ""

	// Create a map for quicker course lookup
	courseMap := make(map[int]m.Course)

	// Get course details for all courses in user progresses
	for _, progress := range userProgresses {
		// Try to get course details via Course relation
		if progress.Course != nil {
			courseMap[progress.CourseID] = *progress.Course
			continue
		}

		// If not available via relation, fetch individually
		course, err := courseRepo.GetCourseByID(progress.CourseID)
		if err != nil {
			logger.Errorf("Failed to fetch course %d: %v", progress.CourseID, err)
			// Create minimal course info with ID from BaseModel
			course = m.Course{}
			course.ID = progress.CourseID
			course.Title = fmt.Sprintf("Course %d", progress.CourseID)
		}
		courseMap[progress.CourseID] = course

		// Track completion date
		if progress.Completed && progress.CompletedDate > latestCompletionDate {
			latestCompletionDate = progress.CompletedDate
		}
	}

	// Process each course progress to build the process info
	for _, progress := range userProgresses {
		// Get course details from our map
		course, exists := courseMap[progress.CourseID]
		if !exists {
			logger.Errorf("Course %d not found in map, skipping", progress.CourseID)
			continue
		}
		courseInfo := resp.CourseInfo{
			CourseID:    course.ID,
			CourseTitle: course.Title,
		}

		// Calculate pending essay reviews for this course (both completed and in-progress)
		pendingReviews, err := quizRepo.GetPendingEssayReviewsCountForCourse(employee.ID, course.ID)
		if err != nil {
			logger.Errorf("Failed to get pending reviews count for course %d, user %d: %v", course.ID, employee.ID, err)
			courseInfo.PendingReviews = 0
		} else {
			courseInfo.PendingReviews = pendingReviews
		}

		if progress.Completed {
			// Course is completed
			completedCourses++
			courseInfo.Status = "completed"
			courseInfo.CompletedDate = progress.CompletedDate
			hasAssessment := progress.ReviewedBy > 0
			courseInfo.HasAssessment = hasAssessment
			if hasAssessment {
				courseInfo.Assessment = &resp.Assessment{
					PerformanceRating:  progress.PerformanceRating,
					PerformanceComment: progress.PerformanceComment,
				}

				// Get reviewer name
				if progress.Reviewer.ID > 0 {
					courseInfo.Assessment.ReviewerName = progress.Reviewer.UserProfile.FirstName + " " + progress.Reviewer.UserProfile.LastName
				}
			}

			// Calculate quiz scores for this course
			userScore, maxScore := calculateCourseQuizScores(course.ID, employee.ID, moduleRepo, moduleItemRepo, quizRepo, logger)
			courseInfo.UserScore = userScore
			courseInfo.TotalScore = maxScore

			totalUserScore += userScore
			totalMaxScore += maxScore
		} else {
			// Course is in progress
			courseInfo.Status = "in_progress"
			// Calculate progress percentage based on module positions
			progressPercent := calculateCourseProgress(progress, moduleRepo, moduleItemRepo, logger)
			courseInfo.Progress = progressPercent

			// Set current module information directly from the progress record
			if progressPercent > 0 && progress.ModulePosition > 0 {
				// Get the module directly
				module, err := moduleRepo.GetModuleByPositionAndCourse(progress.CourseID, progress.ModulePosition)
				if err == nil {
					courseInfo.CurrentModule = module.Title
				} else {
					courseInfo.CurrentModule = fmt.Sprintf("Module %d", progress.ModulePosition)
				}
			}
		}

		processInfo = append(processInfo, courseInfo)
	}

	// Build process stats
	processStats := resp.ProcessStats{
		CompletedCourses: completedCourses,
		TotalCourses:     totalCourses,
		TotalScore:       totalMaxScore,
		UserScore:        totalUserScore,
		CompletedDate:    latestCompletionDate,
	}

	// Build final response
	return resp.EmployeeDetail{
		UserInfo:     userInfo,
		ProcessStats: processStats,
		ProcessInfo:  processInfo,
	}
}

// calculateCourseQuizScores retrieves and calculates quiz scores for a specific course and user
func calculateCourseQuizScores(courseID int, userID int, moduleRepo rp.ModuleRepository, moduleItemRepo rp.ModuleItemRepository, quizRepo rp.QuizRepository, logger echo.Logger) (float64, float64) {
	userScore := float64(0)
	maxScore := float64(0)

	// Step 1: Get all quizzes belonging to the course
	// First get all modules for this course
	modules, err := moduleRepo.GetModulesByCourseID(courseID)
	if err != nil {
		logger.Errorf("Error getting modules for course %d: %v", courseID, err)
		return 0, 0
	}

	// Get all moduleIDs
	moduleIDs := make([]int, len(modules))
	for i, module := range modules {
		moduleIDs[i] = module.ID
	}

	// Get all module items for these modules with quiz relation
	var moduleItems []m.ModuleItem
	if len(moduleIDs) > 0 {
		moduleItems, err = moduleItemRepo.GetModuleItemsByModuleIDs(moduleIDs)
		if err != nil {
			logger.Errorf("Error getting module items for modules: %v", err)
			return 0, 0
		}
	}

	// For each module item of type "quiz"
	for _, item := range moduleItems {
		if item.ItemType == "quiz" && item.QuizID > 0 {
			// Get user quiz submissions
			submissions, err := quizRepo.GetQuizSubmissionsByUser(userID, item.QuizID)
			if err != nil {
				logger.Errorf("Error getting submissions for user %d, quiz %d: %v", userID, item.QuizID, err)
				continue
			}

			if len(submissions) == 0 {
				continue
			}

			// Find max attempt number
			maxAttempt := 0
			for _, submission := range submissions {
				if submission.Attempt > maxAttempt {
					maxAttempt = submission.Attempt
				}
			}

			// Calculate score for the latest attempt
			latestAttemptScore := float64(0)
			for _, submission := range submissions {
				if submission.Attempt == maxAttempt {
					latestAttemptScore += submission.Score
				}
			}

			// Use the quiz relation if available, otherwise fetch separately
			var quizTotalScore float64

			if item.Quiz != nil {
				quizTotalScore = item.Quiz.TotalScore
			} else {
				// Fallback to fetching quiz separately
				quiz, err := quizRepo.GetQuizByID(item.QuizID)
				if err != nil {
					logger.Errorf("Error getting quiz %d: %v", item.QuizID, err)
					continue
				}
				quizTotalScore = quiz.TotalScore
			}

			// Add to total scores
			userScore += latestAttemptScore
			maxScore += quizTotalScore
		}
	}

	return userScore, maxScore
}

// calculateCourseProgress calculates the percentage of course completion
func calculateCourseProgress(progress m.UserProgress, moduleRepo rp.ModuleRepository, moduleItemRepo rp.ModuleItemRepository, logger echo.Logger) int {
	// If course is marked as completed, return 100%
	if progress.Completed {
		return 100
	}

	// If no module position is set, no progress has been made
	if progress.ModulePosition <= 0 {
		return 0
	}
	// Get all modules for this course
	modules, err := moduleRepo.GetModulesByCourseID(progress.CourseID)
	if err != nil {
		logger.Errorf("Error getting modules for course %d: %v", progress.CourseID, err)
		return 0
	}

	// Count total modules and items
	totalModules := len(modules)
	if totalModules == 0 {
		return 0
	}

	// Calculate total module items
	totalItems := 0
	completedItems := 0

	// Get module IDs
	moduleIDs := make([]int, len(modules))
	for i, module := range modules {
		moduleIDs[i] = module.ID
	}

	// Get all module items
	moduleItems, err := moduleItemRepo.GetModuleItemsByModuleIDs(moduleIDs)
	if err != nil {
		logger.Errorf("Error getting module items: %v", err)
		// Fall back to estimation based on progress
		totalModules = progress.CoursePosition
		if totalModules <= 0 {
			totalModules = len(modules)
		}
	} else {
		// Count items by module
		moduleItemCounts := make(map[int]int)
		for _, item := range moduleItems {
			moduleItemCounts[item.ModuleID]++
			totalItems++
		}

		// Count completed items
		for i, module := range modules {
			position := i + 1 // Module positions are 1-indexed

			if position < progress.ModulePosition {
				// All items in this module are completed
				completedItems += moduleItemCounts[module.ID]
			} else if position == progress.ModulePosition {
				// Count completed items in current module
				completedItems += progress.ModuleItemPosition - 1
				if completedItems < 0 {
					completedItems = 0
				}
			}
		}

		// If we have valid item counts, calculate percentage directly
		if totalItems > 0 {
			percentage := (completedItems * 100) / totalItems
			if percentage >= 100 {
				percentage = 99 // Cap at 99% since not marked complete
			}
			return percentage
		}
	}
	// Fall back to estimate if we couldn't calculate based on actual items
	// Calculate the percentage completed based on module position
	completedModules := progress.ModulePosition - 1
	if completedModules < 0 {
		completedModules = 0
	}

	// Base percentage from completed modules
	percentage := (completedModules * 100) / totalModules

	// Add percentage from current module's progress
	if progress.ModulePosition > 0 && progress.ModuleItemPosition > 0 {
		// Get items for the current module
		currentModuleItems := 0
		for _, module := range modules {
			if module.Position == progress.ModulePosition {
				// Get items for this specific module
				items, err := moduleItemRepo.GetModuleItemsByModuleID(module.ID)
				if err == nil && len(items) > 0 {
					currentModuleItems = len(items)
				} else {
					// Fallback to average estimate
					currentModuleItems = 10
				}
				break
			}
		}

		// Calculate current module's contribution to percentage
		modulePercentage := (progress.ModuleItemPosition * 100) / (currentModuleItems * totalModules)
		percentage += modulePercentage
	}
	// Cap at 99% since it's not completed yet
	if percentage >= 100 {
		percentage = 99
	}

	return percentage
}
