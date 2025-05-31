package users

import (
	cm "orientation-training-api/internal/common"
	"orientation-training-api/internal/models"
	"orientation-training-api/internal/platform/utils"

	m "orientation-training-api/internal/models"

	"github.com/labstack/echo/v4"
)

type PgUserRepository struct {
	cm.AppRepository
}

func NewPgUserRepository(logger echo.Logger) (repo *PgUserRepository) {
	repo = &PgUserRepository{}
	repo.Init(logger)
	return
}

func (repo *PgUserRepository) GetLoginUserID(email string, password string) (int, error) {
	user := m.User{}
	err := repo.DB.Model(&user).
		Column("id").
		Where("email = ?", email).
		Where("password = ?", password).
		Where("deleted_at is null").
		Select()

	if err != nil {
		repo.Logger.Errorf("%+v", err)
	}

	return user.ID, err
}

func (repo *PgUserRepository) UpdateLastLogin(userID int) error {
	_, err := repo.DB.Model(&m.User{LastLoginTime: utils.TimeNowUTC()}).
		Column("last_login_time", "updated_at").
		Where("id = ?", userID).
		Update()

	return err
}

func (repo *PgUserRepository) GetUserProfile(id int) (m.User, error) {
	user := m.User{}
	err := repo.DB.Model(&user).
		Column("usr.*").
		Where("usr.id = ?", id).
		Where("usr.deleted_at is null").
		Relation("UserProfile").
		Relation("Role").
		First()

	if err != nil {
		repo.Logger.Errorf("%+v", err)
	}

	return user, err
}

// GetUsersByRoleID retrieves all users with the specified role ID
func (repo *PgUserRepository) GetUsersByRoleID(roleID int) ([]m.User, error) {
	var users []m.User
	err := repo.DB.Model(&users).
		Column("usr.*").
		Where("usr.role_id = ?", roleID).
		Where("usr.deleted_at is null").
		Relation("UserProfile").
		Relation("Role").
		Select()

	if err != nil {
		repo.Logger.Errorf("Error getting users by role ID: %+v", err)
	}

	return users, err
}

// GetUserProgressByUserID retrieves all user progress entries for a specific user
func (repo *PgUserRepository) GetUserProgressByUserID(userID int) ([]models.UserProgress, error) {
	var userProgresses []models.UserProgress
	err := repo.DB.Model(&userProgresses).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Select()
	return userProgresses, err
}
