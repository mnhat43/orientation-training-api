package repository

import (
	m "orientation-training-api/internal/models"
)

type UserRepository interface {
	GetLoginUserID(email string, password string) (int, error)
	UpdateLastLogin(userID int) error
	GetUserProfile(id int) (m.User, error)
	GetUsersByRoleID(roleID int) ([]m.User, error)
}
