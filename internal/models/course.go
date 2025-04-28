package models

import (
	cm "orientation-training-api/internal/common"
)

type Course struct {
	cm.BaseModel

	ID          int    `pg:",pk"`
	Title       string `pg:",notnull"`
	Description string
	Thumbnail   string
	Category    string `pg:",notnull"`
	CreatedBy   int    `pg:",fk:created_by"`

	// User User `pg:"rel:has-one"`
}

type CourseDetail struct {
	cm.BaseModel

	ID          int    `pg:",pk"`
	Title       string `pg:",notnull"`
	Description string
	Thumbnail   string
	Category    string `pg:",notnull"`
	CreatedBy   int    `pg:",fk:created_by"`

	User User `pg:"rel:has-one"`
}
