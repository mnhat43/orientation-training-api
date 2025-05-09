package templatepaths

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

type TemplatePathController struct {
	cm.BaseController
	TempPathRepo rp.TemplatePathRepository
	CourseRepo   rp.CourseRepository
}

func NewTemplatePathController(logger echo.Logger, tempPathRepo rp.TemplatePathRepository, courseRepo rp.CourseRepository) (ctr *TemplatePathController) {
	ctr = &TemplatePathController{cm.BaseController{}, tempPathRepo, courseRepo}
	ctr.Init(logger)
	return
}

// GetTemplatePathList retrieves all template paths
func (ctr *TemplatePathController) GetTemplatePathList(c echo.Context) error {
	tempPaths, err := ctr.TempPathRepo.GetTemplatePathList()
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch template paths: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch template paths",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Template paths retrieved successfully",
		Data:    tempPaths,
	})
}

// GetTemplatePath retrieves a template path by ID
func (ctr *TemplatePathController) GetTemplatePath(c echo.Context) error {
	tempPathIDParam := new(param.TempPathIDParam)
	if err := c.Bind(tempPathIDParam); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
		})
	}

	if _, err := valid.ValidateStruct(tempPathIDParam); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	templatePath, err := ctr.TempPathRepo.GetTemplatePathByID(tempPathIDParam.TempPathID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch template path: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Template path not found",
		})
	}

	courseDetails := []map[string]interface{}{}
	for _, courseID := range templatePath.CourseIds {
		course, err := ctr.CourseRepo.GetCourseByID(courseID)
		if err != nil {
			ctr.Logger.Warnf("Course ID %d in path %d not found: %v", courseID, templatePath.ID, err)
			continue
		}

		courseDetails = append(courseDetails, map[string]interface{}{
			"id":          course.ID,
			"title":       course.Title,
			"thumbnail":   course.Thumbnail,
			"category":    course.Category,
			"duration":    course.Duration,
			"description": course.Description,
		})
	}

	response := map[string]interface{}{
		"id":          templatePath.ID,
		"name":        templatePath.Name,
		"description": templatePath.Description,
		"course_ids":  templatePath.CourseIds,
		"course_list": courseDetails,
		"duration":    templatePath.Duration,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Template path retrieved successfully",
		Data:    response,
	})
}

// CreateTemplatePath creates a new template path
func (ctr *TemplatePathController) CreateTemplatePath(c echo.Context) error {
	createTemplatePathParams := new(param.CreateTemplatePathParams)
	if err := c.Bind(createTemplatePathParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
		})
	}

	if _, err := valid.ValidateStruct(createTemplatePathParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	totalDuration := 0
	for _, courseID := range createTemplatePathParams.CourseIds {
		course, err := ctr.CourseRepo.GetCourseByID(courseID)
		if err != nil {
			ctr.Logger.Errorf("Course ID %d not found: %v", courseID, err)
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "One or more courses do not exist",
			})
		}
		totalDuration += course.Duration
	}

	createTemplatePathParams.Duration = totalDuration

	templatePath, err := ctr.TempPathRepo.CreateTemplatePath(createTemplatePathParams)
	if err != nil {
		ctr.Logger.Errorf("Failed to create template path: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to create template path",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Template path created successfully",
		Data:    templatePath,
	})
}

// UpdateTemplatePath updates an existing template path
func (ctr *TemplatePathController) UpdateTemplatePath(c echo.Context) error {
	updatePathParams := new(param.UpdateTemplatePathParams)
	if err := c.Bind(updatePathParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
		})
	}

	if _, err := valid.ValidateStruct(updatePathParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	existingPath, err := ctr.TempPathRepo.GetTemplatePathByID(updatePathParams.TempPathID)
	if err != nil {
		ctr.Logger.Errorf("Template path not found: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Template path not found",
		})
	}

	if len(updatePathParams.CourseIds) > 0 {
		totalDuration := 0
		for _, courseID := range updatePathParams.CourseIds {
			course, err := ctr.CourseRepo.GetCourseByID(courseID)
			if err != nil {
				ctr.Logger.Errorf("Course ID %d not found: %v", courseID, err)
				return c.JSON(http.StatusOK, cf.JsonResponse{
					Status:  cf.FailResponseCode,
					Message: "One or more courses do not exist",
				})
			}
			totalDuration += course.Duration
		}

		updatePathParams.Duration = totalDuration
	}

	if updatePathParams.Name == "" {
		updatePathParams.Name = existingPath.Name
	}
	if updatePathParams.Description == "" {
		updatePathParams.Description = existingPath.Description
	}

	tempPath, err := ctr.TempPathRepo.UpdateTemplatePath(updatePathParams)
	if err != nil {
		ctr.Logger.Errorf("Failed to update template path: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to update template path",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Template path updated successfully",
		Data:    tempPath,
	})
}

// DeleteTemplatePath deletes a template path
func (ctr *TemplatePathController) DeleteTemplatePath(c echo.Context) error {
	tempPathIDParam := new(param.TempPathIDParam)
	if err := c.Bind(tempPathIDParam); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
		})
	}

	if _, err := valid.ValidateStruct(tempPathIDParam); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	_, err := ctr.TempPathRepo.GetTemplatePathByID(tempPathIDParam.TempPathID)
	if err != nil {
		ctr.Logger.Errorf("Template path not found: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Template path not found",
		})
	}

	err = ctr.TempPathRepo.DeleteTemplatePath(tempPathIDParam.TempPathID)
	if err != nil {
		ctr.Logger.Errorf("Failed to delete template path: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to delete template path",
		})
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Template path deleted successfully",
	})
}
