package repository

import (
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
)

type CourseRepository interface {
	GetCourseByID(id int) (m.Course, error)
	GetCourses(courseListParams *param.CourseListParams) ([]m.Course, int, error)
	SaveCourse(createCourseDBParams *param.CreateCourseDBParams, userCourseRepo UserCourseRepository) (m.Course, error)
	UpdateCourse(courseParams *param.UpdateCourseParams, userCourseRepo UserCourseRepository) error
	DeleteCourse(courseID int) error
}
