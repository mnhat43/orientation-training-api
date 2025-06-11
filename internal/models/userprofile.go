package models

import (
	cm "orientation-training-api/internal/common"
)

// UserProfile : user profile info
type UserProfile struct {
	cm.BaseModel

	UserID            int
	Avatar            string
	FirstName         string
	LastName          string
	Birthday          string
	PhoneNumber       string
	PersonalEmail     string
	Department        string
	CompanyJoinedDate string
	Introduce         string
	Gender            int
}
