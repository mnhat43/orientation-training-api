package usercourse

import (
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
)

type UserCourseController struct {
	cm.BaseController

	UserCourseRepo rp.UserCourseRepository
	UserRepo       rp.UserRepository
	CourseRepo     rp.CourseRepository
}
