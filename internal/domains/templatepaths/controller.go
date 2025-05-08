package templatepaths

import (
	"net/http"
	cf "orientation-training-api/configs"
	cm "orientation-training-api/internal/common"
	rp "orientation-training-api/internal/interfaces/repository"
	param "orientation-training-api/internal/interfaces/requestparams"
	m "orientation-training-api/internal/models"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

type TemplatePathController struct {
	cm.BaseController
	PathRepo   rp.TemplatePathRepository
	CourseRepo rp.CourseRepository
}

func NewTemplatePathController(logger echo.Logger, pathRepo rp.TemplatePathRepository, courseRepo rp.CourseRepository) (ctr *TemplatePathController) {
	ctr = &TemplatePathController{cm.BaseController{}, pathRepo, courseRepo}
	ctr.Init(logger)
	return
}

// GetTemplatePathList retrieves all template paths
func (ctr *TemplatePathController) GetTemplatePathList(c echo.Context) error {
	paths, err := ctr.PathRepo.GetTemplatePathList()
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch template paths: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch template paths",
		})
	}

	if paths == nil {
		paths = []m.TemplatePath{}
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Template paths retrieved successfully",
		Data:    paths,
	})
}

// GetTemplatePath retrieves a template path by ID
func (ctr *TemplatePathController) GetTemplatePath(c echo.Context) error {
	getPathParams := new(param.GetTemplatePathParams)
	if err := c.Bind(getPathParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
		})
	}

	if _, err := valid.ValidateStruct(getPathParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	path, err := ctr.PathRepo.GetTemplatePathByID(getPathParams.TempPathID)
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch template path: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Template path not found",
		})
	}

	courseDetails := []map[string]interface{}{}
	for _, courseID := range path.Courses {
		course, err := ctr.CourseRepo.GetCourseByID(courseID)
		if err != nil {
			ctr.Logger.Warnf("Course ID %d in path %d not found: %v", courseID, path.ID, err)
			continue
		}

		courseDetails = append(courseDetails, map[string]interface{}{
			"id":          course.ID,
			"title":       course.Title,
			"description": course.Description,
			"thumbnail":   course.Thumbnail,
			"duration":    course.Duration,
		})
	}

	response := map[string]interface{}{
		"id":          path.ID,
		"name":        path.Name,
		"description": path.Description,
		"course_ids":  path.Courses,
		"course_list": courseDetails,
		"duration":    path.Duration,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Template path retrieved successfully",
		Data:    response,
	})
}

// CreateTemplatePath creates a new template path
func (ctr *TemplatePathController) CreateTemplatePath(c echo.Context) error {
	createPathParams := new(param.CreateTemplatePathParams)
	if err := c.Bind(createPathParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
		})
	}

	if _, err := valid.ValidateStruct(createPathParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	totalDuration := 0
	for _, courseID := range createPathParams.Courses {
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

	path := &m.TemplatePath{
		Name:        createPathParams.Name,
		Description: createPathParams.Description,
		Courses:     createPathParams.Courses,
		Duration:    totalDuration,
	}

	err := ctr.PathRepo.CreateTemplatePath(path)
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
		Data:    path,
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

	existingPath, err := ctr.PathRepo.GetTemplatePathByID(updatePathParams.TempPathID)
	if err != nil {
		ctr.Logger.Errorf("Template path not found: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Template path not found",
		})
	}

	if len(updatePathParams.Courses) > 0 {
		totalDuration := 0
		for _, courseID := range updatePathParams.Courses {
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

		existingPath.Courses = updatePathParams.Courses
		existingPath.Duration = totalDuration
	}

	if updatePathParams.Name != "" {
		existingPath.Name = updatePathParams.Name
	}
	if updatePathParams.Description != "" {
		existingPath.Description = updatePathParams.Description
	}

	err = ctr.PathRepo.UpdateTemplatePath(&existingPath)
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
		Data:    existingPath,
	})
}

// DeleteTemplatePath deletes a template path
func (ctr *TemplatePathController) DeleteTemplatePath(c echo.Context) error {
	deletePathParams := new(param.DeleteTemplatePathParams)
	if err := c.Bind(deletePathParams); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
		})
	}

	if _, err := valid.ValidateStruct(deletePathParams); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	_, err := ctr.PathRepo.GetTemplatePathByID(deletePathParams.TempPathID)
	if err != nil {
		ctr.Logger.Errorf("Template path not found: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Template path not found",
		})
	}

	err = ctr.PathRepo.DeleteTemplatePath(deletePathParams.TempPathID)
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
