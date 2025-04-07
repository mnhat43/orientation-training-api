package models

import (
	cm "orientation-training-api/internal/common"
)

type Module struct {
	cm.BaseModel

	ID       int    `pg:",pk"`
	Title    string `pg:",notnull"`
	CourseID int    `pg:",fk:course_id"`
	Position int    `pg:",notnull"`
}
