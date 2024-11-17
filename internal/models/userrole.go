package models

import (
	cm "orientation-training-api/internal/common"
)

// UserRole : struct for db table user_roles
type UserRole struct {
	cm.BaseModel

	Name        string
	Description string
}
