package models

import (
	"time"

	cm "orientation-training-api/internal/common"
)

// UserProfile : user profile info
type UserProfile struct {
	cm.BaseModel

	UserID            int
	Avatar            string
	FirstName         string
	LastName          string
	Birthday          time.Time
	PhoneNumber       string
	PersonalEmail     string
	CompanyJoinedDate time.Time
	Introduce         string
	Gender            int
}
