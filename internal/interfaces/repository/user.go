package repository

import (
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
)

type UserRepository interface {
	GetLoginUserID(email string, password string) (int, error)
	UpdateLastLogin(userID int) error
	GetUserProfile(id int) (m.User, error)
	GetUsersByRoleID(roleID int) ([]m.User, error)
	GetUserProgressByUserID(userID int) ([]m.UserProgress, error)
	GetUsersWithoutProgress(roleID int) ([]m.User, error)
	CreateUser(user m.User) (int, error)
	CheckEmailExists(email string) (bool, error)
	UpdateUserProfile(userID int, profileParams *param.UpdateProfileParams) error
}
