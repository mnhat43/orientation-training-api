package models

import (
	cm "orientation-training-api/internal/common"
)

type UserProgress struct {
	cm.BaseModel

	UserID             int     `json:"user_id" pg:"user_id,notnull,on_delete:CASCADE"`
	CourseID           int     `json:"course_id" pg:"course_id,notnull,on_delete:CASCADE"`
	CoursePosition     int     `json:"course_position" pg:"course_position,notnull"`
	ModulePosition     int     `json:"module_position" pg:"module_position,notnull"`
	ModuleItemPosition int     `json:"module_item_position" pg:"module_item_position,notnull"`
	Completed          bool    `json:"completed" pg:"completed,default:false"`
	CompletedDate      string  `json:"completed_date" pg:"completed_date,default:null"`
	PerformanceRating  float64 `json:"performance_rating" pg:"performance_rating,default:null"`
	PerformanceComment string  `json:"performance_comment" pg:"performance_comment,default:null"`
	ReviewedBy         int     `json:"reviewed_by" pg:"reviewed_by,default:null"`

	// Define relationships
	User     *User   `json:"-" pg:"rel:has-one,fk:user_id"`
	Course   *Course `json:"-" pg:"rel:has-one,fk:course_id"`
	Reviewer User    `json:"reviewer" pg:"rel:belongs-to,fk:reviewed_by"`
}
