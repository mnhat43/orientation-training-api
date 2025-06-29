package users

import (
	cm "orientation-training-api/internal/common"
	param "orientation-training-api/internal/interfaces/requestparams"
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

// GetUsersWithoutProgress retrieves all users with specified role ID who don't have any records in user_progresses table
func (repo *PgUserRepository) GetUsersWithoutProgress(roleID int) ([]m.User, error) {
	var users []m.User
	err := repo.DB.Model(&users).
		Column("usr.*").
		Where("usr.role_id = ?", roleID).
		Where("usr.deleted_at is null").
		Where("NOT EXISTS (SELECT 1 FROM user_progresses up WHERE up.user_id = usr.id AND up.deleted_at IS NULL)").
		Relation("UserProfile").
		Relation("Role").
		Select()

	if err != nil {
		repo.Logger.Errorf("Error getting users without progress: %+v", err)
	}

	return users, err
}

// UpdateUserProfile updates a user's profile information and returns the updated user
func (repo *PgUserRepository) UpdateUserProfile(userID int, profileParams *param.UpdateProfileParams) error {
	_, err := repo.GetUserProfile(userID)
	if err != nil {
		repo.Logger.Errorf("Error getting user profile: %+v", err)
		return err
	}

	updateQuery := repo.DB.Model(&m.UserProfile{}).
		Set("first_name = ?", profileParams.FirstName).
		Set("last_name = ?", profileParams.LastName).
		Set("phone_number = ?", profileParams.PhoneNumber).
		Set("birthday = ?", profileParams.Birthday).
		Set("updated_at = ?", utils.TimeNowUTC())

	if profileParams.Avatar != "" {
		updateQuery = updateQuery.Set("avatar = ?", profileParams.Avatar)
	}

	_, err = updateQuery.Where("user_id = ?", userID).Update()

	if err != nil {
		repo.Logger.Errorf("Error updating user profile: %+v", err)
		return err
	}

	return nil
}

// UpdatePassword updates a user's password
func (repo *PgUserRepository) UpdatePassword(userID int, newHashedPassword string) error {
	_, err := repo.DB.Model(&m.User{}).
		Set("password = ?", newHashedPassword).
		Set("updated_at = ?", utils.TimeNowUTC()).
		Where("id = ?", userID).
		Update()

	if err != nil {
		repo.Logger.Errorf("Error updating user password: %+v", err)
		return err
	}

	return nil
}

// GetAllUsersExceptRole retrieves all users except those with the specified role ID
func (repo *PgUserRepository) GetAllUsersExceptRole(roleID int) ([]m.User, error) {
	var users []m.User
	err := repo.DB.Model(&users).
		Column("usr.*").
		Where("usr.role_id != ?", roleID).
		Where("usr.deleted_at is null").
		Relation("UserProfile").
		Relation("Role").
		Select()

	if err != nil {
		repo.Logger.Errorf("Error getting users except role ID %d: %+v", roleID, err)
	}

	return users, err
}

// AdminUpdateUser updates a user's profile and role information by admin
func (repo *PgUserRepository) AdminUpdateUser(userID int, profileParams *param.AdminUpdateUserParams) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Errorf("Error starting transaction: %+v", err)
		return err
	}

	// Update user profile
	_, err = tx.Model(&m.UserProfile{
		FirstName:         profileParams.FirstName,
		LastName:          profileParams.LastName,
		PhoneNumber:       profileParams.PhoneNumber,
		Birthday:          profileParams.Birthday,
		Department:        profileParams.Department,
		Avatar:            profileParams.Avatar,
		Gender:            profileParams.Gender,
		CompanyJoinedDate: profileParams.CompanyJoinedDate,
	}).
		Column("first_name", "last_name", "phone_number", "birthday",
			"department", "avatar", "gender", "company_joined_date", "updated_at").
		Where("user_id = ?", userID).
		Where("deleted_at is null").
		Update()

	if err != nil {
		tx.Rollback()
		repo.Logger.Errorf("Error updating user profile: %+v", err)
		return err
	}

	// If role ID is provided, update user role
	if profileParams.RoleID > 0 {
		_, err = tx.Model(&m.User{
			RoleID: profileParams.RoleID,
		}).
			Column("role_id", "updated_at").
			Where("id = ?", userID).
			Where("deleted_at is null").
			Update()

		if err != nil {
			tx.Rollback()
			repo.Logger.Errorf("Error updating user role: %+v", err)
			return err
		}
	}

	return tx.Commit()
}

// DeleteUser soft deletes a user by setting deleted_at
func (repo *PgUserRepository) DeleteUser(userID int) error {
	now := utils.TimeNowUTC()

	// Begin transaction
	tx, err := repo.DB.Begin()
	if err != nil {
		repo.Logger.Errorf("Error starting transaction: %+v", err)
		return err
	}

	// Soft delete user profile
	_, err = tx.Model(&m.UserProfile{}).
		Set("deleted_at = ?", now).
		Where("user_id = ?", userID).
		Where("deleted_at is null").
		Update()

	if err != nil {
		tx.Rollback()
		repo.Logger.Errorf("Error soft deleting user profile: %+v", err)
		return err
	}

	// Soft delete user
	_, err = tx.Model(&m.User{}).
		Set("deleted_at = ?", now).
		Where("id = ?", userID).
		Where("deleted_at is null").
		Update()

	if err != nil {
		tx.Rollback()
		repo.Logger.Errorf("Error soft deleting user: %+v", err)
		return err
	}

	return tx.Commit()
}
