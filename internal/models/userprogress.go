package models

import (
	cm "orientation-training-api/internal/common"
)

type UserProgress struct {
	cm.BaseModel

	ID                 int  `json:"id" pg:"id,pk"`
	UserID             int  `json:"user_id" pg:"user_id,notnull"`
	CourseID           int  `json:"course_id" pg:"course_id,notnull"`
	ModulePosition     int  `json:"module_position" pg:"module_position,notnull"`
	ModuleItemPosition int  `json:"module_item_position" pg:"module_item_position,notnull"`
	Completed          bool `json:"completed" pg:"completed,default:false"`
}
