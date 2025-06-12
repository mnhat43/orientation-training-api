package appfeedback

import (
	cm "orientation-training-api/internal/common"
	"orientation-training-api/internal/interfaces/repository"
	"orientation-training-api/internal/interfaces/response"
	"orientation-training-api/internal/models"

	m "orientation-training-api/internal/models"

	"github.com/labstack/echo/v4"
)

type PgAppFeedbackRepository struct {
	cm.AppRepository
}

// NewPgAppFeedbackRepository creates a new instance of the PostgreSQL repository for app feedback
func NewPgAppFeedbackRepository(logger echo.Logger) repository.AppFeedbackRepository {
	repo := &PgAppFeedbackRepository{}
	repo.Init(logger)
	return repo
}

// CreateAppFeedback creates a new app feedback record
func (repo *PgAppFeedbackRepository) CreateAppFeedback(appFeedback *models.AppFeedback) (int, error) {
	err := repo.DB.Insert(appFeedback)
	if err != nil {
		repo.Logger.Errorf("Error creating app feedback: %v", err)
		return 0, err
	}

	return appFeedback.ID, nil
}

// GetAppFeedbackList gets a list of app feedbacks with user information
func (repo *PgAppFeedbackRepository) GetAppFeedbackList() ([]*response.FeedbackWithUser, error) {
	var feedbacks []*models.AppFeedback

	query := repo.DB.Model(&feedbacks).
		Where("deleted_at IS NULL").
		Order("created_at DESC")

	err := query.Select()

	if err != nil {
		repo.Logger.Errorf("Error getting app feedback list: %v", err)
		return nil, err
	}

	result := make([]*response.FeedbackWithUser, 0, len(feedbacks))

	for _, feedback := range feedbacks {
		var user models.User
		err := repo.DB.Model(&user).
			Column("usr.*").
			Where("usr.id = ?", feedback.UserID).
			Relation("UserProfile").
			Relation("Role").
			First()

		if err != nil {
			repo.Logger.Errorf("Error getting user for feedback: %v", err)
			continue
		}

		feedbackWithUser := response.CreateFeedbackWithUserFromAppFeedback(
			feedback,
			&user,
			&user.UserProfile,
			user.Role.Name,
		)

		result = append(result, feedbackWithUser)
	}

	return result, nil
}

// GetAppFeedbackByID gets an app feedback by ID
func (repo *PgAppFeedbackRepository) GetAppFeedbackByID(id int) (*models.AppFeedback, error) {
	feedback := new(models.AppFeedback)

	err := repo.DB.Model(feedback).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		First()

	if err != nil {
		repo.Logger.Errorf("Error getting app feedback by ID: %v", err)
		return nil, err
	}

	return feedback, nil
}

// GetAppFeedbackCount gets the total count of app feedbacks
func (repo *PgAppFeedbackRepository) GetAppFeedbackCount() (int, error) {
	count, err := repo.DB.Model((*models.AppFeedback)(nil)).
		Where("deleted_at IS NULL").
		Count()

	if err != nil {
		repo.Logger.Errorf("Error getting app feedback count: %v", err)
		return 0, err
	}

	return count, nil
}

// DeleteAppFeedback soft deletes an app feedback
func (repo *PgAppFeedbackRepository) DeleteAppFeedback(id int) error {
	feedback := &m.AppFeedback{BaseModel: cm.BaseModel{ID: id}}
	_, err := repo.DB.Model(feedback).WherePK().Delete()
	return err
}
