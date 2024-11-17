package repository

import (
	m "orientation-training-api/internal/models"
)

type CourseRepository interface {
	GetCourseByID(id int) (m.Course, error)
}
