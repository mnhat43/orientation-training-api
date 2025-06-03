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

// CheckEmailExists checks if an email already exists in the database
func (repo *PgUserRepository) CheckEmailExists(email string) (bool, error) {
	count, err := repo.DB.Model(&m.User{}).
		Where("email = ?", email).
		Where("deleted_at is null").
		Count()

	if err != nil {
		repo.Logger.Errorf("Error checking email existence: %+v", err)
		return false, err
	}

	return count > 0, nil
}

// CreateUser creates a new user with profile information
func (repo *PgUserRepository) CreateUser(user m.User) (int, error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Errorf("Error starting transaction: %+v", err)
		return 0, err
	}

	now := utils.TimeNowUTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := tx.Insert(&user); err != nil {
		tx.Rollback()
		repo.Logger.Errorf("Error inserting user: %+v", err)
		return 0, err
	}

	user.UserProfile.UserID = user.ID
	user.UserProfile.CreatedAt = now
	user.UserProfile.UpdatedAt = now

	if err := tx.Insert(&user.UserProfile); err != nil {
		tx.Rollback()
		repo.Logger.Errorf("Error inserting user profile: %+v", err)
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		repo.Logger.Errorf("Error committing transaction: %+v", err)
		return 0, err
	}

	return user.ID, nil
}
