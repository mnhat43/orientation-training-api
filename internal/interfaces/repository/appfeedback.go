package repository

import (
	"orientation-training-api/internal/interfaces/response"
	"orientation-training-api/internal/models"
)

// AppFeedbackRepository interface for app feedback operations
type AppFeedbackRepository interface {
	CreateAppFeedback(appFeedback *models.AppFeedback) (int, error)
	GetAppFeedbackList() ([]*response.FeedbackWithUser, error)
	GetAppFeedbackByID(id int) (*models.AppFeedback, error)
	GetAppFeedbackCount() (int, error)
	DeleteAppFeedback(id int) error
}
