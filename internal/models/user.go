package models

import (
	"time"

	cm "orientation-training-api/internal/common"
)

type User struct {
	cm.BaseModel

	tableName     struct{} `sql:"alias:usr"` //lint:ignore U1000 needed by ORM
	ID            int
	Email         string
	Password      string
	RoleID        int
	LastLoginTime time.Time

	UserProfile UserProfile
	Role        UserRole `pg:",fk:role_id"`
	// TargetEvaluation []TargetEvaluation `pg:",fk:user_id"`

}
