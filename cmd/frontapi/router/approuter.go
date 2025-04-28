package router

import (
	"orientation-training-api/internal/domains/auth"
	c "orientation-training-api/internal/domains/courses"
	lec "orientation-training-api/internal/domains/lectures"
	mdi "orientation-training-api/internal/domains/moduleitem"
	md "orientation-training-api/internal/domains/modules"
	uc "orientation-training-api/internal/domains/usercourse"
	up "orientation-training-api/internal/domains/userprogress"
	u "orientation-training-api/internal/domains/users"

	gc "orientation-training-api/internal/platform/cloud"
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
	upCtr         *up.UserProgressController

	// adminCtr *ad.Controller

	userMw *u.UserMiddleware
	gcs    *gc.GcsStorage
}

func NewAppRouter(logger echo.Logger) (r *AppRouter) {
	userRepo := u.NewPgUserRepository(logger)
	courseRepo := c.NewPgCourseRepository(logger)
	moduleRepo := md.NewPgModuleRepository(logger)
	moduleItemRepo := mdi.NewPgModuleItemRepository(logger)
	ucRepo := uc.NewPgUserCourseRepository(logger)
	upRepo := up.NewPgUserProgressRepository(logger)
	// adminRepo := ad.NewPgAdminRepository(logger)

	gcsStorage := gc.NewGcsStorage(logger)

	r = &AppRouter{
		authCtr:       auth.NewAuthController(logger, userRepo),
		userCtr:       u.NewUserController(logger, userRepo),
		courseCtr:     c.NewCourseController(logger, courseRepo, ucRepo, gcsStorage),
		moduleCtr:     md.NewModuleController(logger, moduleRepo, moduleItemRepo, courseRepo),
		moduleItemCtr: mdi.NewModuleItemController(logger, moduleItemRepo, gcsStorage),
		lectureCtr:    lec.NewLectureController(logger, moduleRepo, moduleItemRepo, courseRepo, upRepo, gcsStorage),
		upCtr:         up.NewUserProgressController(logger, upRepo, moduleRepo, moduleItemRepo, userRepo),

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

	g.GET("/profile", r.userCtr.GetLoginUser, isLoggedIn)

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

func (r *AppRouter) UserProgressRoute(g *echo.Group) {
	keyTokenAuth := utils.GetKeyToken()
	isLoggedIn := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(keyTokenAuth),
	})

	g.POST("/update-user-progress", r.upCtr.UpdateUserProgress, isLoggedIn, r.userMw.InitUserProfile)
	g.POST("/get-user-progress", r.upCtr.GetUserProgress, isLoggedIn, r.userMw.InitUserProfile)
	g.POST("/add-user-progress", r.upCtr.AddUserProgress, isLoggedIn, r.userMw.InitUserProfile, r.userMw.CheckManager)
	g.POST("/list-trainee-by-course", r.upCtr.GetListTraineeByCourseID, isLoggedIn, r.userMw.InitUserProfile, r.userMw.CheckManager)
}
