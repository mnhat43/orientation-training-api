package models

import (
	cm "orientation-training-api/internal/common"
)

// TemplatePath represents a learning path template containing multiple courses
type TemplatePath struct {
	cm.BaseModel

	ID          int    `pg:"id,pk" json:"id"`
	Name        string `pg:"name,notnull" json:"name"`
	Description string `pg:"description" json:"description"`
	CourseIds   []int  `pg:"course_ids,array" json:"course_ids"`
	Duration    int    `pg:"duration" json:"duration"`
}
