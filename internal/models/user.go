package models

import (
	"time"

	cm "orientation-training-api/internal/common"
)

type User struct {
	cm.BaseModel

	tableName     struct{} `sql:"alias:usr"` //lint:ignore U1000 needed by ORM
	Email         string
	Password      string
	RoleID        int
	LastLoginTime time.Time

	UserProfile UserProfile `pg:"rel:has-one"`
	Role        UserRole    `pg:"rel:belongs-to,fk:role_id"`
	// TargetEvaluation []TargetEvaluation `pg:",fk:user_id"`

}
