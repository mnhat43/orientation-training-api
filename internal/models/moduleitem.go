package models

import (
	cm "orientation-training-api/internal/common"
)

type ModuleItem struct {
	cm.BaseModel

	Title        string `pg:"title" json:"title"`
	ItemType     string `pg:"item_type" json:"item_type"`
	Resource     string `pg:"resource,null" json:"resource,omitempty"`
	Position     int    `pg:"position" json:"position"`
	RequiredTime int    `pg:"required_time,null" json:"required_time,omitempty"`
	ModuleID     int    `pg:"module_id" json:"module_id"`
	QuizID       int    `pg:"quiz_id,null" json:"quiz_id,omitempty"`

	Module Module `pg:"rel:has-one"`
	Quiz   *Quiz  `pg:"rel:has-one,fk:quiz_id" json:"quiz,omitempty"`
}
