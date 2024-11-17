package auth

import (
	"net/http"

	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	"orientation-training-api/internal/platform/utils"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/go-pg/pg"

	"github.com/dgrijalva/jwt-go"

	// "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type AuthController struct {
	cm.BaseController

	UserRepo rp.UserRepository
}

func NewAuthController(logger echo.Logger, userRepo rp.UserRepository) (ctr *AuthController) {
	ctr = &AuthController{cm.BaseController{}, userRepo}
	ctr.Init(logger)
	return
}

func createTokenLogin(userID int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	claims["exp"] = utils.TimeNowUTC().Add(time.Hour * 72).Unix()

	keyTokenAuth := utils.GetKeyToken()
	t, err := token.SignedString([]byte(keyTokenAuth))

	return t, err
}

func (ctr *AuthController) Login(c echo.Context) error {
	email := c.FormValue("email")
	password := utils.GetSHA256Hash(c.FormValue("password"))

	if !valid.IsEmail(email) {
		return c.JSON(http.StatusUnauthorized, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid email",
		})
	}

	// get ID user login in DB
	idUserLogin, err := ctr.UserRepo.GetLoginUserID(email, password)

	if err != nil {
		if err.Error() == pg.ErrNoRows.Error() {
			//select no rows in database
			return c.JSON(http.StatusUnauthorized, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "User is not exist or password wrong",
			})
		}

		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System error",
		})
	}

	err = ctr.UserRepo.UpdateLastLogin(idUserLogin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "System error",
		})
	}

	tokenLogin, err := createTokenLogin(idUserLogin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to create token",
		})
	}

	objToken := map[string]string{
		"token": tokenLogin,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
		Data:    objToken,
	})
}

func (ctr *AuthController) Logout(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	claims["exp"] = utils.TimeNowUTC()

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Success",
	})
}
