package skillkeyword

import (
	cm "orientation-training-api/internal/common"
	m "orientation-training-api/internal/models"

	"github.com/labstack/echo/v4"
)

type PgSkillKeywordRepository struct {
	cm.AppRepository
}

func NewPgSkillKeywordRepository(logger echo.Logger) (repo *PgSkillKeywordRepository) {
	repo = &PgSkillKeywordRepository{}
	repo.Init(logger)
	return
}

func (repo *PgSkillKeywordRepository) Create(skill *m.SkillKeyword) error {
	return repo.DB.Insert(skill)
}

func (repo *PgSkillKeywordRepository) GetByID(id int) (*m.SkillKeyword, error) {
	var skill m.SkillKeyword
	err := repo.DB.Model(&skill).Where("id = ?", id).Select()
	if err != nil {
		return nil, err
	}
	return &skill, nil
}

func (repo *PgSkillKeywordRepository) List() ([]m.SkillKeyword, error) {
	var skills []m.SkillKeyword
	err := repo.DB.Model(&skills).Select()
	return skills, err
}

func (repo *PgSkillKeywordRepository) Update(skill *m.SkillKeyword) error {
	_, err := repo.DB.Model(skill).WherePK().Update()
	return err
}

func (repo *PgSkillKeywordRepository) Delete(id int) error {
	skill := &m.SkillKeyword{BaseModel: cm.BaseModel{ID: id}}
	_, err := repo.DB.Model(skill).WherePK().Delete()
	return err
}
