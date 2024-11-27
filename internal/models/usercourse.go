package models

import (
	cm "orientation-training-api/internal/common"
)

// UserCourse : struct for db table user_courses
type UserCourse struct {
	cm.BaseModel

	UserID   int
	CourseID int
}
