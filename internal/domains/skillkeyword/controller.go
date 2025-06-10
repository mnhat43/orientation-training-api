package skillkeyword

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

type SkillKeywordController struct {
	cm.BaseController
	sKeyRepo rp.SkillKeywordRepository
}

func NewSkillKeywordController(logger echo.Logger, sKeyRepo rp.SkillKeywordRepository) (ctr *SkillKeywordController) {
	ctr = &SkillKeywordController{cm.BaseController{}, sKeyRepo}
	ctr.Init(logger)
	return
}

func (ctr *SkillKeywordController) GetSkillKeywordList(c echo.Context) error {
	list, err := ctr.sKeyRepo.List()
	if err != nil {
		ctr.Logger.Errorf("Failed to fetch skill keywords: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to fetch skill keywords",
			Data:    nil,
		})
	}

	response := make([]map[string]interface{}, len(list))
	for i, sk := range list {
		response[i] = map[string]interface{}{
			"id":   sk.ID,
			"name": sk.Name,
		}
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Skill keywords retrieved successfully",
		Data:    response,
	})
}

func (ctr *SkillKeywordController) CreateSkillKeyword(c echo.Context) error {
	var req param.CreateSkillKeywordRequest
	if err := c.Bind(&req); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
		})
	}
	if _, err := valid.ValidateStruct(req); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}

	existing, err := ctr.sKeyRepo.List()
	if err != nil {
		ctr.Logger.Errorf("Failed to check existing skill keywords: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to check existing skill keywords",
		})
	}
	for _, sk := range existing {
		if sk.Name == req.Name {
			return c.JSON(http.StatusOK, cf.JsonResponse{
				Status:  cf.FailResponseCode,
				Message: "Skill keyword name already exists",
			})
		}
	}

	skill := &m.SkillKeyword{Name: req.Name}
	if err := ctr.sKeyRepo.Create(skill); err != nil {
		ctr.Logger.Errorf("Failed to create skill keyword: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to create skill keyword",
		})
	}

	response := map[string]interface{}{
		"id":   skill.ID,
		"name": skill.Name,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Skill keyword created successfully",
		Data:    response,
	})
}

func (ctr *SkillKeywordController) UpdateSkillKeyword(c echo.Context) error {
	var req param.UpdateSkillKeywordRequest
	if err := c.Bind(&req); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
		})
	}
	if _, err := valid.ValidateStruct(req); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}
	skill := &m.SkillKeyword{BaseModel: cm.BaseModel{ID: req.ID}, Name: req.Name}
	if err := ctr.sKeyRepo.Update(skill); err != nil {
		ctr.Logger.Errorf("Failed to update skill keyword: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to update skill keyword",
		})
	}

	response := map[string]interface{}{
		"id":   skill.ID,
		"name": skill.Name,
	}

	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Skill keyword updated successfully",
		Data:    response,
	})
}

func (ctr *SkillKeywordController) DeleteSkillKeyword(c echo.Context) error {
	var req param.DeleteSkillKeywordRequest
	if err := c.Bind(&req); err != nil {
		ctr.Logger.Errorf("Failed to bind params: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Invalid Params",
		})
	}
	if _, err := valid.ValidateStruct(req); err != nil {
		ctr.Logger.Errorf("Validation failed: %v", err)
		return c.JSON(http.StatusOK, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: err.Error(),
		})
	}
	if err := ctr.sKeyRepo.Delete(req.ID); err != nil {
		ctr.Logger.Errorf("Failed to delete skill keyword: %v", err)
		return c.JSON(http.StatusInternalServerError, cf.JsonResponse{
			Status:  cf.FailResponseCode,
			Message: "Failed to delete skill keyword",
		})
	}
	return c.JSON(http.StatusOK, cf.JsonResponse{
		Status:  cf.SuccessResponseCode,
		Message: "Skill keyword deleted successfully",
	})
}
