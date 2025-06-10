package models

import (
	cm "orientation-training-api/internal/common"
)

// SkillKeyword : struct for db table skill_keywords
type SkillKeyword struct {
	cm.BaseModel

	Name string `json:"name" pg:"name,unique,notnull"`
}

// CourseSkillKeyword : struct for db table course_skill_keywords
type CourseSkillKeyword struct {
	cm.BaseModel

	CourseID       int `json:"course_id" pg:"course_id,notnull"`
	SkillKeywordID int `json:"skill_keyword_id" pg:"skill_keyword_id,notnull"`

	Course       *Course       `pg:"rel:has-one,fk:course_id"`
	SkillKeyword *SkillKeyword `pg:"rel:has-one,fk:skill_keyword_id"`
}
