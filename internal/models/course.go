package models

import (
	cm "orientation-training-api/internal/common"
)

type Course struct {
	cm.BaseModel

	Title       string `pg:",notnull"`
	Description string
	Thumbnail   string
	Category    string `pg:",notnull"`
	Duration    int    `pg:",default:0"`
	CreatedBy   int    `pg:",fk:created_by"`

	// User User `pg:"rel:has-one"`
}

type CourseDetail struct {
	cm.BaseModel

	Title       string `pg:",notnull"`
	Description string
	Thumbnail   string
	Category    string `pg:",notnull"`
	Duration    int    `pg:",default:0"`
	CreatedBy   int    `pg:",fk:created_by"`

	User User `pg:"rel:has-one"`
}
