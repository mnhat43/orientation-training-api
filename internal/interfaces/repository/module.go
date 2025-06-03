package repository

import (
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
)

type ModuleRepository interface {
	GetModules(moduleListParams *param.ModuleListParams) ([]m.Module, int, error)
	SaveModule(createModuleParams *param.CreateModuleParams) (m.Module, error)
	GetModuleByID(id int) (m.Module, error)
	DeleteModule(moduleID int) error
	GetModuleIDsByCourseID(courseID int) ([]int, error)
	GetModulesByCourseID(courseID int) ([]m.Module, error)
	GetModuleByPositionAndCourse(courseID int, position int) (m.Module, error)
}
