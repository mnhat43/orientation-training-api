package users

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"

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

// GetListTrainee retrieves all users with trainee role
// Params: echo.Context
// Returns: error
func (ctr *UserController) GetListTrainee(c echo.Context) error {
	trainees, err := ctr.UserRepo.GetUsersByRoleID(cf.TraineeRoleID)
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
