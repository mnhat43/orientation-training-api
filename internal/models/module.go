package models

import (
	cm "orientation-training-api/internal/common"
)

// Module represents a module in a course
type Module struct {
	cm.BaseModel

	ID       int    `json:"id" pg:"id,pk"`
	CourseID int    `json:"course_id" pg:"course_id,notnull"`
	Title    string `json:"title" pg:"title,notnull"`
	Duration int    `json:"duration" pg:"duration,default:0"`
	Position int    `json:"position" pg:"position,notnull"`
}
