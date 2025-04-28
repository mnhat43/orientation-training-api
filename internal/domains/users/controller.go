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
