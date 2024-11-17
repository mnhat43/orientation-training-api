package users

import (
	"net/http"

	m "orientation-training-api/internal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg/v9"

	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"

	"github.com/labstack/echo/v4"
)

type UserMiddleware struct {
	cm.AppRepository

	UserRepo rp.UserRepository
}

func NewUserMiddleware(logger echo.Logger, userRepo rp.UserRepository) (userMw *UserMiddleware) {
	userMw = &UserMiddleware{cm.AppRepository{}, userRepo}
	userMw.Init(logger)
	return
}

func (userMw *UserMiddleware) InitUserProfile(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := getUserIDWithToken(c)
		userProfile, err := userMw.UserRepo.GetUserProfile(userID)

		// Check user exists
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

		// add info user profile to echo context(global)
		c.Set("user_profile", userProfile)

		return next(c)
	}
}

func (userMw *UserMiddleware) CheckAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userProfile := c.Get("user_profile").(m.User)

		// Check user is admin
		if userProfile.RoleID != cf.AdminRoleID {
			return c.JSON(http.StatusMethodNotAllowed, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Only admin can do this.",
			})
		}
		return next(c)
	}
}

func (userMw *UserMiddleware) CheckGeneralManager(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userProfile := c.Get("user_profile").(m.User)

		// Check user is general manager
		if userProfile.RoleID != cf.GeneralManagerRoleID {
			return c.JSON(http.StatusMethodNotAllowed, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Only general manager can do this.",
			})
		}
		return next(c)
	}
}

func (userMw *UserMiddleware) CheckManager(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userProfile := c.Get("user_profile").(m.User)

		// Check user is manager
		if userProfile.RoleID != cf.ManagerRoleID {
			return c.JSON(http.StatusMethodNotAllowed, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Only manager can do this.",
			})
		}
		return next(c)
	}
}

func (userMw *UserMiddleware) CheckAllManager(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userProfile := c.Get("user_profile").(m.User)

		// Check user is admin
		if userProfile.RoleID != cf.GeneralManagerRoleID && userProfile.RoleID != cf.ManagerRoleID {
			return c.JSON(http.StatusMethodNotAllowed, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Only managers can do this.",
			})
		}
		return next(c)
	}
}

func getUserIDWithToken(c echo.Context) int {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)

	return int(userID)
}
