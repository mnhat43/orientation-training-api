package models

import (
	cm "orientation-training-api/internal/common"
)

type Course struct {
	cm.BaseModel

	ID          int
	Title       string
	Description string
	CreatedBy   int

	User User `pg:",fk:created_by"`
}
