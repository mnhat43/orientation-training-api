package router

import (
	"orientation-training-api/internal/domains/auth"
	c "orientation-training-api/internal/domains/courses"
	lec "orientation-training-api/internal/domains/lectures"
	mdi "orientation-training-api/internal/domains/moduleitem"
	md "orientation-training-api/internal/domains/modules"
	uc "orientation-training-api/internal/domains/usercourse"
	u "orientation-training-api/internal/domains/users"
	cld "orientation-training-api/internal/platform/cloud"
	"orientation-training-api/internal/platform/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type AppRouter struct {
	authCtr       *auth.AuthController
	userCtr       *u.UserController
	courseCtr     *c.CourseController
	moduleCtr     *md.ModuleController
	moduleItemCtr *mdi.ModuleItemController
	ucCtr         *uc.UserCourseController
	lectureCtr    *lec.LectureController

	// adminCtr *ad.Controller

	userMw *u.UserMiddleware
	cld    *cld.CloudinaryStorage
}

func NewAppRouter(logger echo.Logger) (r *AppRouter) {
	userRepo := u.NewPgUserRepository(logger)
	courseRepo := c.NewPgCourseRepository(logger)
	moduleRepo := md.NewPgModuleRepository(logger)
	moduleItemRepo := mdi.NewPgModuleItemRepository(logger)
	ucRepo := uc.NewPgUserCourseRepository(logger)
	// adminRepo := ad.NewPgAdminRepository(logger)

	cldStorage := cld.NewCloudinaryStorage(logger)

	r = &AppRouter{
		authCtr:       auth.NewAuthController(logger, userRepo),
		userCtr:       u.NewUserController(logger, userRepo),
		courseCtr:     c.NewCourseController(logger, courseRepo, ucRepo, cldStorage),
		moduleCtr:     md.NewModuleController(logger, moduleRepo, moduleItemRepo, courseRepo),
		moduleItemCtr: mdi.NewModuleItemController(logger, moduleItemRepo, cldStorage),
		lectureCtr:    lec.NewLectureController(logger, moduleRepo, moduleItemRepo, courseRepo),

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
	// g.POST("/get-course-detail", r.courseCtr.GetCourseDetail, isLoggedIn, r.userMw.InitUserProfile, r.userMw.CheckManager)

	// g.GET("/getuser", r.userCtr.GetLoginUser, isLoggedIn)

}

func (r *AppRouter) ModuleRoute(g *echo.Group) {
	keyTokenAuth := utils.GetKeyToken()
	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(keyTokenAuth),
	})

	g.POST("/get-module-list", r.moduleCtr.GetModuleList, isLoggedIn, r.userMw.InitUserProfile)
	g.POST("/get-module-details", r.moduleCtr.GetModuleDetails, isLoggedIn, r.userMw.InitUserProfile)
	g.POST("/add-module", r.moduleCtr.AddModule, isLoggedIn, r.userMw.InitUserProfile, r.userMw.CheckManager)
	g.POST("/delete-module", r.moduleCtr.DeleteModule, isLoggedIn, r.userMw.InitUserProfile, r.userMw.CheckManager)
}

func (r *AppRouter) ModuleItemRoute(g *echo.Group) {
	keyTokenAuth := utils.GetKeyToken()
	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(keyTokenAuth),
	})

	g.POST("/get-module-item-list", r.moduleItemCtr.GetModuleItemList, isLoggedIn, r.userMw.InitUserProfile)
	g.POST("/add-module-item", r.moduleItemCtr.AddModuleItem, isLoggedIn, r.userMw.InitUserProfile, r.userMw.CheckManager)
	g.POST("/delete-module-item", r.moduleItemCtr.DeleteModuleItem, isLoggedIn, r.userMw.InitUserProfile, r.userMw.CheckManager)

	// g.POST("/add-module-item-video", r.moduleItemCtr.AddModuleItemVideo, isLoggedIn, r.userMw.InitUserProfile, r.userMw.CheckManager)
}

func (r *AppRouter) LectureRoute(g *echo.Group) {
	keyTokenAuth := utils.GetKeyToken()
	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(keyTokenAuth),
	})

	g.POST("/get-lecture-list", r.lectureCtr.GetLectureList, isLoggedIn, r.userMw.InitUserProfile)

}
