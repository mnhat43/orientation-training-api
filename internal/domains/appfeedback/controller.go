package appfeedback

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"
	"time"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

// AppFeedbackController handles HTTP requests related to app feedback
type AppFeedbackController struct {
	cm.BaseController
	AppFeedbackRepo rp.AppFeedbackRepository
}

// NewAppFeedbackController creates a new instance of AppFeedbackController
func NewAppFeedbackController(logger echo.Logger, appFeedbackRepo rp.AppFeedbackRepository) (ctr *AppFeedbackController) {
	ctr = &AppFeedbackController{cm.BaseController{}, appFeedbackRepo}
	ctr.Init(logger)
	return
}

// SubmitAppFeedback handles app feedback submission
func (ctr *AppFeedbackController) SubmitAppFeedback(c echo.Context) error {
	userProfile := c.Get("user_profile").(m.User)

	req := new(param.SubmitAppFeedbackRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid request format",
		})
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}
	now := time.Now()
	submitTime := now

	if req.SubmittedAt != "" {
		if parsedTime, err := time.Parse(time.RFC3339, req.SubmittedAt); err == nil {
			submitTime = parsedTime
		} else {
			ctr.Logger.Warn("Invalid submittedAt format:", err)
		}
	}
	appFeedback := &m.AppFeedback{
		Rating:   req.Rating,
		Feedback: req.Feedback,
		SubmitAt: submitTime,
		UserID:   userProfile.ID,
	}

	id, err := ctr.AppFeedbackRepo.CreateAppFeedback(appFeedback)
	if err != nil {
		ctr.Logger.Errorf("Error submitting app feedback: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to submit feedback",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Feedback submitted successfully",
		Data: map[string]interface{}{
			"id": id,
		},
	})
}

// GetAppFeedbackList handles retrieval of app feedback list
func (ctr *AppFeedbackController) GetAppFeedbackList(c echo.Context) error {
	feedbacks, err := ctr.AppFeedbackRepo.GetAppFeedbackList()
	if err != nil {
		ctr.Logger.Errorf("Error getting app feedback list: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to get feedback list",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Feedback list retrieved successfully",
		Data:    feedbacks,
	})
}

// DeleteAppFeedback handles deleting app feedback
func (ctr *AppFeedbackController) DeleteAppFeedback(c echo.Context) error {
	req := new(param.DeleteAppFeedbackRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid request format",
		})
	}

	if _, err := valid.ValidateStruct(req); err != nil {
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	err := ctr.AppFeedbackRepo.DeleteAppFeedback(req.ID)
	if err != nil {
		ctr.Logger.Errorf("Error deleting app feedback: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to delete feedback",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Feedback deleted successfully",
	})
}
