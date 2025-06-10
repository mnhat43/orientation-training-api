package courseskillkeyword

import (
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"

	"github.com/labstack/echo/v4"
)

type CourseSkillKeywordController struct {
	cm.BaseController

	CourseSkillKeywordRepo rp.CourseSkillKeywordRepository
}

func NewCourseSkillKeywordController(logger echo.Logger, courseSkillKeywordRepo rp.CourseSkillKeywordRepository) (ctr *CourseSkillKeywordController) {
	ctr = &CourseSkillKeywordController{
		cm.BaseController{},
		courseSkillKeywordRepo,
	}
	ctr.Init(logger)
	return
}
