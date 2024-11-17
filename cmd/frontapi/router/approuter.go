package router

import (
	"orientation-training-api/internal/domains/auth"
	c "orientation-training-api/internal/domains/courses"
	u "orientation-training-api/internal/domains/users"
	"orientation-training-api/internal/platform/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type AppRouter struct {
	authCtr   *auth.AuthController
	userCtr   *u.UserController
	courseCtr *c.CourseController
	// adminCtr *ad.Controller

	userMw *u.UserMiddleware
}

func NewAppRouter(logger echo.Logger) (r *AppRouter) {
	userRepo := u.NewPgUserRepository(logger)
	courseRepo := c.NewPgCourseRepository(logger)
	// adminRepo := ad.NewPgAdminRepository(logger)

	r = &AppRouter{
		authCtr:   auth.NewAuthController(logger, userRepo),
		userCtr:   u.NewUserController(logger, userRepo),
		courseCtr: c.NewCourseController(logger, courseRepo),
		// adminCtr: ad.NewAdminController(),

		// userMw: u.NewUserMiddleware(),
	}

	return
}

func (r *AppRouter) UserRoute(g *echo.Group) {
	keyTokenAuth := utils.GetKeyToken()
	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(keyTokenAuth),
	})

	g.GET("/getuser", r.userCtr.GetLoginUser, isLoggedIn)

}

func (r *AppRouter) AuthRoute(g *echo.Group) {
	keyTokenAuth := utils.GetKeyToken()
	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(keyTokenAuth),
	})

	g.POST("/login", r.authCtr.Login)
	g.GET("/logout", r.authCtr.Logout, isLoggedIn)

}

func (r *AppRouter) CourseRoute(g *echo.Group) {
	keyTokenAuth := utils.GetKeyToken()
	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(keyTokenAuth),
	})

	g.GET("/:id", r.courseCtr.GetCourse, isLoggedIn)

	// g.GET("/getuser", r.userCtr.GetLoginUser, isLoggedIn)

}
