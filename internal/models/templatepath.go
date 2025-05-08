package models

import (
	cm "orientation-training-api/internal/common"
)

// TemplatePath represents a learning path template containing multiple courses
type TemplatePath struct {
	cm.BaseModel

	ID          int    `json:"id" pg:"id,pk"`
	Name        string `json:"name" pg:"name,notnull"`
	Description string `json:"description" pg:"description"`
	Courses     []int  `json:"courses" pg:"courses,array"`
	Duration    int    `json:"duration" pg:"duration"`
}
