package users

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	resp "orientation-training-api/internal/interfaces/response"
	m "orientation-training-api/internal/models"
	gc "orientation-training-api/internal/platform/cloud"
	"orientation-training-api/internal/platform/utils"

	valid "github.com/asaskevich/govalidator"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg"
	"github.com/labstack/echo/v4"
)

type UserController struct {
	cm.BaseController

	UserRepo               rp.UserRepository
	UserProgressRepo       rp.UserProgressRepository
	CourseRepo             rp.CourseRepository
	ModuleRepo             rp.ModuleRepository
	ModuleItemRepo         rp.ModuleItemRepository
	QuizRepo               rp.QuizRepository
	CourseSkillKeywordRepo rp.CourseSkillKeywordRepository
	cloud                  gc.StorageUtility
}

func NewUserController(
	logger echo.Logger,
	userRepo rp.UserRepository,
	userProgressRepo rp.UserProgressRepository,
	courseRepo rp.CourseRepository,
	moduleRepo rp.ModuleRepository,
	moduleItemRepo rp.ModuleItemRepository,
	quizRepo rp.QuizRepository,
	courseSkillKeywordRepo rp.CourseSkillKeywordRepository,
	cloud gc.StorageUtility,

) (ctr *UserController) {
	ctr = &UserController{
		cm.BaseController{},
		userRepo,
		userProgressRepo,
		courseRepo,
		moduleRepo,
		moduleItemRepo,
		quizRepo,
		courseSkillKeywordRepo,
		cloud,
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

	if registerParams.Avatar != "" {
		parts := strings.SplitN(registerParams.Avatar, ",", 2)
		if len(parts) != 2 {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid Avatar Format",
			})
		}

		mimeType := parts[0]
		base64Data := parts[1]

		formatImageAvatar := ""
		if strings.HasPrefix(mimeType, "data:image/") {
			formatImageAvatar = strings.TrimPrefix(mimeType, "data:image/")
			formatImageAvatar = strings.Split(formatImageAvatar, ";")[0]
		}

		if formatImageAvatar == "" {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid Image Format",
			})
		}

		if _, check := utils.FindStringInArray(cf.AllowFormatImageList, formatImageAvatar); !check {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "The Avatar field must be an image",
			})
		}

		millisecondTimeNow := int(time.Now().UnixNano() / int64(time.Millisecond))
		nameAvatar := fmt.Sprintf("%d_%d.%s", 99, millisecondTimeNow, formatImageAvatar)

		err := ctr.cloud.UploadFileToCloud(base64Data, nameAvatar, cf.AvatarFolderGCS)
		if err != nil {
			ctr.Logger.Error(err)
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Upload Avatar error",
			})
		}

		registerParams.Avatar = ctr.cloud.GetURL(nameAvatar, cf.AvatarFolderGCS)
	}

	newUser := m.User{
		Email:    registerParams.Email,
		Password: hashedPassword,
		RoleID:   registerParams.RoleID,
		UserProfile: m.UserProfile{
			FirstName:         registerParams.FirstName,
			LastName:          registerParams.LastName,
			PhoneNumber:       registerParams.PhoneNumber,
			PersonalEmail:     registerParams.PersonnalEmail,
			Department:        registerParams.Department,
			Avatar:            registerParams.Avatar,
			Gender:            registerParams.Gender,
			CompanyJoinedDate: registerParams.CompanyJoinedDate,
			Birthday:          registerParams.Birthday,
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

// GetListTrainee retrieves all users with trainee role who don't have any assigned courses
// Params: echo.Context
// Returns: error
func (ctr *UserController) GetListTrainee(c echo.Context) error {
	trainees, err := ctr.UserRepo.GetUsersWithoutProgress(cf.EmployeeRoleID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch trainees without assigned courses: %v", err)
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
			"birthday":    trainee.UserProfile.Birthday,
			"department":  trainee.UserProfile.Department,
			"gender":      cf.Gender[trainee.UserProfile.Gender],
			"joinedDate":  trainee.UserProfile.CompanyJoinedDate,
		}
		traineeList = append(traineeList, traineeInfo)
	}
	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Unassigned trainee list retrieved successfully",
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
		// Initialize with empty skill keywords list
		skillKeywords := []string{}

		// For employees with InProgress or Completed status, get skill keywords from completed courses
		if status == resp.StatusInProgress || status == resp.StatusCompleted {
			// Find completed courses
			var completedCourseIDs []int
			for _, progress := range userProgresses {
				if progress.Completed {
					completedCourseIDs = append(completedCourseIDs, progress.CourseID)
				}
			}

			// Get skill keywords for completed courses
			if len(completedCourseIDs) > 0 {
				skillKeywordsMap := make(map[string]bool) // To avoid duplicates
				for _, courseID := range completedCourseIDs {
					keywords, err := ctr.CourseSkillKeywordRepo.GetSkillKeywordsByCourseID(courseID)
					if err != nil {
						ctr.Logger.Warnf("Failed to fetch skill keywords for course ID %d: %v", courseID, err)
						continue
					}

					for _, keyword := range keywords {
						skillKeywordsMap[keyword.Name] = true
					}
				}

				// Convert map keys to slice
				for keyword := range skillKeywordsMap {
					skillKeywords = append(skillKeywords, keyword)
				}
			}
		}

		employeeInfo := resp.EmployeeOverview{
			UserID:        employee.ID,
			Fullname:      employee.UserProfile.FirstName + " " + employee.UserProfile.LastName,
			Email:         employee.Email,
			PhoneNumber:   employee.UserProfile.PhoneNumber,
			Avatar:        employee.UserProfile.Avatar,
			Department:    employee.UserProfile.Department,
			Status:        status,
			SkillKeywords: skillKeywords,
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

	employee, err := ctr.UserRepo.GetUserProfile(employeeDetailParams.UserID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch employee details: %v", err)
		return c.JSON(http.StatusNotFound, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Employee not found",
		})
	}
	userProgresses, err := ctr.UserProgressRepo.GetAllUserProgressByUserID(employeeDetailParams.UserID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch user progress: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch user progress",
		})
	}
	response := buildEmployeeDetailResponse(employee, userProgresses, ctr.CourseRepo, ctr.ModuleRepo, ctr.ModuleItemRepo, ctr.QuizRepo, ctr.CourseSkillKeywordRepo, ctr.Logger)

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Employee details retrieved successfully",
		Data:    response,
	})
}

func buildEmployeeDetailResponse(
	employee m.User,
	userProgresses []m.UserProgress,
	courseRepo rp.CourseRepository,
	moduleRepo rp.ModuleRepository,
	moduleItemRepo rp.ModuleItemRepository,
	quizRepo rp.QuizRepository,
	courseSkillKeywordRepo rp.CourseSkillKeywordRepository,
	logger echo.Logger,
) resp.EmployeeDetail {
	userInfo := resp.UserInfo{
		ID:          employee.ID,
		Fullname:    employee.UserProfile.FirstName + " " + employee.UserProfile.LastName,
		Email:       employee.Email,
		PhoneNumber: employee.UserProfile.PhoneNumber,
		Department:  employee.UserProfile.Department,
		Avatar:      employee.UserProfile.Avatar,
		JoinedDate:  employee.UserProfile.CompanyJoinedDate,
	}

	totalCourses := len(userProgresses)
	processInfo := []resp.CourseInfo{}
	completedCourses := 0
	totalUserScore := float64(0)
	totalMaxScore := float64(0)
	latestCompletionDate := ""

	totalPerformanceRating := float64(0)
	ratedCourses := 0

	courseMap := make(map[int]m.Course)

	// Map to store all user skills (from completed courses)
	allUserSkills := make(map[string]bool)

	for _, progress := range userProgresses {
		if progress.Course != nil {
			courseMap[progress.CourseID] = *progress.Course
			continue
		}

		course, err := courseRepo.GetCourseByID(progress.CourseID)
		if err != nil {
			logger.Errorf("Failed to fetch course %d: %v", progress.CourseID, err)
			course = m.Course{}
			course.ID = progress.CourseID
			course.Title = fmt.Sprintf("Course %d", progress.CourseID)
		}
		courseMap[progress.CourseID] = course

		if progress.Completed && progress.CompletedDate > latestCompletionDate {
			latestCompletionDate = progress.CompletedDate
		}
	}

	for _, progress := range userProgresses {
		course, exists := courseMap[progress.CourseID]
		if !exists {
			logger.Errorf("Course %d not found in map, skipping", progress.CourseID)
			continue
		}
		courseInfo := resp.CourseInfo{
			CourseID:    course.ID,
			CourseTitle: course.Title,
		}

		// Get skill keywords for this course
		skillKeywords := []string{}
		keywords, err := courseSkillKeywordRepo.GetSkillKeywordsByCourseID(course.ID)
		if err != nil {
			logger.Warnf("Failed to fetch skill keywords for course ID %d: %v", course.ID, err)
		} else {
			for _, keyword := range keywords {
				skillKeywords = append(skillKeywords, keyword.Name)
			}
			courseInfo.SkillKeywords = skillKeywords
		}

		pendingReviews, err := quizRepo.GetPendingEssayReviewsCountForCourse(employee.ID, course.ID)
		if err != nil {
			logger.Errorf("Failed to get pending reviews count for course %d, user %d: %v", course.ID, employee.ID, err)
			courseInfo.PendingReviews = 0
		} else {
			courseInfo.PendingReviews = pendingReviews
		}
		if progress.Completed {
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

				if progress.ReviewedBy > 0 {
					courseInfo.Assessment.ReviewerName = progress.Reviewer.UserProfile.FirstName + " " + progress.Reviewer.UserProfile.LastName
				}

				if progress.PerformanceRating > 0 {
					totalPerformanceRating += progress.PerformanceRating
					ratedCourses++
				}
			}

			userScore, maxScore := calculateCourseQuizScores(course.ID, employee.ID, moduleRepo, moduleItemRepo, quizRepo, logger)
			courseInfo.UserScore = userScore
			courseInfo.TotalScore = maxScore

			totalUserScore += userScore
			totalMaxScore += maxScore

			// Add skills from completed courses to the user's skill set
			for _, keyword := range skillKeywords {
				allUserSkills[keyword] = true
			}
		} else {
			courseInfo.Status = "in_progress"
			progressPercent := calculateCourseProgress(progress, moduleRepo, moduleItemRepo, logger)
			courseInfo.Progress = progressPercent

			if progressPercent > 0 && progress.ModulePosition > 0 {
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
	averagePerformanceRating := float64(0)
	if ratedCourses > 0 {
		averagePerformanceRating = totalPerformanceRating / float64(ratedCourses)
		averagePerformanceRating = float64(int(averagePerformanceRating*100+0.5)) / 100
	}

	// Convert accumulated user skills to a slice
	userSkills := []string{}
	for skill := range allUserSkills {
		userSkills = append(userSkills, skill)
	}

	processStats := resp.ProcessStats{
		CompletedCourses:         completedCourses,
		TotalCourses:             totalCourses,
		TotalScore:               totalMaxScore,
		UserScore:                totalUserScore,
		CompletedDate:            latestCompletionDate,
		AveragePerformanceRating: averagePerformanceRating,
		UserSkills:               userSkills,
	}

	return resp.EmployeeDetail{
		UserInfo:     userInfo,
		ProcessStats: processStats,
		ProcessInfo:  processInfo,
	}
}

func calculateCourseQuizScores(courseID int, userID int, moduleRepo rp.ModuleRepository, moduleItemRepo rp.ModuleItemRepository, quizRepo rp.QuizRepository, logger echo.Logger) (float64, float64) {
	userScore := float64(0)
	maxScore := float64(0)

	modules, err := moduleRepo.GetModulesByCourseID(courseID)
	if err != nil {
		logger.Errorf("Error getting modules for course %d: %v", courseID, err)
		return 0, 0
	}

	moduleIDs := make([]int, len(modules))
	for i, module := range modules {
		moduleIDs[i] = module.ID
	}

	var moduleItems []m.ModuleItem
	if len(moduleIDs) > 0 {
		moduleItems, err = moduleItemRepo.GetModuleItemsByModuleIDs(moduleIDs)
		if err != nil {
			logger.Errorf("Error getting module items for modules: %v", err)
			return 0, 0
		}
	}

	for _, item := range moduleItems {
		if item.ItemType == "quiz" && item.QuizID > 0 {
			submissions, err := quizRepo.GetQuizSubmissionsByUser(userID, item.QuizID)
			if err != nil {
				logger.Errorf("Error getting submissions for user %d, quiz %d: %v", userID, item.QuizID, err)
				continue
			}

			if len(submissions) == 0 {
				continue
			}

			maxAttempt := 0
			for _, submission := range submissions {
				if submission.Attempt > maxAttempt {
					maxAttempt = submission.Attempt
				}
			}

			latestAttemptScore := float64(0)
			for _, submission := range submissions {
				if submission.Attempt == maxAttempt {
					latestAttemptScore += submission.Score
				}
			}

			var quizTotalScore float64

			if item.Quiz != nil {
				quizTotalScore = item.Quiz.TotalScore
			} else {
				quiz, err := quizRepo.GetQuizByID(item.QuizID)
				if err != nil {
					logger.Errorf("Error getting quiz %d: %v", item.QuizID, err)
					continue
				}
				quizTotalScore = quiz.TotalScore
			}

			userScore += latestAttemptScore
			maxScore += quizTotalScore
		}
	}

	return userScore, maxScore
}

func calculateCourseProgress(progress m.UserProgress, moduleRepo rp.ModuleRepository, moduleItemRepo rp.ModuleItemRepository, logger echo.Logger) int {
	if progress.Completed {
		return 100
	}

	if progress.ModulePosition <= 0 {
		return 0
	}
	modules, err := moduleRepo.GetModulesByCourseID(progress.CourseID)
	if err != nil {
		logger.Errorf("Error getting modules for course %d: %v", progress.CourseID, err)
		return 0
	}

	totalModules := len(modules)
	if totalModules == 0 {
		return 0
	}

	totalItems := 0
	completedItems := 0

	moduleIDs := make([]int, len(modules))
	for i, module := range modules {
		moduleIDs[i] = module.ID
	}

	moduleItems, err := moduleItemRepo.GetModuleItemsByModuleIDs(moduleIDs)
	if err != nil {
		logger.Errorf("Error getting module items: %v", err)
		totalModules = progress.CoursePosition
		if totalModules <= 0 {
			totalModules = len(modules)
		}
	} else {
		moduleItemCounts := make(map[int]int)
		for _, item := range moduleItems {
			moduleItemCounts[item.ModuleID]++
			totalItems++
		}

		for i, module := range modules {
			position := i + 1

			if position < progress.ModulePosition {
				completedItems += moduleItemCounts[module.ID]
			} else if position == progress.ModulePosition {
				completedItems += progress.ModuleItemPosition - 1
				if completedItems < 0 {
					completedItems = 0
				}
			}
		}

		if totalItems > 0 {
			percentage := (completedItems * 100) / totalItems
			if percentage >= 100 {
				percentage = 99
			}
			return percentage
		}
	}
	completedModules := progress.ModulePosition - 1
	if completedModules < 0 {
		completedModules = 0
	}

	percentage := (completedModules * 100) / totalModules

	if progress.ModulePosition > 0 && progress.ModuleItemPosition > 0 {
		currentModuleItems := 0
		for _, module := range modules {
			if module.Position == progress.ModulePosition {
				items, err := moduleItemRepo.GetModuleItemsByModuleID(module.ID)
				err = nil
				if err == nil && len(items) > 0 {
					currentModuleItems = len(items)
				} else {
					currentModuleItems = 10
				}
				break
			}
		}

		modulePercentage := (progress.ModuleItemPosition * 100) / (currentModuleItems * totalModules)
		percentage += modulePercentage
	}
	if percentage >= 100 {
		percentage = 99
	}

	return percentage
}

// UpdateProfile updates a user's profile information
// Params: echo.Context
// Returns: error
func (ctr *UserController) UpdateProfile(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	userID := userProfile.ID

	updateProfileParams := new(param.UpdateProfileParams)
	if err := c.Bind(updateProfileParams); err != nil {
		ctr.Logger.Errorf("Error binding request params: %v", err)
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid request parameters",
		})
	}

	if _, err := valid.ValidateStruct(updateProfileParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	if updateProfileParams.Avatar != "" {
		parts := strings.SplitN(updateProfileParams.Avatar, ",", 2)
		if len(parts) != 2 {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid Avatar Format",
			})
		}

		mimeType := parts[0]
		base64Data := parts[1]

		formatImageAvatar := ""
		if strings.HasPrefix(mimeType, "data:image/") {
			formatImageAvatar = strings.TrimPrefix(mimeType, "data:image/")
			formatImageAvatar = strings.Split(formatImageAvatar, ";")[0]
		}

		if formatImageAvatar == "" {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid Image Format",
			})
		}

		if _, check := utils.FindStringInArray(cf.AllowFormatImageList, formatImageAvatar); !check {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "The Avatar field must be an image",
			})
		}

		millisecondTimeNow := int(time.Now().UnixNano() / int64(time.Millisecond))
		nameAvatar := fmt.Sprintf("%d_%d.%s", userID, millisecondTimeNow, formatImageAvatar)

		err := ctr.cloud.UploadFileToCloud(base64Data, nameAvatar, cf.AvatarFolderGCS)
		if err != nil {
			ctr.Logger.Error(err)
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Upload Avatar error",
			})
		}

		updateProfileParams.Avatar = ctr.cloud.GetURL(nameAvatar, cf.AvatarFolderGCS)
	}

	err := ctr.UserRepo.UpdateUserProfile(userID, updateProfileParams)
	if err != nil {
		ctr.Logger.Errorf("Error updating user profile: %v", err)

		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to update profile",
		})
	}

	updatedUser, err := ctr.UserRepo.GetUserProfile(int(userID))
	if err != nil {
		ctr.Logger.Errorf("Error getting updated user profile: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.SuccessResponseCode,
			Message: "Profile updated successfully, but failed to fetch updated data",
		})
	}

	dataResponse := map[string]interface{}{
		"id":           updatedUser.ID,
		"email":        updatedUser.Email,
		"phone_number": updatedUser.UserProfile.PhoneNumber,
		"first_name":   updatedUser.UserProfile.FirstName,
		"last_name":    updatedUser.UserProfile.LastName,
		"fullname":     updatedUser.UserProfile.FirstName + " " + updatedUser.UserProfile.LastName,
		"avatar":       updatedUser.UserProfile.Avatar,
		"birthday":     updatedUser.UserProfile.Birthday,
		"role_id":      updatedUser.RoleID,
		"role_name":    updatedUser.Role.Name,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Profile updated successfully",
		Data:    dataResponse,
	})
}

// ChangePassword allows a user to change their password
// Params: echo.Context
// Returns: error
func (ctr *UserController) ChangePassword(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)
	userID := userProfile.ID

	changePasswordParams := new(param.ChangePasswordParams)
	if err := c.Bind(changePasswordParams); err != nil {
		ctr.Logger.Errorf("Error binding request params: %v", err)
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid request parameters",
		})
	}

	if _, err := valid.ValidateStruct(changePasswordParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	if changePasswordParams.NewPassword != changePasswordParams.ConfirmPassword {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "New password and confirmation password do not match",
		})
	}

	user, err := ctr.UserRepo.GetUserProfile(userID)
	if err != nil {
		ctr.Logger.Errorf("Error getting user profile: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to retrieve user information",
		})
	}

	currentPasswordHash := utils.GetSHA256Hash(changePasswordParams.CurrentPassword)
	if currentPasswordHash != user.Password {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Current password is incorrect",
		})
	}

	newPasswordHash := utils.GetSHA256Hash(changePasswordParams.NewPassword)

	err = ctr.UserRepo.UpdatePassword(userID, newPasswordHash)
	if err != nil {
		ctr.Logger.Errorf("Error updating password: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to update password",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Password updated successfully",
	})
}

// GetAllUsers retrieves all users except those with admin role
// Params: echo.Context
// Returns: error
func (ctr *UserController) GetAllUsers(c echo.Context) error {
	users, err := ctr.UserRepo.GetAllUsersExceptRole(cf.AdminRoleID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch users: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch users",
		})
	}
	userList := []map[string]interface{}{}
	for _, user := range users {
		userInfo := map[string]interface{}{
			"user_id":             user.ID,
			"first_name":          user.UserProfile.FirstName,
			"last_name":           user.UserProfile.LastName,
			"fullname":            user.UserProfile.FirstName + " " + user.UserProfile.LastName,
			"email":               user.Email,
			"role_name":           user.Role.Name,
			"role_id":             user.RoleID,
			"department":          user.UserProfile.Department,
			"phone_number":        user.UserProfile.PhoneNumber,
			"birthday":            user.UserProfile.Birthday,
			"gender":              user.UserProfile.Gender,
			"company_joined_date": user.UserProfile.CompanyJoinedDate,
			"last_login":          user.LastLoginTime,
			"created_at":          user.CreatedAt,
			"avatar":              user.UserProfile.Avatar,
		}
		userList = append(userList, userInfo)
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Users retrieved successfully",
		Data:    userList,
	})
}

// AdminUpdateUser allows an admin to update another user's information
// Params: echo.Context
// Returns: error
func (ctr *UserController) AdminUpdateUser(c echo.Context) error {
	updateParams := new(param.AdminUpdateUserParams)
	if err := c.Bind(updateParams); err != nil {
		ctr.Logger.Errorf("Error binding request params: %v", err)
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid request parameters",
		})
	}

	if _, err := valid.ValidateStruct(updateParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}
	// Check if user exists
	_, err := ctr.UserRepo.GetUserProfile(updateParams.UserID)
	if err != nil {
		if err.Error() == pg.ErrNoRows.Error() {
			return c.JSON(http.StatusNotFound, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "User not found",
			})
		}
		ctr.Logger.Errorf("Error checking user existence: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System error",
		})
	}

	// Process avatar if it's in base64 format
	if updateParams.Avatar != "" && strings.HasPrefix(updateParams.Avatar, "data:") {
		parts := strings.SplitN(updateParams.Avatar, ",", 2)
		if len(parts) != 2 {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid Avatar Format",
			})
		}

		mimeType := parts[0]
		base64Data := parts[1]

		formatImageAvatar := ""
		if strings.HasPrefix(mimeType, "data:image/") {
			formatImageAvatar = strings.TrimPrefix(mimeType, "data:image/")
			formatImageAvatar = strings.Split(formatImageAvatar, ";")[0]
		}

		if formatImageAvatar == "" {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Invalid Image Format",
			})
		}

		if _, check := utils.FindStringInArray(cf.AllowFormatImageList, formatImageAvatar); !check {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "The Avatar field must be an image",
			})
		}

		millisecondTimeNow := int(time.Now().UnixNano() / int64(time.Millisecond))
		nameAvatar := fmt.Sprintf("%d_%d.%s", updateParams.UserID, millisecondTimeNow, formatImageAvatar)

		err := ctr.cloud.UploadFileToCloud(base64Data, nameAvatar, cf.AvatarFolderGCS)
		if err != nil {
			ctr.Logger.Error(err)
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Upload Avatar error",
			})
		}

		updateParams.Avatar = ctr.cloud.GetURL(nameAvatar, cf.AvatarFolderGCS)
	}

	// Update user information
	err = ctr.UserRepo.AdminUpdateUser(updateParams.UserID, updateParams)
	if err != nil {
		ctr.Logger.Errorf("Error updating user: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to update user information",
		})
	}

	// Get updated user data
	updatedUser, err := ctr.UserRepo.GetUserProfile(updateParams.UserID)
	if err != nil {
		ctr.Logger.Warnf("User updated but failed to get updated data: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.SuccessResponseCode,
			Message: "User information updated successfully, but failed to fetch updated data",
		})
	}
	// Return updated user data
	dataResponse := map[string]interface{}{
		"user_id":             updatedUser.ID,
		"email":               updatedUser.Email,
		"first_name":          updatedUser.UserProfile.FirstName,
		"last_name":           updatedUser.UserProfile.LastName,
		"fullname":            updatedUser.UserProfile.FirstName + " " + updatedUser.UserProfile.LastName,
		"phone_number":        updatedUser.UserProfile.PhoneNumber,
		"birthday":            updatedUser.UserProfile.Birthday,
		"department":          updatedUser.UserProfile.Department,
		"gender":              updatedUser.UserProfile.Gender,
		"company_joined_date": updatedUser.UserProfile.CompanyJoinedDate,
		"avatar":              updatedUser.UserProfile.Avatar,
		"role_id":             updatedUser.RoleID,
		"role_name":           updatedUser.Role.Name,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "User information updated successfully",
		Data:    dataResponse,
	})
}

// DeleteUser allows an admin to delete a user from the system (soft delete)
// Params: echo.Context
// Returns: error
func (ctr *UserController) DeleteUser(c echo.Context) error {
	deleteUserParams := new(param.DeleteUserParams)

	if err := c.Bind(deleteUserParams); err != nil {
		ctr.Logger.Errorf("Error binding request params: %v", err)
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid request parameters",
		})
	}

	if _, err := valid.ValidateStruct(deleteUserParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	// Check if user exists
	targetUser, err := ctr.UserRepo.GetUserProfile(deleteUserParams.UserId)
	if err != nil {
		if err.Error() == pg.ErrNoRows.Error() {
			return c.JSON(http.StatusNotFound, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "User not found",
			})
		}
		ctr.Logger.Errorf("Error checking user existence: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System error",
		})
	}

	// Don't allow deletion of the current user
	userProfile := c.Get("user_profile").(m.User)
	if userProfile.ID == deleteUserParams.UserId {
		return c.JSON(http.StatusBadRequest, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Cannot delete your own account",
		})
	}

	// Delete the user
	err = ctr.UserRepo.DeleteUser(deleteUserParams.UserId)
	if err != nil {
		ctr.Logger.Errorf("Error deleting user: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to delete user",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: fmt.Sprintf("User %s (%s) deleted successfully", targetUser.Email, targetUser.UserProfile.FirstName+" "+targetUser.UserProfile.LastName),
		Data: map[string]interface{}{
			"user_id": deleteUserParams.UserId,
		},
	})
}
