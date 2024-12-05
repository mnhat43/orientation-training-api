package router

import (
	"orientation-training-api/internal/domains/auth"
	c "orientation-training-api/internal/domains/courses"
	uc "orientation-training-api/internal/domains/usercourse"
	u "orientation-training-api/internal/domains/users"
	cld "orientation-training-api/internal/platform/cloud"
	"orientation-training-api/internal/platform/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type AppRouter struct {
	authCtr   *auth.AuthController
	userCtr   *u.UserController
	courseCtr *c.CourseController
	ucCtr     *uc.UserCourseController
	// adminCtr *ad.Controller

	userMw *u.UserMiddleware
	cld    *cld.CloudinaryStorage
}

func NewAppRouter(logger echo.Logger) (r *AppRouter) {
	userRepo := u.NewPgUserRepository(logger)
	courseRepo := c.NewPgCourseRepository(logger)
	ucRepo := uc.NewPgUserCourseRepository(logger)
	// adminRepo := ad.NewPgAdminRepository(logger)

	cldStorage := cld.NewCloudinaryStorage(logger)

	r = &AppRouter{
		authCtr:   auth.NewAuthController(logger, userRepo),
		userCtr:   u.NewUserController(logger, userRepo),
		courseCtr: c.NewCourseController(logger, courseRepo, ucRepo, cldStorage),
		// adminCtr: ad.NewAdminController(),

		userMw: u.NewUserMiddleware(logger, userRepo),
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

	g.POST("/get-course-list", r.courseCtr.GetCourseList, isLoggedIn, r.userMw.InitUserProfile)
	g.POST("/add-course", r.courseCtr.AddCourse, isLoggedIn, r.userMw.InitUserProfile, r.userMw.CheckManager)
	// g.POST("/update-course", r.courseCtr.UpdateCourse, isLoggedIn, r.userMw.InitUserProfile, r.userMw.CheckManager)
	g.POST("/delete-course", r.courseCtr.DeleteCourse, isLoggedIn, r.userMw.InitUserProfile, r.userMw.CheckManager)

	// g.GET("/getuser", r.userCtr.GetLoginUser, isLoggedIn)

}
