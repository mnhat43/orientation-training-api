package users

import (
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

	UserRepo rp.UserRepository
}

func NewUserController(
	logger echo.Logger,
	userRepo rp.UserRepository,
) (ctr *UserController) {
	ctr = &UserController{
		cm.BaseController{},
		userRepo,
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
